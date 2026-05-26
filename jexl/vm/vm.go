// Copyright (c) 2026 Harness Inc.
// Copyright (c) 2018 Anton Medvedev
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package vm

//go:generate sh -c "go run ./func_types > ./func_types[generated].go"

import (
	"fmt"
	"math"
	"math/big"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/harness/go-jexl/jexl/classes/java/util"
	"github.com/harness/go-jexl/jexl/coerce"
	"github.com/harness/go-jexl/jexl/config"
	"github.com/harness/go-jexl/jexl/internal/deref"
	"github.com/harness/go-jexl/jexl/token"
)

const maxFnArgsBuf = 256

func Run(program *Program, env any) (any, error) {
	if program == nil {
		return nil, fmt.Errorf("program is nil")
	}
	vm := VM{}
	return vm.Run(program, env)
}

// tryFrame records the catch-handler entry point and the stack/scope depth at the time
// OpTry was executed, so OpThrow can restore them before jumping to the handler.
type tryFrame struct {
	catchIP    int    // bytecode offset of the first catch-body instruction
	stackDepth int    // len(vm.Stack) at the time OpTry fired
	scopeDepth int    // len(vm.Scopes) at the time OpTry fired
	catchVar   string // name of the catch variable (empty if unused)
	varSlot    int    // Variables slot for the catch variable (-1 if unused)
}

type VM struct {
	Stack        []any
	Scopes       []*Scope
	Variables    []any
	MemoryBudget uint
	ip           int
	memory       uint
	iterations   uint
	scopePool    []Scope // Pre-allocated pool of Scope values; grows as needed but never shrinks
	scopePoolIdx int     // Current index into scopePool for allocation
	currScope    *Scope  // Cached pointer to the current scope (optimization)
	tryStack     []tryFrame
}

func (vm *VM) Run(program *Program, env any) (_ any, err error) {
	if vm.Stack == nil {
		vm.Stack = make([]any, 0, 2)
	} else {
		clearSlice(vm.Stack)
		vm.Stack = vm.Stack[0:0]
	}
	if vm.Scopes != nil {
		clearSlice(vm.Scopes)
		vm.Scopes = vm.Scopes[0:0]
	}
	vm.scopePoolIdx = 0 // Reset pool index for reuse
	vm.currScope = nil
	vm.tryStack = vm.tryStack[:0]
	if len(vm.Variables) < program.variables {
		vm.Variables = make([]any, program.variables)
	}
	if vm.MemoryBudget == 0 {
		if program.MaxMemory > 0 {
			vm.MemoryBudget = program.MaxMemory
		} else {
			vm.MemoryBudget = config.MaxMemory
		}
	}
	vm.memory = 0
	vm.iterations = 0
	vm.ip = 0

	var fnArgsBuf []any

	for vm.ip < len(program.Bytecode) {
		// runLoop executes until completion, panic, or an active try frame catches the panic.
		// If a panic is caught by an active try frame, vm state is restored to the catch
		// handler and this outer loop continues.
		panicked, panicVal := vm.runLoop(program, env, &fnArgsBuf)
		if !panicked {
			break
		}
		// A panic occurred. If there's an active try frame, route it to the catch handler.
		if len(vm.tryStack) == 0 {
			// No try frame — convert to error.
			var location token.Range
			if vm.ip-1 < len(program.locations) {
				location = program.locations[vm.ip-1]
			}
			f := &token.Error{
				Range:   location,
				Message: fmt.Sprintf("%v", panicVal),
			}
			if e, ok := panicVal.(error); ok {
				f.Wrap(e)
			}
			return nil, f.Bind(program.source)
		}
		// Pop the try frame and jump to the catch handler.
		frame := vm.tryStack[len(vm.tryStack)-1]
		vm.tryStack = vm.tryStack[:len(vm.tryStack)-1]
		// Restore stack and scopes.
		vm.Stack = vm.Stack[:frame.stackDepth]
		vm.Scopes = vm.Scopes[:frame.scopeDepth]
		if len(vm.Scopes) > 0 {
			vm.currScope = vm.Scopes[len(vm.Scopes)-1]
		} else {
			vm.currScope = nil
		}
		// Push the panic value as a string for the catch variable.
		vm.push(fmt.Sprintf("%v", panicVal))
		vm.ip = frame.catchIP
	}

	if len(vm.Stack) > 0 {
		result := vm.pop()
		// BigInteger/BigDecimal values serialize to their canonical string representation.
		if bi, ok := asBigInt(result); ok {
			return bi.String(), nil
		}
		if bd, ok := asBigDec(result); ok {
			return bd.String(), nil
		}
		return result, nil
	}

	return nil, nil
}

// runLoop runs the VM main loop starting at vm.ip. It returns (false, nil) on
// clean completion. On panic it returns (true, recovered_value).
func (vm *VM) runLoop(program *Program, env any, fnArgsBuf *[]any) (panicked bool, panicVal any) {
	defer func() {
		if r := recover(); r != nil {
			panicked = true
			panicVal = r
		}
	}()

	for vm.ip < len(program.Bytecode) {
		op := program.Bytecode[vm.ip]
		arg := program.Arguments[vm.ip]
		vm.ip += 1

		switch op {

		case OpInvalid:
			panic("invalid opcode")

		case OpPush:
			vm.push(program.Constants[arg])

		case OpInt:
			vm.push(arg)

		case OpPop:
			vm.pop()

		case OpStore:
			vm.Variables[arg] = vm.pop()

		case OpLoadVar:
			vm.push(vm.Variables[arg])

		case OpLoadConst:
			key := program.Constants[arg]
			val := Fetch(env, key)
			// Antish fallback: if root env is a flat map and key is absent but a dotted
			// key starting with this prefix exists, use a cursor for chained lookup.
			if val == nil {
				if m, ok := env.(map[string]any); ok {
					if keyStr, isStr := key.(string); isStr {
						if _, exists := m[keyStr]; !exists && HasAntishPrefix(m, keyStr) {
							val = &AntishCursor{Env: m, Path: keyStr}
						}
					}
				}
			}
			vm.push(val)

		case OpLoadField:
			vm.push(FetchField(env, program.Constants[arg].(*Field)))

		case OpLoadFast:
			m := env.(map[string]any)
			key := program.Constants[arg].(string)
			val, exists := m[key]
			if !exists && HasAntishPrefix(m, key) {
				// Key not in context but a dotted key exists — use antish cursor.
				val = &AntishCursor{Env: m, Path: key}
			}
			vm.push(val)

		case OpLoadMethod:
			vm.push(FetchMethod(env, program.Constants[arg].(*Method)))

		case OpLoadFunc:
			vm.push(program.functions[arg])

		case OpFetch:
			b := vm.pop()
			a := vm.pop()
			if methodName, ok := b.(string); ok {
				if obj, ok := a.(interface {
					Call(string, ...any) (any, error)
				}); ok {
					obj := obj
					vm.push(func(args ...any) any {
						result, err := obj.Call(methodName, args...)
						if err != nil {
							panic(err)
						}
						return result
					})
					break
				}
			}
			vm.push(Fetch(a, b))

		case OpFetchField:
			a := vm.pop()
			vm.push(FetchField(a, program.Constants[arg].(*Field)))

		case OpLoadEnv:
			vm.push(env)

		case OpMethod:
			a := vm.pop()
			vm.push(FetchMethod(a, program.Constants[arg].(*Method)))

		case OpTrue:
			vm.push(true)

		case OpFalse:
			vm.push(false)

		case OpNil:
			vm.push(nil)

		case OpNegate:
			val := vm.pop()
			if bi, ok := asBigInt(val); ok {
				vm.push(new(big.Int).Neg(bi))
			} else if bd, ok := asBigDec(val); ok {
				vm.push(bd.Neg())
			} else {
				vm.push(Negate(val))
			}

		case OpNot:
			v := vm.pop().(bool)
			vm.push(!v)

		case OpEqual:
			b := vm.pop()
			a := vm.pop()
			vm.push(looseEqual(a, b))

		case OpEqualInt:
			b := vm.pop()
			a := vm.pop()
			// Fast path for int/int; fall back to equal for float/mixed types.
			ai, aOk := a.(int)
			bi, bOk := b.(int)
			if aOk && bOk {
				vm.push(ai == bi)
			} else {
				vm.push(looseEqual(a, b))
			}

		case OpEqualString:
			b := vm.pop()
			a := vm.pop()
			vm.push(a.(string) == b.(string))

		case OpJump:
			if arg < 0 {
				panic("negative jump offset is invalid")
			}
			vm.ip += arg

		case OpJumpIfTrue:
			if arg < 0 {
				panic("negative jump offset is invalid")
			}
			if !IsFalsy(vm.current()) {
				vm.ip += arg
			}

		case OpJumpIfFalse:
			if arg < 0 {
				panic("negative jump offset is invalid")
			}
			if IsFalsy(vm.current()) {
				vm.ip += arg
			}

		case OpJumpIfNil:
			if arg < 0 {
				panic("negative jump offset is invalid")
			}
			if IsNil(vm.current()) {
				vm.ip += arg
			}

		case OpJumpIfNotNil:
			if arg < 0 {
				panic("negative jump offset is invalid")
			}
			if !IsNil(vm.current()) {
				vm.ip += arg
			}

		case OpJumpIfFalsy:
			if arg < 0 {
				panic("negative jump offset is invalid")
			}
			if IsFalsy(vm.current()) {
				vm.ip += arg
			}

		case OpJumpIfEnd:
			if arg < 0 {
				panic("negative jump offset is invalid")
			}
			if vm.currScope.Index >= vm.currScope.Len {
				vm.ip += arg
			} else if vm.currScope.VarSlot >= 0 {
				vm.Variables[vm.currScope.VarSlot] = vm.currScope.Item()
			}

		case OpJumpBackward:
			vm.iterations++
			if program.MaxIterations > 0 && vm.iterations > program.MaxIterations {
				panic(fmt.Sprintf("maximum iterations (%d) exceeded", program.MaxIterations))
			}
			vm.ip -= arg

		case OpIn:
			b := vm.pop()
			a := vm.pop()
			vm.push(In(a, b))

		case OpInstanceOf:
			typeName := vm.pop().(string)
			a := vm.pop()
			vm.push(instanceOf(a, typeName))

		case OpLess:
			b := vm.pop()
			a := vm.pop()
			if af, ok := a.(float64); ok && math.IsNaN(af) {
				vm.push(true)
			} else if cmp, ok := bigCompare(a, b); ok {
				vm.push(cmp < 0)
			} else {
				vm.push(Less(a, b))
			}

		case OpMore:
			b := vm.pop()
			a := vm.pop()
			if af, ok := a.(float64); ok && math.IsNaN(af) {
				vm.push(false)
			} else if cmp, ok := bigCompare(a, b); ok {
				vm.push(cmp > 0)
			} else {
				vm.push(More(a, b))
			}

		case OpLessOrEqual:
			b := vm.pop()
			a := vm.pop()
			if af, ok := a.(float64); ok && math.IsNaN(af) {
				vm.push(true)
			} else if cmp, ok := bigCompare(a, b); ok {
				vm.push(cmp <= 0)
			} else {
				vm.push(LessOrEqual(a, b))
			}

		case OpMoreOrEqual:
			b := vm.pop()
			a := vm.pop()
			if af, ok := a.(float64); ok && math.IsNaN(af) {
				vm.push(false)
			} else if cmp, ok := bigCompare(a, b); ok {
				vm.push(cmp >= 0)
			} else {
				vm.push(MoreOrEqual(a, b))
			}

		case OpAdd:
			b := vm.pop()
			a := vm.pop()
			if a == nil || b == nil {
				vm.push(nil)
			} else if result, ok := bigArith(a, b, "+"); ok {
				vm.push(result)
			} else if _, aStr := a.(string); aStr {
				// String + anything: coerce b to string.
				vm.push(a.(string) + jexlStringOf(b))
			} else if _, bStr := b.(string); bStr {
				// Anything + string: coerce a to string.
				vm.push(jexlStringOf(a) + b.(string))
			} else {
				a = coerceToNumeric(a)
				b = coerceToNumeric(b)
				// Detect int64 overflow and promote to BigInteger.
				if ai, aOk := toInt64Safe(a); aOk {
					if bi, bOk := toInt64Safe(b); bOk {
						result := ai + bi
						if (bi > 0 && result < ai) || (bi < 0 && result > ai) {
							// overflow: promote to BigInteger
							r := new(big.Int).Add(big.NewInt(ai), big.NewInt(bi))
							vm.push(r)
							break
						}
					}
				}
				vm.push(Add(a, b))
			}

		case OpSubtract:
			b := vm.pop()
			a := vm.pop()
			if result, ok := bigArith(a, b, "-"); ok {
				vm.push(result)
			} else {
				a = coerceToNumeric(a)
				b = coerceToNumeric(b)
				vm.push(Subtract(a, b))
			}

		case OpMultiply:
			b := vm.pop()
			a := vm.pop()
			if result, ok := bigArith(a, b, "*"); ok {
				vm.push(result)
			} else {
				a = coerceToNumeric(a)
				b = coerceToNumeric(b)
				vm.push(Multiply(a, b))
			}

		case OpDivide:
			b := vm.pop()
			a := vm.pop()
			if result, ok := bigArith(a, b, "/"); ok {
				vm.push(result)
			} else {
				if ToFloat64(b) == 0 {
					panic(fmt.Errorf("division by zero"))
				}
				ak := reflect.TypeOf(a).Kind()
				bk := reflect.TypeOf(b).Kind()
				if ak == reflect.Float32 || ak == reflect.Float64 || bk == reflect.Float32 || bk == reflect.Float64 {
					vm.push(Divide(a, b))
				} else {
					vm.push(IntDivide(a, b))
				}
			}

		case OpModulo:
			b := vm.pop()
			a := vm.pop()
			if result, ok := bigArith(a, b, "%"); ok {
				vm.push(result)
			} else {
				// Use float modulo when either operand is a float.
				af, aIsFloat := a.(float64)
				bf, bIsFloat := b.(float64)
				if aIsFloat || bIsFloat {
					if !aIsFloat {
						af = ToFloat64(a)
					}
					if !bIsFloat {
						bf = ToFloat64(b)
					}
					vm.push(math.Mod(af, bf))
				} else {
					vm.push(Modulo(a, b))
				}
			}

		case OpExponent:
			b := vm.pop()
			a := vm.pop()
			vm.push(Exponent(a, b))

		case OpRange:
			b := vm.pop()
			a := vm.pop()
			min := ToInt(a)
			max := ToInt(b)
			var size int
			if min <= max {
				size = max - min + 1
			} else {
				size = min - max + 1
			}
			if size > 0 {
				vm.memGrow(uint(size))
			}
			vm.push(ToRange(min, max))

		case OpMatches:
			b := vm.pop()
			a := vm.pop()
			if IsNil(a) || IsNil(b) {
				vm.push(false)
				break
			}
			if rv, isRegexVal := b.(regexp.Regexp); isRegexVal {
				b = &rv
			}
			vm.push(jexlRegexMatch(a, b))

		case OpInOrMatches:
			b := vm.pop()
			a := vm.pop()
			if IsNil(a) || IsNil(b) {
				vm.push(false)
				break
			}
			if _, isRegex := b.(*regexp.Regexp); isRegex {
				vm.push(jexlRegexMatch(a, b))
			} else if rv, isRegexVal := b.(regexp.Regexp); isRegexVal {
				vm.push(jexlRegexMatch(a, &rv))
			} else if bs, ok := b.(string); ok {
				vm.push(jexlRegexFullMatch(a, bs))
			} else {
				vm.push(In(a, b))
			}

		case OpMatchesConst:
			a := vm.pop()
			if IsNil(a) {
				vm.push(false)
				break
			}
			r := program.Constants[arg].(*regexp.Regexp)
			vm.push(jexlRegexMatch(a, r))

		case OpStartsWith:
			b := vm.pop()
			a := vm.pop()
			if IsNil(a) || IsNil(b) {
				vm.push(false)
				break
			}
			vm.push(strings.HasPrefix(a.(string), b.(string)))

		case OpEndsWith:
			b := vm.pop()
			a := vm.pop()
			if IsNil(a) || IsNil(b) {
				vm.push(false)
				break
			}
			vm.push(strings.HasSuffix(a.(string), b.(string)))

		case OpSlice:
			from := vm.pop()
			to := vm.pop()
			node := vm.pop()
			vm.push(Slice(node, from, to))

		case OpCall:
			v := vm.pop()
			if v == nil {
				panic("invalid operation: cannot call nil")
			}
			if closure, ok := v.(*Closure); ok {
				numArgs := arg
				args := make([]any, numArgs)
				for i := numArgs - 1; i >= 0; i-- {
					args[i] = vm.pop()
				}
				bodyProgram := closure.Program.(*Program)
				callEnv := make(map[string]any, len(closure.Captures)+len(closure.Params)+program.variables)
				if envMap, ok := env.(map[string]any); ok {
					for k, v := range envMap {
						callEnv[k] = v
					}
				}
				for i := 0; i < program.variables; i++ {
					if name, ok := program.debugInfo[fmt.Sprintf("var_%d", i)]; ok {
						callEnv[name] = vm.Variables[i]
					}
				}
				for k, cv := range closure.Captures {
					callEnv[k] = cv
				}
				for i, p := range closure.Params {
					if i < len(args) {
						callEnv[p] = args[i]
					}
				}
				result, err := Run(bodyProgram, callEnv)
				if err != nil {
					panic(err)
				}
				vm.push(result)
				break
			}
			fn := reflect.ValueOf(v)
			if fn.Kind() != reflect.Func {
				panic(fmt.Sprintf("invalid operation: cannot call non-function of type %T", v))
			}
			fnType := fn.Type()
			size := arg
			isVariadic := fnType.IsVariadic()
			numIn := fnType.NumIn()
			if isVariadic {
				if size < numIn-1 {
					panic(fmt.Sprintf("invalid number of arguments: expected at least %d, got %d", numIn-1, size))
				}
			} else {
				if size != numIn {
					panic(fmt.Sprintf("invalid number of arguments: expected %d, got %d", numIn, size))
				}
			}
			in := make([]reflect.Value, size)
			for i := int(size) - 1; i >= 0; i-- {
				param := vm.pop()
				var inType reflect.Type
				if isVariadic && i >= numIn-1 {
					inType = fnType.In(numIn - 1).Elem()
				} else {
					inType = fnType.In(i)
				}
				if param == nil {
					in[i] = reflect.Zero(inType)
				} else {
					rv := reflect.ValueOf(param)
					if rv.Type().ConvertibleTo(inType) && rv.Type() != inType {
						rv = rv.Convert(inType)
					}
					in[i] = rv
				}
			}
			out := fn.Call(in)
			if len(out) == 2 && out[1].Type() == errorType && !out[1].IsNil() {
				panic(out[1].Interface().(error))
			}
			vm.push(out[0].Interface())

		case OpCall0:
			out, err := program.functions[arg]()
			if err != nil {
				panic(err)
			}
			vm.push(out)

		case OpCall1:
			var args []any
			args, *fnArgsBuf = vm.getArgsForFunc(*fnArgsBuf, program, 1)
			out, err := program.functions[arg](args...)
			if err != nil {
				panic(err)
			}
			vm.push(out)

		case OpCall2:
			var args []any
			args, *fnArgsBuf = vm.getArgsForFunc(*fnArgsBuf, program, 2)
			out, err := program.functions[arg](args...)
			if err != nil {
				panic(err)
			}
			vm.push(out)

		case OpCall3:
			var args []any
			args, *fnArgsBuf = vm.getArgsForFunc(*fnArgsBuf, program, 3)
			out, err := program.functions[arg](args...)
			if err != nil {
				panic(err)
			}
			vm.push(out)

		case OpCallN:
			fn := vm.pop().(Function)
			var args []any
			args, *fnArgsBuf = vm.getArgsForFunc(*fnArgsBuf, program, arg)
			out, err := fn(args...)
			if err != nil {
				panic(err)
			}
			vm.push(out)

		case OpCallFast:
			fn := vm.pop().(func(...any) any)
			var args []any
			args, *fnArgsBuf = vm.getArgsForFunc(*fnArgsBuf, program, arg)
			vm.push(fn(args...))

		case OpCallSafe:
			fn := vm.pop().(SafeFunction)
			var args []any
			args, *fnArgsBuf = vm.getArgsForFunc(*fnArgsBuf, program, arg)
			out, mem, err := fn(args...)
			if err != nil {
				panic(err)
			}
			vm.memGrow(mem)
			vm.push(out)

		case OpCallTyped:
			v := vm.pop()
			if closure, ok := v.(*Closure); ok {
				// Determine number of args from the FuncType signature.
				numArgs := 0
				if arg < len(program.Constants) {
					// arg is a FuncTypes index — peek at the type to count params
					_ = arg
				}
				// Fall back: collect args by counting params from closure.
				numArgs = len(closure.Params)
				args := make([]any, numArgs)
				for i := numArgs - 1; i >= 0; i-- {
					args[i] = vm.pop()
				}
				bodyProgram := closure.Program.(*Program)
				callEnv := make(map[string]any, len(closure.Captures)+numArgs+program.variables)
				if envMap, ok2 := env.(map[string]any); ok2 {
					for k, ev := range envMap {
						callEnv[k] = ev
					}
				}
				for i := 0; i < program.variables; i++ {
					if name, ok2 := program.debugInfo[fmt.Sprintf("var_%d", i)]; ok2 {
						callEnv[name] = vm.Variables[i]
					}
				}
				for k, cv := range closure.Captures {
					callEnv[k] = cv
				}
				for i, p := range closure.Params {
					if i < len(args) {
						callEnv[p] = args[i]
					}
				}
				result, err := Run(bodyProgram, callEnv)
				if err != nil {
					panic(err)
				}
				vm.push(result)
				break
			}
			vm.push(vm.call(v, arg))

		case OpArray:
			size := vm.pop().(int)
			vm.memGrow(uint(size))
			array := make([]any, size)
			for i := size - 1; i >= 0; i-- {
				array[i] = vm.pop()
			}
			vm.push(array)

		case OpMap:
			size := vm.pop().(int)
			vm.memGrow(uint(size))
			// Collect pairs first to determine key types.
			pairs := make([][2]any, size)
			hasNonStringKey := false
			for i := size - 1; i >= 0; i-- {
				pairs[i][1] = vm.pop()
				pairs[i][0] = vm.pop()
				if _, isStr := pairs[i][0].(string); !isStr {
					hasNonStringKey = true
				}
			}
			if hasNonStringKey {
				m := make(map[any]any, size)
				for _, p := range pairs {
					m[p[0]] = p[1]
				}
				vm.push(m)
			} else {
				m := make(map[string]any, size)
				for _, p := range pairs {
					m[p[0].(string)] = p[1]
				}
				vm.push(m)
			}

		case OpDeref:
			a := vm.pop()
			vm.push(deref.Interface(a))

		case OpIncrementIndex:
			vm.currScope.Index++

		case OpDecrementIndex:
			vm.currScope.Index--

		case OpIncrementCount:
			vm.currScope.Count++

		case OpGetIndex:
			vm.push(vm.currScope.Index)

		case OpGetCount:
			vm.push(vm.currScope.Count)

		case OpGetLen:
			vm.push(vm.currScope.Len)

		case OpGetAcc:
			vm.push(vm.currScope.Acc)

		case OpSetAcc:
			vm.currScope.Acc = vm.pop()

		case OpSetIndex:
			vm.currScope.Index = vm.pop().(int)

		case OpPointer:
			vm.push(vm.currScope.Item())

		case OpTry:
			if arg == 0 {
				// Pop try frame: try body exited normally.
				if len(vm.tryStack) > 0 {
					vm.tryStack = vm.tryStack[:len(vm.tryStack)-1]
				}
			} else {
				// Push try frame; catchIP is relative to the instruction following OpTry.
				vm.tryStack = append(vm.tryStack, tryFrame{
					catchIP:    vm.ip + arg,
					stackDepth: len(vm.Stack),
					scopeDepth: len(vm.Scopes),
				})
			}

		case OpThrow:
			val := vm.pop()
			if len(vm.tryStack) > 0 {
				frame := vm.tryStack[len(vm.tryStack)-1]
				vm.tryStack = vm.tryStack[:len(vm.tryStack)-1]
				// Restore stack to the depth it was when OpTry fired.
				vm.Stack = vm.Stack[:frame.stackDepth]
				// Restore scopes to the depth they were when OpTry fired.
				vm.Scopes = vm.Scopes[:frame.scopeDepth]
				if len(vm.Scopes) > 0 {
					vm.currScope = vm.Scopes[len(vm.Scopes)-1]
				} else {
					vm.currScope = nil
				}
				// Push the thrown value for the catch block to consume (e.g., OpStore).
				vm.push(val)
				// Jump to catch handler (absolute IP).
				vm.ip = frame.catchIP
			} else {
				// No enclosing try: escape the script as a panic.
				panic(&ThrownValue{Value: val})
			}

		case OpBegin:
			a := vm.pop()
			s := vm.allocScope()
			switch v := a.(type) {
			case []int:
				s.Ints = v
				s.Len = len(v)
			case []float64:
				s.Floats = v
				s.Len = len(v)
			case []string:
				s.Strings = v
				s.Len = len(v)
			case []any:
				s.Anys = v
				s.Len = len(v)
			default:
				s.Array = reflect.ValueOf(a)
				s.Len = s.Array.Len()
			}
			vm.Scopes = append(vm.Scopes, s)
			vm.currScope = s

		case OpAnd:
			a := vm.pop()
			b := vm.pop()
			vm.push(a.(bool) && b.(bool))

		case OpOr:
			a := vm.pop()
			b := vm.pop()
			vm.push(a.(bool) || b.(bool))

		case OpBitOr:
			b := vm.pop()
			a := vm.pop()
			result, err := BitwiseOr(a, b)
			if err != nil {
				panic(err)
			}
			vm.push(result)

		case OpBitXor:
			b := vm.pop()
			a := vm.pop()
			result, err := BitwiseXor(a, b)
			if err != nil {
				panic(err)
			}
			vm.push(result)

		case OpBitAnd:
			b := vm.pop()
			a := vm.pop()
			result, err := BitwiseAnd(a, b)
			if err != nil {
				panic(err)
			}
			vm.push(result)

		case OpBitNot:
			a := vm.pop()
			result, err := BitwiseNot(a)
			if err != nil {
				panic(err)
			}
			vm.push(result)

		case OpShiftLeft:
			b := vm.pop()
			a := vm.pop()
			result, err := ShiftLeft(a, b)
			if err != nil {
				panic(err)
			}
			vm.push(result)

		case OpShiftRight:
			b := vm.pop()
			a := vm.pop()
			result, err := ShiftRight(a, b)
			if err != nil {
				panic(err)
			}
			vm.push(result)

		case OpShiftRightU:
			b := vm.pop()
			a := vm.pop()
			result, err := ShiftRightUnsigned(a, b)
			if err != nil {
				panic(err)
			}
			vm.push(result)

		case OpStrictEqual:
			b := vm.pop()
			a := vm.pop()
			vm.push(StrictEqual(a, b))

		case OpStrictNotEqual:
			b := vm.pop()
			a := vm.pop()
			vm.push(!StrictEqual(a, b))

		case OpLambda:
			tmpl := program.Constants[arg].(*Closure)
			closure := &Closure{Params: tmpl.Params, Program: tmpl.Program}
			// Capture outer variables at definition time (capture-by-value semantics).
			// First snapshot the env (which contains params and outer env vars when inside a closure call).
			if envMap, ok := env.(map[string]any); ok && len(envMap) > 0 {
				closure.Captures = make(map[string]any, len(envMap)+len(tmpl.CaptureVars))
				for k, v := range envMap {
					closure.Captures[k] = v
				}
			} else if len(tmpl.CaptureVars) > 0 {
				closure.Captures = make(map[string]any, len(tmpl.CaptureVars))
			}
			// Then snapshot local variables (these override env to get current values).
			for _, cv := range tmpl.CaptureVars {
				if closure.Captures == nil {
					closure.Captures = make(map[string]any, len(tmpl.CaptureVars))
				}
				closure.Captures[cv.Name] = vm.Variables[cv.Slot]
			}
			vm.push(closure)

		case OpCallLambda:
			numArgs := arg
			args := make([]any, numArgs)
			for i := numArgs - 1; i >= 0; i-- {
				args[i] = vm.pop()
			}
			rawClosure := vm.pop()
			closure, ok := rawClosure.(*Closure)
			if !ok {
				// Not a closure — treat as a constant value (e.g., filter(arr, false)).
				vm.push(rawClosure)
				break
			}
			bodyProgram := closure.Program.(*Program)
			callEnv := make(map[string]any, len(closure.Captures)+len(closure.Params)+program.variables)
			// Merge current env so closures and recursive calls can see outer names.
			if envMap, ok := env.(map[string]any); ok {
				for k, v := range envMap {
					callEnv[k] = v
				}
			}
			// Expose outer-scope named variables (take precedence over env).
			for i := 0; i < program.variables; i++ {
				if name, ok := program.debugInfo[fmt.Sprintf("var_%d", i)]; ok {
					callEnv[name] = vm.Variables[i]
				}
			}
			for k, v := range closure.Captures {
				callEnv[k] = v
			}
			// Params override everything.
			for i, p := range closure.Params {
				if i < len(args) {
					callEnv[p] = args[i]
				}
			}
			result, err := Run(bodyProgram, callEnv)
			if err != nil {
				panic(err)
			}
			vm.push(result)

		case OpAssign:
			// arg is the Variables slot index; pop value, store, push it back (assignment is an expression)
			val := vm.pop()
			vm.Variables[arg] = val
			vm.push(val)

		case OpCompoundAssign:
			// arg indexes program.Constants for a CompoundAssignOp
			op := program.Constants[arg].(CompoundAssignOp)
			rhs := vm.pop()
			lhs := vm.Variables[op.VarIndex]
			var result any
			switch op.Op {
			case "+=":
				if _, lhsStr := lhs.(string); lhsStr {
					if r, ok := rhs.(rune); ok {
						result = lhs.(string) + string(r)
					} else {
						result = lhs.(string) + coerce.ToString(rhs)
					}
				} else if _, rhsStr := rhs.(string); rhsStr {
					if r, ok := lhs.(rune); ok {
						result = string(r) + rhs.(string)
					} else {
						result = coerce.ToString(lhs) + rhs.(string)
					}
				} else {
					result = Add(lhs, rhs)
				}
			case "-=":
				result = Subtract(lhs, rhs)
			case "*=":
				result = Multiply(lhs, rhs)
			case "/=":
				result = Divide(lhs, rhs)
			case "%=":
				result = Modulo(lhs, rhs)
			case "&=":
				r, err := BitwiseAnd(lhs, rhs)
				if err != nil {
					panic(err)
				}
				result = r
			case "|=":
				r, err := BitwiseOr(lhs, rhs)
				if err != nil {
					panic(err)
				}
				result = r
			case "^=":
				r, err := BitwiseXor(lhs, rhs)
				if err != nil {
					panic(err)
				}
				result = r
			case "<<=":
				r, err := ShiftLeft(lhs, rhs)
				if err != nil {
					panic(err)
				}
				result = r
			case ">>=":
				r, err := ShiftRight(lhs, rhs)
				if err != nil {
					panic(err)
				}
				result = r
			case ">>>=":
				r, err := ShiftRightUnsigned(lhs, rhs)
				if err != nil {
					panic(err)
				}
				result = r
			default:
				panic(fmt.Sprintf("unknown compound assignment operator %q", op.Op))
			}
			vm.Variables[op.VarIndex] = result
			vm.push(result)

		case OpIncrement:
			old := vm.Variables[arg]
			newVal, err := Increment(old)
			if err != nil {
				panic(err)
			}
			vm.Variables[arg] = newVal
			vm.push(newVal)

		case OpDecrement:
			old := vm.Variables[arg]
			newVal, err := Decrement(old)
			if err != nil {
				panic(err)
			}
			vm.Variables[arg] = newVal
			vm.push(newVal)

		case OpForEach:
			// arg is the Variables slot for the loop variable.
			a := vm.pop()
			// Convert Java collections to slices for iteration.
			a = toIterableSlice(a)
			s := vm.allocScope()
			switch v := a.(type) {
			case []int:
				s.Ints = v
				s.Len = len(v)
			case []float64:
				s.Floats = v
				s.Len = len(v)
			case []string:
				s.Strings = v
				s.Len = len(v)
			case []any:
				s.Anys = v
				s.Len = len(v)
			default:
				s.Array = reflect.ValueOf(a)
				s.Len = s.Array.Len()
			}
			s.VarSlot = arg
			vm.Scopes = append(vm.Scopes, s)
			vm.currScope = s

		case OpEnd:
			vm.Scopes = vm.Scopes[:len(vm.Scopes)-1]
			if len(vm.Scopes) > 0 {
				vm.currScope = vm.Scopes[len(vm.Scopes)-1]
			} else {
				vm.currScope = nil
			}

		case OpBreak:
			// arg is a forward jump offset; compiler emits any necessary scope cleanup before this.
			if arg < 0 {
				panic("negative jump offset is invalid")
			}
			vm.ip += arg

		case OpContinue:
			// arg is a backward jump offset; compiler emits any necessary scope cleanup before this.
			vm.ip -= arg

		case OpReturn:
			// Unwind the script: set ip past the end so the run loop exits.
			vm.ip = len(program.Bytecode)

		case OpEmpty:
			a := vm.pop()
			vm.push(Empty(a))

		case OpSize:
			a := vm.pop()
			vm.push(Size(a))

		case OpSet:
			size := arg
			vm.memGrow(uint(size))
			elements := make([]any, size)
			for i := size - 1; i >= 0; i-- {
				elements[i] = vm.pop()
			}
			vm.push(util.NewHashSetFrom(elements))

		case OpNew:
			numArgs := arg
			args := make([]any, numArgs)
			for i := numArgs - 1; i >= 0; i-- {
				args[i] = vm.pop()
			}
			name := vm.pop().(string)
			obj, ok := program.Registry.Lookup(name)
			if !ok {
				panic(fmt.Sprintf("unknown class %q", name))
			}
			out, err := obj.Call("new", args...)
			if err != nil {
				panic(err)
			}
			vm.push(out)

		case OpStoreIndex:
			val := vm.pop()
			key := vm.pop()
			obj := vm.pop()
			SetIndex(obj, key, val)

		case OpCompoundStoreIndex:
			op := program.Constants[arg].(CompoundAssignOp)
			rhs := vm.pop()
			key := vm.pop()
			obj := vm.pop()
			lhs := Fetch(obj, key)
			result := applyCompoundOp(op.Op, lhs, rhs)
			SetIndex(obj, key, result)
			vm.push(result)

		case OpCompoundStoreEnv:
			op := program.Constants[arg].(CompoundEnvAssignOp)
			rhs := vm.pop()
			lhs := Fetch(env, op.Key)
			result := applyCompoundOp(op.Op, lhs, rhs)
			if m, ok := env.(map[string]any); ok {
				m[op.Key] = result
			}
			vm.push(result)

		default:
			panic(fmt.Sprintf("unknown bytecode %#x", op))
		}

	}

	return false, nil
}

func (vm *VM) push(value any) {
	vm.Stack = append(vm.Stack, value)
}

func (vm *VM) current() any {
	if len(vm.Stack) == 0 {
		panic("stack underflow")
	}
	return vm.Stack[len(vm.Stack)-1]
}

func (vm *VM) pop() any {
	if len(vm.Stack) == 0 {
		panic("stack underflow")
	}
	value := vm.Stack[len(vm.Stack)-1]
	vm.Stack = vm.Stack[:len(vm.Stack)-1]
	return value
}

func (vm *VM) memGrow(size uint) {
	vm.memory += size
	if vm.memory >= vm.MemoryBudget {
		panic("memory budget exceeded")
	}
}

func (vm *VM) scope() *Scope {
	return vm.Scopes[len(vm.Scopes)-1]
}

// allocScope returns a pointer to a Scope from the pool, growing the pool if needed.
// Callers must set Len and exactly one of: Ints, Floats, Strings, Anys, or Array.
func (vm *VM) allocScope() *Scope {
	if vm.scopePoolIdx >= len(vm.scopePool) {
		vm.scopePool = append(vm.scopePool, Scope{})
	}
	s := &vm.scopePool[vm.scopePoolIdx]
	vm.scopePoolIdx++
	// Reset iteration state
	s.Index = 0
	s.Count = 0
	s.Acc = nil
	s.VarSlot = -1
	// Clear typed slice pointers to avoid stale fast-path matches
	s.Ints = nil
	s.Floats = nil
	s.Strings = nil
	s.Anys = nil
	// Clear Array to release reference for GC (only matters for fallback path)
	s.Array = reflect.Value{}
	return s
}

// getArgsForFunc lazily initializes the buffer the first time it is called for
// a given program (thus, it also needs "program" to run). It will
// take "needed" elements from the buffer and populate them with vm.pop() in
// reverse order. Because the estimation can fall short, this function can
// occasionally make a new allocation.
func (vm *VM) getArgsForFunc(argsBuf []any, program *Program, needed int) (args []any, argsBufOut []any) {
	if needed == 0 || program == nil {
		return nil, argsBuf
	}

	// Step 1: fix estimations and preallocate
	if argsBuf == nil {
		estimatedFnArgsCount := estimateFnArgsCount(program)
		if estimatedFnArgsCount > maxFnArgsBuf {
			// put a practical limit to avoid excessive preallocation
			estimatedFnArgsCount = maxFnArgsBuf
		}
		if estimatedFnArgsCount < needed {
			// in the case that the first call is for example OpCallN with a large
			// number of arguments, then make sure we will be able to serve them at
			// least.
			estimatedFnArgsCount = needed
		}

		// in the case that we are preparing the arguments for the first
		// function call of the program, then argsBuf will be nil, so we
		// initialize it. We delay this initial allocation here because a
		// program could have many function calls but exit earlier than the
		// first call, so in that case we avoid allocating unnecessarily
		argsBuf = make([]any, estimatedFnArgsCount)
	}

	// Step 2: get the final slice that will be returned
	var buf []any
	if len(argsBuf) >= needed {
		// in this case, we are successfully using the single preallocation. We
		// use the full slice expression [low : high : max] because in that way
		// a function that receives this slice as variadic arguments will not be
		// able to make modifications to contiguous elements with append(). If
		// they call append on their variadic arguments they will make a new
		// allocation.
		buf = (argsBuf)[:needed:needed]
		argsBuf = (argsBuf)[needed:] // advance the buffer
	} else {
		// if we have been making calls to something like OpCallN with many more
		// arguments than what we estimated, then we will need to allocate
		// separately
		buf = make([]any, needed)
	}

	// Step 3: populate the final slice bulk copying from the stack. This is the
	// exact order and copy() is a highly optimized operation
	copy(buf, vm.Stack[len(vm.Stack)-needed:])
	vm.Stack = vm.Stack[:len(vm.Stack)-needed]

	return buf, argsBuf
}

// CompoundAssignOp is stored in program.Constants for OpCompoundAssign.
// VarIndex is the Variables slot; Op is the operator string (e.g. "+=").
type CompoundAssignOp struct {
	VarIndex int
	Op       string
}

// CompoundEnvAssignOp is stored in program.Constants for OpCompoundStoreEnv.
type CompoundEnvAssignOp struct {
	Key string
	Op  string
}

func applyCompoundOp(op string, lhs, rhs any) any {
	switch op {
	case "+=":
		if _, lhsStr := lhs.(string); lhsStr {
			return lhs.(string) + coerce.ToString(rhs)
		} else if _, rhsStr := rhs.(string); rhsStr {
			return coerce.ToString(lhs) + rhs.(string)
		}
		return Add(lhs, rhs)
	case "-=":
		return Subtract(lhs, rhs)
	case "*=":
		return Multiply(lhs, rhs)
	case "/=":
		return Divide(lhs, rhs)
	case "%=":
		return Modulo(lhs, rhs)
	case "&=":
		r, err := BitwiseAnd(lhs, rhs)
		if err != nil {
			panic(err)
		}
		return r
	case "|=":
		r, err := BitwiseOr(lhs, rhs)
		if err != nil {
			panic(err)
		}
		return r
	case "^=":
		r, err := BitwiseXor(lhs, rhs)
		if err != nil {
			panic(err)
		}
		return r
	case "<<=":
		r, err := ShiftLeft(lhs, rhs)
		if err != nil {
			panic(err)
		}
		return r
	case ">>=":
		r, err := ShiftRight(lhs, rhs)
		if err != nil {
			panic(err)
		}
		return r
	case ">>>=":
		r, err := ShiftRightUnsigned(lhs, rhs)
		if err != nil {
			panic(err)
		}
		return r
	default:
		panic(fmt.Sprintf("unknown compound assignment operator %q", op))
	}
}

func clearSlice[S ~[]E, E any](s S) {
	var zero E
	for i := range s {
		s[i] = zero // clear mem, optimized by the compiler, in Go 1.21 the "clear" builtin can be used
	}
}

// estimateFnArgsCount inspects a *Program and estimates how many function
// arguments will be required to run it.
func estimateFnArgsCount(program *Program) int {
	// Implementation note: a program will not necessarily go through all
	// operations, but this is just an estimation
	var count int
	for _, op := range program.Bytecode {
		if int(op) < len(opArgLenEstimation) {
			count += opArgLenEstimation[op]
		}
	}
	return count
}

func instanceOf(a any, typeName string) bool {
	if a == nil {
		return false
	}
	switch typeName {
	case "String":
		_, ok := a.(string)
		return ok
	case "Boolean":
		_, ok := a.(bool)
		return ok
	case "Integer", "Long", "Short", "Byte":
		switch a.(type) {
		case int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32, uint64,
			float64:
			rv := reflect.ValueOf(a)
			switch rv.Kind() {
			case reflect.Float64, reflect.Float32:
				f := rv.Float()
				return f == float64(int64(f))
			default:
				return true
			}
		}
		return false
	case "Float", "Double":
		switch a.(type) {
		case float32, float64, int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32, uint64:
			return true
		}
		return false
	case "Number":
		switch a.(type) {
		case int, int8, int16, int32, int64,
			uint, uint8, uint16, uint32, uint64,
			float32, float64:
			return true
		}
		return false
	case "Map":
		return reflect.TypeOf(a).Kind() == reflect.Map
	case "List", "Array":
		k := reflect.TypeOf(a).Kind()
		return k == reflect.Slice || k == reflect.Array
	default:
		// For unknown type names, check the Go type name.
		return reflect.TypeOf(a).Name() == typeName
	}
}

// jexlStringOf converts a value to string following JEXL/Java conventions.
// Float64 whole numbers are displayed as "N.0" (Java Double.toString behavior).
func jexlStringOf(v any) string {
	if f, ok := v.(float64); ok {
		if math.IsNaN(f) {
			return "NaN"
		}
		if math.IsInf(f, 1) {
			return "Infinity"
		}
		if math.IsInf(f, -1) {
			return "-Infinity"
		}
		s := strconv.FormatFloat(f, 'f', -1, 64)
		if !strings.Contains(s, ".") {
			s += ".0"
		}
		return s
	}
	return coerce.ToString(v)
}

// looseEqual implements JEXL loose equality: number vs string comparison uses toString coercion.
func looseEqual(a, b any) bool {
	// BigDec/BigInt equality: use bigCompare.
	if cmp, ok := bigCompare(a, b); ok {
		return cmp == 0
	}
	// String vs numeric: compare by converting number to string.
	aStr, aIsStr := a.(string)
	bStr, bIsStr := b.(string)
	if aIsStr && !bIsStr && b != nil {
		return aStr == coerce.ToString(b)
	}
	if bIsStr && !aIsStr && a != nil {
		return coerce.ToString(a) == bStr
	}
	return Equal(a, b)
}

// jexlRegexFullMatch performs Java JEXL-style full-string regex match using a pattern string.
func jexlRegexFullMatch(a any, pattern string) bool {
	anchored := "^(?:" + pattern + ")$"
	r, err := regexp.Compile(anchored)
	if err != nil {
		panic(err)
	}
	if s, ok := a.(string); ok {
		return r.MatchString(s)
	}
	if b, ok := a.([]byte); ok {
		return r.Match(b)
	}
	return false
}

// jexlRegexMatch matches using a compiled *regexp.Regexp or a string pattern (full-match).
func jexlRegexMatch(a, b any) bool {
	switch r := b.(type) {
	case *regexp.Regexp:
		fullPattern := "^(?:" + r.String() + ")$"
		fr, err := regexp.Compile(fullPattern)
		if err != nil {
			// Fallback to original
			fr = r
		}
		if s, ok := a.(string); ok {
			return fr.MatchString(s)
		}
		if bs, ok := a.([]byte); ok {
			return fr.Match(bs)
		}
		return false
	case string:
		return jexlRegexFullMatch(a, r)
	}
	return false
}

var opArgLenEstimation = [...]int{
	OpCall1: 1,
	OpCall2: 2,
	OpCall3: 3,
	// we don't know exactly but we know at least 4, so be conservative as this
	// is only an optimization and we also want to avoid excessive preallocation
	OpCallN: 4,
	// here we don't know either, but we can guess it could be common to receive
	// up to 3 arguments in a function
	OpCallFast: 3,
	OpCallSafe: 3,
}
