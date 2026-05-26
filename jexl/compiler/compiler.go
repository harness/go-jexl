// Copyright (c) 2026 Harness Inc.
// Copyright (c) 2018 Anton Medvedev
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package compiler

import (
	"fmt"
	"math"
	"math/big"
	"reflect"
	"regexp"
	"runtime/debug"

	"github.com/harness/go-jexl/jexl/ast"
	"github.com/harness/go-jexl/jexl/checker"
	. "github.com/harness/go-jexl/jexl/checker/nature"
	"github.com/harness/go-jexl/jexl/config"
	"github.com/harness/go-jexl/jexl/token"
	. "github.com/harness/go-jexl/jexl/vm"
	"github.com/harness/go-jexl/jexl/vm/runtime"
)

const (
	// placeholder is a dummy jump argument emitted before the
	// real target is known; always overwritten by patchJump.
	placeholder = 12345
)

var bigIntPtrType = reflect.TypeOf((*big.Int)(nil))

// Compile translates an AST into a Program. Any panic from the
// compiler (e.g. overflow, unknown node) is caught and returned
// as an error together with a stack trace.
func Compile(tree *ast.Tree, config *config.Config) (program *Program, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v\n%s", r, debug.Stack())
		}
	}()

	c := &compiler{
		config:         config,
		locations:      make([]token.Range, 0),
		constantsIndex: make(map[any]int),
		functionsIndex: make(map[string]int),
		debugInfo:      make(map[string]string),
	}

	// Share the config's Nature cache when available, otherwise
	// allocate a fresh one so reflect lookups still work.
	if config != nil {
		c.ntCache = &c.config.Cache
	} else {
		c.ntCache = new(Cache)
	}

	c.compile(tree.Node)

	// Optimise only when we have full config (type info is needed
	// to safely chain jump targets).
	if c.config != nil {
		c.optimize()
	}

	program = NewProgram(
		tree.Source,
		c.locations,
		c.variables,
		c.constants,
		c.bytecode,
		c.arguments,
		c.functions,
		c.debugInfo,
	)
	if c.config != nil {
		program.Registry = c.config.Registry
		program.MaxIterations = c.config.MaxIterations
		program.MaxMemory = c.config.MaxMemory
	}
	return
}

// compiler emits bytecode for a single expression or program.
// One compiler instance is used per compilation; sub-compilers
// are created for lambda bodies (see LambdaNode).
type compiler struct {
	config         *config.Config
	ntCache        *Cache        // shared Nature cache for type lookups
	locations      []token.Range // source location per instruction
	bytecode       []Opcode
	variables      int           // total local variable slots allocated
	scopes         []scope       // lexical variable stack (innermost last)
	loopStack      []loopFrame   // one frame per active loop/switch
	constants      []any         // constant pool (literals, Field/Method, etc.)
	constantsIndex map[any]int   // dedup index for hashable constants
	functions      []Function    // function table (index used by OpCall*)
	functionsIndex map[string]int
	debugInfo      map[string]string // human-readable names for vars/funcs
	nodes          []ast.Node        // current compile path (for location tracking)
	chains         [][]int           // stack of optional-chain jump placeholders
	arguments      []int             // argument per instruction (parallel to bytecode)
}

// scope is one entry on the lexical variable stack —
// a name bound to a VM slot index.
type scope struct {
	variableName string
	index        int
	isConst      bool
	isClosure    bool // variable holds a Closure (from a lambda/function declaration)
}

// loopFrame tracks the bytecode bookkeeping for one active loop
// or switch: where to jump back (begin), and where break/continue land.
type loopFrame struct {
	isForeach     bool
	hasPostExpr   bool  // for-loop has a post expression; continue must skip to it
	begin         int   // bytecode index of the loop's back-jump target
	breakHoles    []int // forward-jump placeholders for break; patched to loop end
	continueHoles []int // forward-jump placeholders for continue; patched to continueTarget
}

// nodeParent returns the node one level up in the current
// compile path, used to peek at the call context (e.g. detect
// whether a ChainExpr sits inside a ?? operator).
func (c *compiler) nodeParent() ast.Node {
	if len(c.nodes) > 1 {
		return c.nodes[len(c.nodes)-2]
	}
	return nil
}

// emitLocation appends one instruction with an explicit source
// location and returns the index *after* the opcode (i.e. the
// index of its argument), which is used as a patch handle.
func (c *compiler) emitLocation(loc token.Range, op Opcode, arg int) int {
	c.bytecode = append(c.bytecode, op)
	current := len(c.bytecode)
	c.arguments = append(c.arguments, arg)
	c.locations = append(c.locations, loc)
	return current
}

// emit appends one instruction using the current node's location.
// Returns the patch handle (argument index) for forward jumps.
func (c *compiler) emit(op Opcode, args ...int) int {
	arg := 0
	if len(args) > 1 {
		panic("too many arguments")
	}
	if len(args) == 1 {
		arg = args[0]
	}
	var loc token.Range
	if len(c.nodes) > 0 {
		loc = c.nodes[len(c.nodes)-1].Location()
	}
	return c.emitLocation(loc, op, arg)
}

// emitPush interns value into the constant pool and emits OpPush.
func (c *compiler) emitPush(value any) int {
	return c.emit(OpPush, c.addConstant(value))
}

// addConstant interns constant into the pool and returns its index.
// Slices, maps, structs, and funcs are not hashable so each call
// appends a new entry. *Field and *Method are stringified for
// deduplication so two field references with the same index path
// share one pool slot.
func (c *compiler) addConstant(constant any) int {
	indexable := true
	hash := constant
	switch reflect.TypeOf(constant).Kind() {
	case reflect.Slice, reflect.Map, reflect.Struct, reflect.Func:
		indexable = false
	}
	if field, ok := constant.(*Field); ok {
		indexable = true
		hash = fmt.Sprintf("%v", field)
	}
	if method, ok := constant.(*Method); ok {
		indexable = true
		hash = fmt.Sprintf("%v", method)
	}
	if indexable {
		if p, ok := c.constantsIndex[hash]; ok {
			return p
		}
	}
	c.constants = append(c.constants, constant)
	p := len(c.constants) - 1
	if indexable {
		c.constantsIndex[hash] = p
	}
	return p
}

// addVariable allocates the next variable slot and records its
// human-readable name in debugInfo for the disassembler.
func (c *compiler) addVariable(name string) int {
	c.variables++
	c.debugInfo[fmt.Sprintf("var_%d", c.variables-1)] = name
	return c.variables - 1
}

// emitFunction registers fn and emits the most specific call
// opcode for argsLen (OpCall0–OpCall3 for 0–3 args; OpCallN
// for more, which requires a preceding OpLoadFunc).
func (c *compiler) emitFunction(fn *runtime.Function, argsLen int) {
	switch argsLen {
	case 0:
		c.emit(OpCall0, c.addFunction(fn.Name, fn.Func))
	case 1:
		c.emit(OpCall1, c.addFunction(fn.Name, fn.Func))
	case 2:
		c.emit(OpCall2, c.addFunction(fn.Name, fn.Func))
	case 3:
		c.emit(OpCall3, c.addFunction(fn.Name, fn.Func))
	default:
		c.emit(OpLoadFunc, c.addFunction(fn.Name, fn.Func))
		c.emit(OpCallN, argsLen)
	}
}

// addFunction adds Function.Func to the program.functions and returns its index.
func (c *compiler) addFunction(name string, fn Function) int {
	if fn == nil {
		panic("function is nil")
	}
	if p, ok := c.functionsIndex[name]; ok {
		return p
	}
	p := len(c.functions)
	c.functions = append(c.functions, fn)
	c.functionsIndex[name] = p
	c.debugInfo[fmt.Sprintf("func_%d", p)] = name
	return p
}

// patchJump writes the correct forward-jump offset into the
// argument slot identified by placeholder (the patch handle
// returned by emit). The VM computes: ip += arg + 1.
func (c *compiler) patchJump(placeholder int) {
	offset := len(c.bytecode) - placeholder
	c.arguments[placeholder-1] = offset
}

// calcBackwardJump returns the argument for an OpJumpBackward
// that should land at bytecode index "to".
// The VM computes: ip -= arg - 1; so arg = (current+1) - to.
func (c *compiler) calcBackwardJump(to int) int {
	return len(c.bytecode) + 1 - to
}

// compile pushes node onto the path for location tracking,
// dispatches to the typed handler, then pops it.
func (c *compiler) compile(node ast.Node) {
	c.nodes = append(c.nodes, node)
	defer func() {
		c.nodes = c.nodes[:len(c.nodes)-1]
	}()

	switch n := node.(type) {
	case *ast.NilLit:
		c.NilNode(n)
	case *ast.Ident:
		c.IdentifierNode(n)
	case *ast.IntLit:
		c.IntegerNode(n)
	case *ast.FloatLit:
		c.FloatNode(n)
	case *ast.BoolLit:
		c.BoolNode(n)
	case *ast.StringLit:
		c.StringNode(n)
	case *ast.ConstantExpr:
		c.ConstantNode(n)
	case *ast.UnaryExpr:
		c.UnaryNode(n)
	case *ast.BinaryExpr:
		c.BinaryNode(n)
	case *ast.ChainExpr:
		c.ChainNode(n)
	case *ast.MemberExpr:
		c.MemberNode(n)
	case *ast.CallExpr:
		c.CallNode(n)
	case *ast.PredicateExpr:
		c.PredicateNode(n)
	case *ast.ConditionalExpr:
		c.ConditionalNode(n)
	case *ast.ArrayLit:
		c.ArrayNode(n)
	case *ast.MapLit:
		c.MapNode(n)
	case *ast.KeyValueExpr:
		c.PairNode(n)
	case *ast.BitExpr:
		c.BitOpNode(n)
	case *ast.StrictEqualExpr:
		c.StrictEqualNode(n)
	case *ast.LambdaExpr:
		c.LambdaNode(n)
	case *ast.BlockStmt:
		c.BlockNode(n)
	case *ast.VarDecl:
		c.VarNode(n)
	case *ast.AssignStmt:
		c.AssignNode(n)
	case *ast.IncDecStmt:
		c.IncrDecrNode(n)
	case *ast.ForeachStmt:
		c.ForeachNode(n)
	case *ast.WhileStmt:
		c.WhileNode(n)
	case *ast.ForStmt:
		c.ForNode(n)
	case *ast.DoWhileStmt:
		c.DoWhileNode(n)
	case *ast.BreakStmt:
		c.BreakNode(n)
	case *ast.ContinueStmt:
		c.ContinueNode(n)
	case *ast.ReturnStmt:
		c.ReturnNode(n)
	case *ast.TryStmt:
		c.TryNode(n)
	case *ast.ThrowStmt:
		c.ThrowNode(n)
	case *ast.SwitchStmt:
		c.SwitchNode(n)
	case *ast.SetLit:
		c.SetNode(n)
	case *ast.EmptyExpr:
		c.EmptyNode(n)
	case *ast.SizeExpr:
		c.SizeNode(n)
	case *ast.RegexLit:
		c.RegexNode(n)
	case *ast.NewExpr:
		c.NewNode(n)
	case *ast.NamespaceCallExpr:
		c.NamespaceCallNode(n)
	default:
		panic(fmt.Sprintf("undefined node type (%T)", node))
	}
}

// NilNode pushes the nil sentinel onto the stack.
func (c *compiler) NilNode(_ *ast.NilLit) {
	c.emit(OpNil)
}

// IdentifierNode resolves a name to the fastest available opcode:
// local variable slot, $env shortcut, numeric constant (NaN/Infinity),
// map key, struct field index, method index, or generic name lookup.
func (c *compiler) IdentifierNode(node *ast.Ident) {
	// Local variable takes priority over the environment.
	if index, ok := c.lookupVariable(node.Value); ok {
		c.emit(OpLoadVar, index)
		return
	}
	// $env is a special reference to the whole environment map.
	if node.Value == "$env" {
		c.emit(OpLoadEnv)
		return
	}
	// Well-known numeric constants.
	switch node.Value {
	case "NaN":
		c.emitPush(math.NaN())
		return
	case "Infinity":
		c.emitPush(math.Inf(1))
		return
	}

	var env Nature
	if c.config != nil {
		env = c.config.Context
	}

	if env.IsFastMap() {
		// map[string]any env — key lookup at runtime by name.
		c.emit(OpLoadFast, c.addConstant(node.Value))
	} else if ok, index, name := checker.FieldIndex(c.ntCache, env, node); ok {
		// Struct field — emit a direct field access by reflect index.
		c.emit(OpLoadField, c.addConstant(&Field{
			Index: index,
			Path:  []string{name},
		}))
	} else if ok, index, name := checker.MethodIndex(c.ntCache, env, node); ok {
		// Method on the env — emit a direct method load by index.
		c.emit(OpLoadMethod, c.addConstant(&Method{
			Name:  name,
			Index: index,
		}))
	} else {
		// Generic env lookup by name (slow path via reflection).
		c.emit(OpLoadConst, c.addConstant(node.Value))
	}
}

// IntegerNode pushes an integer literal. If the checker narrowed
// the type (e.g. int32, float64), it casts and range-checks first.
func (c *compiler) IntegerNode(node *ast.IntLit) {
	t := node.Type()
	if t == nil {
		// No type annotation — push as plain int.
		c.emitPush(node.Value)
		return
	}
	// The checker may have narrowed the type to match a function
	// parameter (e.g. int32 or float64). Cast and range-check here
	// so the VM sees the correct concrete type.
	switch t.Kind() {
	case reflect.Float32:
		c.emitPush(float32(node.Value))
	case reflect.Float64:
		c.emitPush(float64(node.Value))
	case reflect.Int:
		c.emitPush(node.Value)
	case reflect.Int8:
		if node.Value > math.MaxInt8 || node.Value < math.MinInt8 {
			panic(fmt.Sprintf("constant %d overflows int8", node.Value))
		}
		c.emitPush(int8(node.Value))
	case reflect.Int16:
		if node.Value > math.MaxInt16 || node.Value < math.MinInt16 {
			panic(fmt.Sprintf("constant %d overflows int16", node.Value))
		}
		c.emitPush(int16(node.Value))
	case reflect.Int32:
		if node.Value > math.MaxInt32 || node.Value < math.MinInt32 {
			panic(fmt.Sprintf("constant %d overflows int32", node.Value))
		}
		c.emitPush(int32(node.Value))
	case reflect.Int64:
		c.emitPush(int64(node.Value))
	case reflect.Uint:
		if node.Value < 0 {
			panic(fmt.Sprintf("constant %d overflows uint", node.Value))
		}
		c.emitPush(uint(node.Value))
	case reflect.Uint8:
		if node.Value > math.MaxUint8 || node.Value < 0 {
			panic(fmt.Sprintf("constant %d overflows uint8", node.Value))
		}
		c.emitPush(uint8(node.Value))
	case reflect.Uint16:
		if node.Value > math.MaxUint16 || node.Value < 0 {
			panic(fmt.Sprintf("constant %d overflows uint16", node.Value))
		}
		c.emitPush(uint16(node.Value))
	case reflect.Uint32:
		if node.Value < 0 {
			panic(fmt.Sprintf("constant %d overflows uint32", node.Value))
		}
		c.emitPush(uint32(node.Value))
	case reflect.Uint64:
		if node.Value < 0 {
			panic(fmt.Sprintf("constant %d overflows uint64", node.Value))
		}
		c.emitPush(uint64(node.Value))
	default:
		c.emitPush(node.Value)
	}
}

// FloatNode pushes a float literal, narrowing to float32 when the
// checker annotated the node with that type.
func (c *compiler) FloatNode(node *ast.FloatLit) {
	switch node.Type().Kind() {
	case reflect.Float32:
		c.emitPush(float32(node.Value))
	case reflect.Float64:
		c.emitPush(node.Value)
	default:
		c.emitPush(node.Value)
	}
}

// BoolNode emits OpTrue or OpFalse — no constant pool needed.
func (c *compiler) BoolNode(node *ast.BoolLit) {
	if node.Value {
		c.emit(OpTrue)
	} else {
		c.emit(OpFalse)
	}
}

// StringNode pushes the string value into the constant pool.
func (c *compiler) StringNode(node *ast.StringLit) {
	c.emitPush(node.Value)
}

// ConstantNode pushes a pre-evaluated constant (e.g. from a patcher).
// nil values emit OpNil instead of going through the constant pool.
func (c *compiler) ConstantNode(node *ast.ConstantExpr) {
	if node.Value == nil {
		c.emit(OpNil)
		return
	}
	c.emitPush(node.Value)
}

// compileRegex converts a RegexLit (with optional flags) into a
// compiled *regexp.Regexp. Only i, m, s flags are mapped to Go's
// inline flag syntax (?ims); unsupported flags are silently dropped.
func (c *compiler) compileRegex(node *ast.RegexLit) *regexp.Regexp {
	pattern := node.Pattern
	if node.Flags != "" {
		goFlags := ""
		for _, f := range node.Flags {
			switch f {
			case 'i', 'm', 's':
				goFlags += string(f)
			}
		}
		if goFlags != "" {
			pattern = "(?" + goFlags + ")" + pattern
		}
	}
	re, err := regexp.Compile(pattern)
	if err != nil {
		panic(fmt.Errorf("invalid regex literal /%s/%s: %w", node.Pattern, node.Flags, err))
	}
	return re
}

// RegexNode compiles the regex literal and pushes the *regexp.Regexp
// into the constant pool; any syntax error panics here at compile time.
func (c *compiler) RegexNode(node *ast.RegexLit) {
	c.emitPush(c.compileRegex(node))
}

// UnaryNode compiles the operand and emits the operator opcode.
// Unary + is a no-op at the bytecode level.
func (c *compiler) UnaryNode(node *ast.UnaryExpr) {
	c.compile(node.Node)
	c.derefInNeeded(node.Node)

	switch node.Operator {

	case "!", "not":
		c.emit(OpNot)

	case "+":
		// Do nothing

	case "-":
		c.emit(OpNegate)

	default:
		panic(fmt.Sprintf("unknown operator (%v)", node.Operator))
	}
}

// BinaryNode compiles both operands and emits the corresponding opcode.
// Short-circuit operators (||, &&, ??) use conditional jump sequences
// rather than evaluating both sides unconditionally.
func (c *compiler) BinaryNode(node *ast.BinaryExpr) {
	switch node.Operator {
	case "==":
		c.equalBinaryNode(node)

	case "!=":
		// Reuse the equality path then negate the result.
		c.equalBinaryNode(node)
		c.emit(OpNot)

	case "or", "||":
		// Short-circuit: if left is truthy, leave it on stack and skip right.
		// OpJumpIfTrue does NOT pop; we pop manually on the false path.
		c.compile(node.Left)
		c.derefInNeeded(node.Left)
		end := c.emit(OpJumpIfTrue, placeholder)
		c.emit(OpPop)
		c.compile(node.Right)
		c.derefInNeeded(node.Right)
		c.patchJump(end)

	case "and", "&&":
		// Short-circuit: if left is falsy, leave it on stack and skip right.
		c.compile(node.Left)
		c.derefInNeeded(node.Left)
		end := c.emit(OpJumpIfFalse, placeholder)
		c.emit(OpPop)
		c.compile(node.Right)
		c.derefInNeeded(node.Right)
		c.patchJump(end)

	case "<":
		c.compile(node.Left)
		c.derefInNeeded(node.Left)
		c.compile(node.Right)
		c.derefInNeeded(node.Right)
		c.emit(OpLess)

	case ">":
		c.compile(node.Left)
		c.derefInNeeded(node.Left)
		c.compile(node.Right)
		c.derefInNeeded(node.Right)
		c.emit(OpMore)

	case "<=":
		c.compile(node.Left)
		c.derefInNeeded(node.Left)
		c.compile(node.Right)
		c.derefInNeeded(node.Right)
		c.emit(OpLessOrEqual)

	case ">=":
		c.compile(node.Left)
		c.derefInNeeded(node.Left)
		c.compile(node.Right)
		c.derefInNeeded(node.Right)
		c.emit(OpMoreOrEqual)

	case "+":
		c.compile(node.Left)
		c.derefInNeeded(node.Left)
		c.compile(node.Right)
		c.derefInNeeded(node.Right)
		c.emit(OpAdd)

	case "-":
		c.compile(node.Left)
		c.derefInNeeded(node.Left)
		c.compile(node.Right)
		c.derefInNeeded(node.Right)
		c.emit(OpSubtract)

	case "*":
		c.compile(node.Left)
		c.derefInNeeded(node.Left)
		c.compile(node.Right)
		c.derefInNeeded(node.Right)
		c.emit(OpMultiply)

	case "/":
		c.compile(node.Left)
		c.derefInNeeded(node.Left)
		c.compile(node.Right)
		c.derefInNeeded(node.Right)
		c.emit(OpDivide)

	case "%":
		c.compile(node.Left)
		c.derefInNeeded(node.Left)
		c.compile(node.Right)
		c.derefInNeeded(node.Right)
		c.emit(OpModulo)

	case "**", "^":
		c.compile(node.Left)
		c.derefInNeeded(node.Left)
		c.compile(node.Right)
		c.derefInNeeded(node.Right)
		c.emit(OpExponent)

	case "in":
		c.compile(node.Left)
		c.derefInNeeded(node.Left)
		c.compile(node.Right)
		c.derefInNeeded(node.Right)
		c.emit(OpIn)

	case "=~":
		// Regex match operator. Dispatch depends on what the RHS is:
		//   /pattern/  → compile once, emit OpMatchesConst
		//   "string"   → try to compile as regex; fall back to OpIn
		//   [array]    → membership test (OpIn)
		//   other      → OpInOrMatches decides at runtime
		if rn, ok := node.Right.(*ast.RegexLit); ok {
			re := c.compileRegex(rn)
			c.compile(node.Left)
			c.derefInNeeded(node.Left)
			c.emit(OpMatchesConst, c.addConstant(re))
		} else if str, ok := node.Right.(*ast.StringLit); ok {
			re, err := regexp.Compile(str.Value)
			if err != nil {
				// Invalid regex string — treat as OpIn membership.
				c.compile(node.Left)
				c.derefInNeeded(node.Left)
				c.compile(node.Right)
				c.derefInNeeded(node.Right)
				c.emit(OpIn)
			} else {
				c.compile(node.Left)
				c.derefInNeeded(node.Left)
				c.emit(OpMatchesConst, c.addConstant(re))
			}
		} else if _, ok := node.Right.(*ast.ArrayLit); ok {
			c.compile(node.Left)
			c.derefInNeeded(node.Left)
			c.compile(node.Right)
			c.derefInNeeded(node.Right)
			c.emit(OpIn)
		} else {
			// Dynamic RHS: runtime decides if it's a collection or regex string.
			c.compile(node.Left)
			c.derefInNeeded(node.Left)
			c.compile(node.Right)
			c.derefInNeeded(node.Right)
			c.emit(OpInOrMatches)
		}

	case "!~":
		// Same dispatch as =~ then negate.
		if rn, ok := node.Right.(*ast.RegexLit); ok {
			re := c.compileRegex(rn)
			c.compile(node.Left)
			c.derefInNeeded(node.Left)
			c.emit(OpMatchesConst, c.addConstant(re))
			c.emit(OpNot)
		} else if str, ok := node.Right.(*ast.StringLit); ok {
			re, err := regexp.Compile(str.Value)
			if err != nil {
				c.compile(node.Left)
				c.derefInNeeded(node.Left)
				c.compile(node.Right)
				c.derefInNeeded(node.Right)
				c.emit(OpIn)
				c.emit(OpNot)
			} else {
				c.compile(node.Left)
				c.derefInNeeded(node.Left)
				c.emit(OpMatchesConst, c.addConstant(re))
				c.emit(OpNot)
			}
		} else if _, ok := node.Right.(*ast.ArrayLit); ok {
			c.compile(node.Left)
			c.derefInNeeded(node.Left)
			c.compile(node.Right)
			c.derefInNeeded(node.Right)
			c.emit(OpIn)
			c.emit(OpNot)
		} else {
			c.compile(node.Left)
			c.derefInNeeded(node.Left)
			c.compile(node.Right)
			c.derefInNeeded(node.Right)
			c.emit(OpInOrMatches)
			c.emit(OpNot)
		}

	case "instanceof":
		// RHS is typically a bare Ident (class name string); push
		// it as a constant rather than resolving it as a variable.
		c.compile(node.Left)
		c.derefInNeeded(node.Left)
		if id, ok := node.Right.(*ast.Ident); ok {
			c.emitPush(id.Value)
		} else {
			c.compile(node.Right)
			c.derefInNeeded(node.Right)
		}
		c.emit(OpInstanceOf)

	case "!instanceof":
		c.compile(node.Left)
		c.derefInNeeded(node.Left)
		if id, ok := node.Right.(*ast.Ident); ok {
			c.emitPush(id.Value)
		} else {
			c.compile(node.Right)
			c.derefInNeeded(node.Right)
		}
		c.emit(OpInstanceOf)
		c.emit(OpNot)

	case "=^":
		c.compile(node.Left)
		c.derefInNeeded(node.Left)
		c.compile(node.Right)
		c.derefInNeeded(node.Right)
		c.emit(OpStartsWith)

	case "!^":
		c.compile(node.Left)
		c.derefInNeeded(node.Left)
		c.compile(node.Right)
		c.derefInNeeded(node.Right)
		c.emit(OpStartsWith)
		c.emit(OpNot)

	case "=$":
		c.compile(node.Left)
		c.derefInNeeded(node.Left)
		c.compile(node.Right)
		c.derefInNeeded(node.Right)
		c.emit(OpEndsWith)

	case "!$":
		c.compile(node.Left)
		c.derefInNeeded(node.Left)
		c.compile(node.Right)
		c.derefInNeeded(node.Right)
		c.emit(OpEndsWith)
		c.emit(OpNot)

	case "..":
		c.compile(node.Left)
		c.derefInNeeded(node.Left)
		c.compile(node.Right)
		c.derefInNeeded(node.Right)
		c.emit(OpRange)

	case "??":
		// Null-coalescing: if left is non-nil, keep it; else use right.
		// OpJumpIfNotNil does NOT pop, so we pop manually on the nil path.
		c.compile(node.Left)
		c.derefInNeeded(node.Left)
		end := c.emit(OpJumpIfNotNil, placeholder)
		c.emit(OpPop)
		c.compile(node.Right)
		c.derefInNeeded(node.Right)
		c.patchJump(end)

	default:
		panic(fmt.Sprintf("unknown operator (%v)", node.Operator))

	}
}

// equalBinaryNode compiles == (and != via the caller adding OpNot).
// When both sides are the same primitive kind with no package path
// (i.e. a built-in type, not a named type from a package) we emit
// the faster OpEqualInt / OpEqualString instead of the generic
// reflect-based OpEqual.
func (c *compiler) equalBinaryNode(node *ast.BinaryExpr) {
	l := kind(node.Left.Type())
	r := kind(node.Right.Type())

	leftIsSimple := isSimpleType(node.Left)
	rightIsSimple := isSimpleType(node.Right)
	leftAndRightAreSimple := leftIsSimple && rightIsSimple

	c.compile(node.Left)
	c.derefInNeeded(node.Left)
	c.compile(node.Right)
	c.derefInNeeded(node.Right)

	if l == r && l == reflect.Int && leftAndRightAreSimple {
		c.emit(OpEqualInt)
	} else if l == r && l == reflect.String && leftAndRightAreSimple {
		c.emit(OpEqualString)
	} else {
		c.emit(OpEqual)
	}
}

// isSimpleType returns true when the node has a non-nil type with
// no package path — i.e. a built-in Go kind (int, string, bool …).
// Named types (e.g. type MyInt int) have a PkgPath and must go
// through the generic OpEqual path.
func isSimpleType(node ast.Node) bool {
	if node == nil {
		return false
	}
	t := node.Type()
	if t == nil {
		return false
	}
	return t.PkgPath() == ""
}

// ChainNode compiles an optional-access chain (a?.b?.c).
// It pushes a new jump-hole list onto c.chains so that each
// optional member/call inside the chain can register its nil-exit
// jump. When any step is nil all registered jumps land here and
// the chain produces nil.
func (c *compiler) ChainNode(node *ast.ChainExpr) {
	c.chains = append(c.chains, []int{})
	c.compile(node.Node)
	// Patch all nil-exit holes to land at the current position.
	for _, ph := range c.chains[len(c.chains)-1] {
		c.patchJump(ph)
	}
	parent := c.nodeParent()
	if binary, ok := parent.(*ast.BinaryExpr); ok && binary.Operator == "??" {
		// Parent is a ?? operator — it will consume whatever is on
		// the stack, so we don't need to push an explicit nil here.
	} else {
		// Replace whatever "typed nil" landed on the stack with an
		// untyped nil so callers get a clean nil result.
		j := c.emit(OpJumpIfNotNil, placeholder)
		c.emit(OpPop)
		c.emit(OpNil)
		c.patchJump(j)
	}
	c.chains = c.chains[:len(c.chains)-1]
}

// MemberNode compiles a property or index access (a.b, a[i]).
// It tries to use the fastest available opcode:
//
//  1. Direct method call → OpMethod with pre-computed index.
//  2. Struct field access → walk up the base chain to merge
//     multi-level field index paths into a single OpLoadField,
//     skipping intermediate loads entirely.
//  3. Everything else → compile base + property, emit OpFetch.
func (c *compiler) MemberNode(node *ast.MemberExpr) {
	var env Nature
	if c.config != nil {
		env = c.config.Context
	}

	// Fast path 1: method with a known index on the env type.
	if ok, index, name := checker.MethodIndex(c.ntCache, env, node); ok {
		c.compile(node.Node)
		c.emit(OpMethod, c.addConstant(&Method{
			Name:  name,
			Index: index,
		}))
		return
	}
	op := OpFetch
	base := node.Node

	ok, index, nodeName := checker.FieldIndex(c.ntCache, env, node)
	path := []string{nodeName}

	if ok {
		// Fast path 2: struct field. Walk up the chain of base
		// nodes merging field indices so the VM can reach the
		// final field in one step without intermediate loads.
		// Stop if any step in the chain is optional (?.) because
		// we can't skip the nil check.
		op = OpFetchField
		for !node.Optional {
			if ident, isIdent := base.(*ast.Ident); isIdent {
				if ok, identIndex, name := checker.FieldIndex(c.ntCache, env, ident); ok {
					index = append(identIndex, index...)
					path = append([]string{name}, path...)
					c.emitLocation(ident.Location(), OpLoadField, c.addConstant(
						&Field{Index: index, Path: path},
					))
					return
				}
			}

			if member, isMember := base.(*ast.MemberExpr); isMember {
				if ok, memberIndex, name := checker.FieldIndex(c.ntCache, env, member); ok {
					index = append(memberIndex, index...)
					path = append([]string{name}, path...)
					node = member
					base = member.Node
				} else {
					break
				}
			} else {
				break
			}
		}
	}

	c.compile(base)
	// Optional access inside a chain: emit a nil-guard jump and
	// register the placeholder so ChainNode can patch it.
	// Outside a chain (no c.chains) the optional flag is ignored.
	if node.Optional && len(c.chains) > 0 {
		ph := c.emit(OpJumpIfNil, placeholder)
		c.chains[len(c.chains)-1] = append(c.chains[len(c.chains)-1], ph)
	}

	if op == OpFetch {
		c.compile(node.Property)
		// When the map key type is a pointer, skip the deref so
		// the pointer value itself is used as the key.
		deref := true
		if node.Node.Type() != nil && node.Node.Type().Kind() == reflect.Map {
			keyType := node.Node.Type().Key()
			propType := node.Property.Type()
			if propType != nil && propType.AssignableTo(keyType) {
				deref = false
			}
		}
		if deref {
			c.derefInNeeded(node.Property)
		}
		c.emit(OpFetch)
	} else {
		c.emitLocation(node.Location(), op, c.addConstant(
			&Field{Index: index, Path: path},
		))
	}
}

// CallNode compiles a function or method call. Dispatch priority:
//
//  1. Inline lambda call / closure variable → OpCallLambda
//  2. Named builtin (config.Functions) → emitFunction
//  3. Registry class method call → emitFunction wrapping obj.Call
//  4. Callee with known func type → choose the fastest call opcode:
//     OpCallTyped (pre-matched signature), OpCallFast (variadic any),
//     or generic OpCall.
//
// For cases 4+ the callee itself is compiled first, leaving the
// function value on the stack for the call opcode to consume.
func (c *compiler) CallNode(node *ast.CallExpr) {
	// Immediately-invoked lambda: (x => x+1)(2)
	if _, ok := node.Callee.(*ast.LambdaExpr); ok {
		c.compile(node.Callee)
		for _, arg := range node.Arguments {
			c.compile(arg)
		}
		c.emit(OpCallLambda, len(node.Arguments))
		return
	}
	// Variable that holds a Closure (declared with var f = x => …)
	if ident, ok := node.Callee.(*ast.Ident); ok && c.isClosureVariable(ident.Value) {
		c.compile(node.Callee)
		for _, arg := range node.Arguments {
			c.compile(arg)
		}
		c.emit(OpCallLambda, len(node.Arguments))
		return
	}

	fn := node.Callee.Type()
	if fn.Kind() == reflect.Func {
		// Compile arguments with parameter-type deref coercion.
		// If the callee is a method, the receiver is In(0) and
		// user args start at In(1) — fnInOffset accounts for this.
		fnInOffset := 0
		fnNumIn := fn.NumIn()
		switch callee := node.Callee.(type) {
		case *ast.MemberExpr:
			if prop, ok := callee.Property.(*ast.StringLit); ok {
				if _, ok = callee.Node.Type().MethodByName(prop.Value); ok && callee.Node.Type().Kind() != reflect.Interface {
					fnInOffset = 1
					fnNumIn--
				}
			}
		case *ast.Ident:
			if t, ok := c.config.Context.MethodByName(c.ntCache, callee.Value); ok && t.Method {
				fnInOffset = 1
				fnNumIn--
			}
		}
		for i, arg := range node.Arguments {
			c.compile(arg)
			var in reflect.Type
			if fn.IsVariadic() && i >= fnNumIn-1 {
				// Variadic tail: use the elem type of the last param.
				in = fn.In(fn.NumIn() - 1).Elem()
			} else {
				in = fn.In(i + fnInOffset)
			}
			c.derefParam(in, arg)
		}
	} else {
		// Unknown func type — compile args without coercion.
		for _, arg := range node.Arguments {
			c.compile(arg)
		}
	}

	// Named builtin from config.Functions.
	if ident, ok := node.Callee.(*ast.Ident); ok {
		if c.config != nil {
			if fn, ok := c.config.Functions[ident.Value]; ok {
				c.emitFunction(fn, len(node.Arguments))
				return
			}
		}
	}

	// Registry class static/instance call (e.g. MyClass.staticFn()).
	// Wrap the registry dispatch in a Function closure so the VM
	// can call it via the normal function table.
	if c.config != nil && c.config.Registry != nil {
		if className, method, ok := extractClassCall(node.Callee); ok {
			if obj, found := c.config.Registry.LookupClass(className); found {
				m := method // capture loop var for closure
				rfn := &runtime.Function{
					Name: className + "." + method,
					Func: func(args ...any) (any, error) { return obj.Call(m, args...) },
				}
				c.emitFunction(rfn, len(node.Arguments))
				return
			}
		}
	}

	// Generic path: push the callee onto the stack, then pick the
	// most specific call opcode available.
	c.compile(node.Callee)

	if c.config != nil {
		isMethod, _, _ := checker.MethodIndex(c.ntCache, c.config.Context, node.Callee)
		if index, ok := checker.TypedFuncIndex(node.Callee.Type(), isMethod); ok {
			// Signature matched a pre-compiled FuncTypes slot — no reflection.
			c.emit(OpCallTyped, index)
			return
		} else if checker.IsFastFunc(node.Callee.Type(), isMethod) {
			// Variadic func(...any) any — slightly faster than generic call.
			c.emit(OpCallFast, len(node.Arguments))
		} else {
			c.emit(OpCall, len(node.Arguments))
		}
	} else {
		c.emit(OpCall, len(node.Arguments))
	}
}

// PredicateNode compiles the inner expression; the predicate wrapper
// itself carries no bytecode — it exists only for type-checking.
func (c *compiler) PredicateNode(node *ast.PredicateExpr) {
	c.compile(node.Node)
}

// beginScope pushes a new scope entry for the given variable slot.
func (c *compiler) beginScope(name string, index int, isConst ...bool) {
	s := scope{variableName: name, index: index}
	if len(isConst) > 0 {
		s.isConst = isConst[0]
	}
	c.scopes = append(c.scopes, s)
}

// endScope pops the innermost scope entry.
func (c *compiler) endScope() {
	c.scopes = c.scopes[:len(c.scopes)-1]
}

// lookupVariable searches scopes innermost-first and returns
// the slot index and true if the name is a declared local variable.
func (c *compiler) lookupVariable(name string) (int, bool) {
	for i := len(c.scopes) - 1; i >= 0; i-- {
		if c.scopes[i].variableName == name {
			return c.scopes[i].index, true
		}
	}
	return 0, false
}

// isConstVariable returns true when name was declared with `const`,
// so AssignNode can reject a write to it at compile time.
func (c *compiler) isConstVariable(name string) bool {
	for i := len(c.scopes) - 1; i >= 0; i-- {
		if c.scopes[i].variableName == name {
			return c.scopes[i].isConst
		}
	}
	return false
}

// isClosureVariable returns true when name was declared with a
// lambda initializer so CallNode emits OpCallLambda instead of OpCall.
func (c *compiler) isClosureVariable(name string) bool {
	for i := len(c.scopes) - 1; i >= 0; i-- {
		if c.scopes[i].variableName == name {
			return c.scopes[i].isClosure
		}
	}
	return false
}

// ConditionalNode compiles a ternary (cond ? a : b) or elvis (a ?: b).
// Elvis uses OpJumpIfFalsy so any falsy value (nil, 0, "") triggers
// the right branch, matching JEXL semantics.
func (c *compiler) ConditionalNode(node *ast.ConditionalExpr) {
	c.compile(node.Cond)
	c.derefInNeeded(node.Cond)

	// elvis (?:): Exp1 is the same node as Cond; use falsy jump for any type
	isElvis := node.Exp1 == node.Cond
	var otherwise int
	if isElvis {
		otherwise = c.emit(OpJumpIfFalsy, placeholder)
	} else {
		otherwise = c.emit(OpJumpIfFalse, placeholder)
	}

	c.emit(OpPop)
	c.compile(node.Exp1)
	end := c.emit(OpJump, placeholder)

	c.patchJump(otherwise)
	c.emit(OpPop)
	if node.Exp2 != nil {
		c.compile(node.Exp2)
	} else {
		c.emit(OpNil)
	}

	c.patchJump(end)
}

// ArrayNode compiles each element onto the stack, then emits
// OpArray with the element count so the VM builds the slice.
func (c *compiler) ArrayNode(node *ast.ArrayLit) {
	for _, node := range node.Nodes {
		c.compile(node)
	}

	c.emitPush(len(node.Nodes))
	c.emit(OpArray)
}

// MapNode compiles each key-value pair onto the stack via PairNode,
// then emits OpMap with the pair count so the VM builds the map.
func (c *compiler) MapNode(node *ast.MapLit) {
	for _, pair := range node.Pairs {
		c.compile(pair)
	}

	c.emitPush(len(node.Pairs))
	c.emit(OpMap)
}

// PairNode pushes the key then the value onto the stack;
// MapNode's OpMap pops them in pairs to build the map.
func (c *compiler) PairNode(node *ast.KeyValueExpr) {
	c.compile(node.Key)
	c.compile(node.Value)
}

// BitOpNode compiles bitwise and shift operators. Unary ~ has
// no RHS node; all other operators compile both operands.
func (c *compiler) BitOpNode(node *ast.BitExpr) {
	c.compile(node.Left)
	c.derefInNeeded(node.Left)
	if node.Right != nil {
		c.compile(node.Right)
		c.derefInNeeded(node.Right)
	}
	switch node.Operator {
	case "|":
		c.emit(OpBitOr)
	case "^":
		c.emit(OpBitXor)
	case "&":
		c.emit(OpBitAnd)
	case "~":
		c.emit(OpBitNot)
	case "<<":
		c.emit(OpShiftLeft)
	case ">>":
		c.emit(OpShiftRight)
	case ">>>":
		c.emit(OpShiftRightU)
	default:
		panic(fmt.Sprintf("unknown bitwise operator (%v)", node.Operator))
	}
}

// StrictEqualNode emits OpStrictEqual or OpStrictNotEqual. Unlike
// OpEqual, these never coerce types — identity comparison only.
func (c *compiler) StrictEqualNode(node *ast.StrictEqualExpr) {
	c.compile(node.Left)
	c.derefInNeeded(node.Left)
	c.compile(node.Right)
	c.derefInNeeded(node.Right)
	if node.Negated {
		c.emit(OpStrictNotEqual)
	} else {
		c.emit(OpStrictEqual)
	}
}

// collectOuterVarRefs walks the lambda body and returns one
// CaptureVar for each identifier that resolves to a local
// variable in the enclosing scope (a "free variable"). Lambda
// params are excluded — they are not captured, they are passed.
func (c *compiler) collectOuterVarRefs(body ast.Node, params []string) []CaptureVar {
	paramSet := make(map[string]bool, len(params))
	for _, p := range params {
		paramSet[p] = true
	}
	seen := make(map[string]bool)
	var captures []CaptureVar
	c.walkForCaptures(body, paramSet, seen, &captures)
	return captures
}

// walkForCaptures recurses into node looking for identifiers that
// resolve to outer local variables and adds them to out.
func (c *compiler) walkForCaptures(node ast.Node, params, seen map[string]bool, out *[]CaptureVar) {
	if node == nil {
		return
	}
	if ident, ok := node.(*ast.Ident); ok {
		name := ident.Value
		if !seen[name] && !params[name] {
			for i := len(c.scopes) - 1; i >= 0; i-- {
				if c.scopes[i].variableName == name {
					seen[name] = true
					*out = append(*out, CaptureVar{Name: name, Slot: c.scopes[i].index})
					break
				}
			}
		}
	}
	// Recurse into child nodes via the visitor walker.
	w := &captureWalker{c: c, params: params, seen: seen, out: out}
	nd := ast.Node(node)
	ast.Walk(&nd, w)
}

// captureWalker implements ast.Visitor for collectOuterVarRefs,
// recording every outer-scope identifier referenced in a lambda body.
type captureWalker struct {
	c      *compiler
	params map[string]bool
	seen   map[string]bool
	out    *[]CaptureVar
}

// Visit is called by ast.Walk for every node in the lambda body.
// It records outer-scope variable references into w.out.
func (w *captureWalker) Visit(node *ast.Node) {
	if *node == nil {
		return
	}
	if ident, ok := (*node).(*ast.Ident); ok {
		name := ident.Value
		if !w.seen[name] && !w.params[name] {
			for i := len(w.c.scopes) - 1; i >= 0; i-- {
				if w.c.scopes[i].variableName == name {
					w.seen[name] = true
					*w.out = append(*w.out, CaptureVar{Name: name, Slot: w.c.scopes[i].index})
					break
				}
			}
		}
	}
}

// LambdaNode compiles the lambda body into a separate sub-Program
// and emits OpLambda with a Closure template as the constant.
// At runtime OpLambda instantiates the closure, copying the current
// values of any captured outer variables into the new Closure.
func (c *compiler) LambdaNode(node *ast.LambdaExpr) {
	// Use a fresh compiler for the body so bytecode, constants,
	// and scopes don't bleed into the enclosing program.
	sub := &compiler{
		config:         c.config,
		locations:      make([]token.Range, 0),
		constantsIndex: make(map[any]int),
		functionsIndex: make(map[string]int),
		debugInfo:      make(map[string]string),
	}
	if c.config != nil {
		sub.ntCache = &c.config.Cache
	} else {
		sub.ntCache = new(Cache)
	}
	// Params are not in sub.scopes — the VM passes them via a
	// callEnv map, and the body reads them with OpLoadConst (name
	// lookup), not OpLoadVar (slot lookup).
	sub.compile(node.Body)

	bodyProgram := NewProgram(
		token.File{},
		sub.locations,
		sub.variables,
		sub.constants,
		sub.bytecode,
		sub.arguments,
		sub.functions,
		sub.debugInfo,
	)
	// Collect free variables from the enclosing scope so OpLambda
	// can snapshot their current slot values into the Closure.
	captureVars := c.collectOuterVarRefs(node.Body, node.Params)
	tmpl := &Closure{
		Params:      node.Params,
		Program:     bodyProgram,
		CaptureVars: captureVars,
	}
	c.emit(OpLambda, c.addConstant(tmpl))
}

// BlockNode compiles each statement in sequence. The final
// statement's value remains on the stack as the block's result.
func (c *compiler) BlockNode(node *ast.BlockStmt) {
	for _, stmt := range node.Statements {
		c.compile(stmt)
	}
}

// VarNode compiles the initializer, stores the value into a new
// slot, and pushes a scope entry. Lambda initializers set isClosure
// so CallNode can emit OpCallLambda instead of OpCall.
func (c *compiler) VarNode(node *ast.VarDecl) {
	c.compile(node.Value)
	index := c.addVariable(node.Name)
	c.emit(OpStore, index) // pop value into variable slot
	// Mark as closure so CallNode knows to use OpCallLambda.
	_, isClosure := node.Value.(*ast.LambdaExpr)
	s := scope{variableName: node.Name, index: index, isConst: node.Keyword == "const", isClosure: isClosure}
	c.scopes = append(c.scopes, s)
}

// AssignNode compiles simple and compound assignment (=, +=, -= …).
// Targets can be local variable slots or member/index lvalues.
// In non-strict mode, assigning an unknown name writes to the env map.
func (c *compiler) AssignNode(node *ast.AssignStmt) {
	switch target := node.Target.(type) {
	case *ast.Ident:
		if c.isConstVariable(target.Value) {
			panic(fmt.Sprintf("cannot assign to const variable: %s", target.Value))
		}
		index, found := c.lookupVariable(target.Value)
		if !found {
			// Name not in local scope. In non-strict mode, fall back
			// to env (context) assignment so scripts can mutate the
			// environment map passed in from the host.
			isStrict := c.config == nil || c.config.Strict
			if isStrict {
				panic(fmt.Sprintf("undefined variable: %s", target.Value))
			}
			if c.config != nil && c.config.Safe {
				panic(fmt.Sprintf("cannot assign to context variable in safe mode: %s", target.Value))
			}
			if node.Op == "=" {
				// Stack layout for OpStoreIndex: [env, key, value]
				c.emit(OpLoadEnv)
				c.emit(OpPush, c.addConstant(target.Value))
				c.compile(node.Value)
				c.derefInNeeded(node.Value)
				c.emit(OpStoreIndex, 0)
			} else {
				// Compound env assignment (+=, -= …); the opcode
				// reads the current value, applies the op, and stores back.
				op := CompoundEnvAssignOp{Key: target.Value, Op: node.Op}
				c.compile(node.Value)
				c.derefInNeeded(node.Value)
				c.emit(OpCompoundStoreEnv, c.addConstant(op))
			}
			return
		}
		if node.Op == "=" {
			c.compile(node.Value)
			c.derefInNeeded(node.Value)
			c.emit(OpAssign, index)
		} else {
			// Compound local assignment: opcode reads slot, applies op, stores back.
			op := CompoundAssignOp{VarIndex: index, Op: node.Op}
			c.compile(node.Value)
			c.derefInNeeded(node.Value)
			c.emit(OpCompoundAssign, c.addConstant(op))
		}
	case *ast.MemberExpr:
		// obj.field = val  or  obj[key] = val
		if c.config != nil && c.config.Safe {
			if root := memberRoot(target); root != "" {
				if _, found := c.lookupVariable(root); !found {
					panic(fmt.Sprintf("cannot assign to context variable in safe mode: %s", root))
				}
			}
		}
		if node.Op == "=" {
			// Push [obj, key, val] then store.
			c.compile(target.Node)
			c.derefInNeeded(target.Node)
			c.compile(target.Property)
			c.derefInNeeded(target.Property)
			c.compile(node.Value)
			c.derefInNeeded(node.Value)
			c.emit(OpStoreIndex, 0)
		} else {
			// Compound member assignment: opcode reads current value,
			// applies op, stores back — all in one instruction.
			c.compile(target.Node)
			c.derefInNeeded(target.Node)
			c.compile(target.Property)
			c.derefInNeeded(target.Property)
			c.compile(node.Value)
			c.derefInNeeded(node.Value)
			op := CompoundAssignOp{Op: node.Op}
			c.emit(OpCompoundStoreIndex, c.addConstant(op))
		}
	default:
		panic("assignment target must be an identifier")
	}
}

// IncrDecrNode compiles ++ and --. Postfix form saves the pre-increment
// value by loading it before the opcode, then pops the new value.
func (c *compiler) IncrDecrNode(node *ast.IncDecStmt) {
	ident, ok := node.Target.(*ast.Ident)
	if !ok {
		panic("increment/decrement target must be an identifier")
	}
	index, found := c.lookupVariable(ident.Value)
	if !found {
		panic(fmt.Sprintf("undefined variable: %s", ident.Value))
	}
	if !node.Prefix {
		// Postfix (x++): push the old value so the expression
		// result is the pre-increment value. OpIncrement will push
		// the new value on top; we pop it away below.
		c.emit(OpLoadVar, index)
	}
	if node.Op == "++" {
		c.emit(OpIncrement, index)
	} else {
		c.emit(OpDecrement, index)
	}
	if !node.Prefix {
		// Discard the new value pushed by OpIncrement/OpDecrement;
		// the old value sitting below it is the expression result.
		c.emit(OpPop)
	}
}

// WhileNode emits:
//
//	begin: <cond> → OpJumpIfFalse end → OpPop → <body>
//	       → OpJumpBackward begin
//	end:   (break lands here)
//
// continue in a while loop jumps directly back to begin (no post
// expression), so no continueHoles need patching.
func (c *compiler) WhileNode(node *ast.WhileStmt) {
	begin := len(c.bytecode)
	c.loopStack = append(c.loopStack, loopFrame{isForeach: false, begin: begin})

	c.compile(node.Cond)
	c.derefInNeeded(node.Cond)
	end := c.emit(OpJumpIfFalse, placeholder)
	c.emit(OpPop) // pop the truthy condition value before body runs

	c.compile(node.Body)

	c.emit(OpJumpBackward, c.calcBackwardJump(begin))
	c.patchJump(end) // false-condition exit lands here

	frame := c.loopStack[len(c.loopStack)-1]
	c.loopStack = c.loopStack[:len(c.loopStack)-1]
	for _, ph := range frame.breakHoles {
		c.arguments[ph-1] = len(c.bytecode) - ph
	}
}

// ForeachNode emits:
//
//	<collection> → OpForEach varSlot (initialise iterator)
//	begin: OpJumpIfEnd end → <body>
//	continueTarget: OpIncrementIndex → OpJumpBackward begin
//	end: OpEnd (release iterator)
//	(break lands after OpEnd)
//
// continue jumps to continueTarget so the iterator always advances
// before the next iteration check.
func (c *compiler) ForeachNode(node *ast.ForeachStmt) {
	c.compile(node.Collection)
	c.derefInNeeded(node.Collection)

	index := c.addVariable(node.Var)
	c.emit(OpForEach, index)   // initialise foreach iterator, store first elem
	c.beginScope(node.Var, index)

	begin := len(c.bytecode)
	c.loopStack = append(c.loopStack, loopFrame{isForeach: true, begin: begin})
	end := c.emit(OpJumpIfEnd, placeholder) // exit when iterator exhausted

	c.compile(node.Body)

	// Patch continue holes before emitting OpIncrementIndex so that
	// continue correctly advances the iterator.
	continueTarget := len(c.bytecode)
	frame := c.loopStack[len(c.loopStack)-1]
	c.loopStack = c.loopStack[:len(c.loopStack)-1]
	for _, ph := range frame.continueHoles {
		c.arguments[ph-1] = continueTarget - ph
	}

	c.emit(OpIncrementIndex)
	c.emit(OpJumpBackward, c.calcBackwardJump(begin))
	c.patchJump(end)
	c.endScope()
	c.emit(OpEnd) // release the iterator

	// break exits past OpEnd
	for _, ph := range frame.breakHoles {
		c.arguments[ph-1] = len(c.bytecode) - ph
	}
}

// ForNode emits a C-style for (init; cond; post) loop:
//
//	<init> (once)
//	begin: [<cond> → OpJumpIfFalse end → OpPop]?  <body>
//	continueTarget: [<post> → OpPop]?
//	       → OpJumpBackward begin
//	end:   (break/false-cond land here)
//
// init VarDecls push their value and OpStore it; non-var init
// expressions have their value popped since they are used only
// for side effects. Post expressions are the same.
func (c *compiler) ForNode(node *ast.ForStmt) {
	if node.Init != nil {
		c.compile(node.Init)
		c.derefInNeeded(node.Init)
		// VarDecl stores the value itself via OpStore; plain
		// expressions leave a value we must discard.
		if _, ok := node.Init.(*ast.VarDecl); !ok {
			c.emit(OpPop)
		}
	}

	begin := len(c.bytecode)
	c.loopStack = append(c.loopStack, loopFrame{isForeach: false, hasPostExpr: node.Post != nil, begin: begin})

	var end int
	if node.Cond != nil {
		c.compile(node.Cond)
		c.derefInNeeded(node.Cond)
		end = c.emit(OpJumpIfFalse, placeholder)
		c.emit(OpPop) // pop truthy condition before body
	}

	c.compile(node.Body)

	// continue holes patch here — before post so the post
	// expression runs on every continue, matching C semantics.
	continueTarget := len(c.bytecode)
	frame := c.loopStack[len(c.loopStack)-1]
	c.loopStack = c.loopStack[:len(c.loopStack)-1]
	for _, ph := range frame.continueHoles {
		c.arguments[ph-1] = continueTarget - ph
	}

	if node.Post != nil {
		c.compile(node.Post)
		c.derefInNeeded(node.Post)
		c.emit(OpPop) // post result is unused
	}

	c.emit(OpJumpBackward, c.calcBackwardJump(begin))

	if node.Cond != nil {
		c.patchJump(end)
	}

	for _, ph := range frame.breakHoles {
		c.arguments[ph-1] = len(c.bytecode) - ph
	}
}

// DoWhileNode emits:
//
//	begin: <body>
//	continueTarget: <cond> → OpJumpIfFalse end → OpPop
//	       → OpJumpBackward begin
//	end:   (break/false-cond land here)
//
// continue patches to continueTarget (re-checks condition),
// not back to begin (which would re-run the body first).
func (c *compiler) DoWhileNode(node *ast.DoWhileStmt) {
	begin := len(c.bytecode)
	c.loopStack = append(c.loopStack, loopFrame{isForeach: false, begin: begin})

	c.compile(node.Body)

	// continue holes patch here — before the condition — so that
	// continue re-evaluates the condition, not the body.
	continueTarget := len(c.bytecode)
	frame := c.loopStack[len(c.loopStack)-1]
	c.loopStack = c.loopStack[:len(c.loopStack)-1]
	for _, ph := range frame.continueHoles {
		c.arguments[ph-1] = continueTarget - ph
	}

	c.compile(node.Cond)
	c.derefInNeeded(node.Cond)
	end := c.emit(OpJumpIfFalse, placeholder)
	c.emit(OpPop)

	c.emit(OpJumpBackward, c.calcBackwardJump(begin))
	c.patchJump(end)

	for _, ph := range frame.breakHoles {
		c.arguments[ph-1] = len(c.bytecode) - ph
	}
}

// BreakNode emits OpEnd (for foreach) then a forward OpBreak jump.
// The enclosing loop compiler patches the jump to land after the loop.
func (c *compiler) BreakNode(_ *ast.BreakStmt) {
	if len(c.loopStack) == 0 {
		panic("break outside loop")
	}
	top := len(c.loopStack) - 1
	frame := &c.loopStack[top]
	if frame.isForeach {
		// Foreach's normal exit path emits OpEnd to release the
		// iterator. break skips that, so we emit our own OpEnd here
		// before jumping out.
		c.emit(OpEnd)
	}
	// Emit a forward jump with a placeholder; the loop compiler
	// patches it to land after the loop end.
	ph := c.emit(OpBreak, placeholder)
	frame.breakHoles = append(frame.breakHoles, ph)
}

// ContinueNode jumps to the loop's continue target. For while/do-while
// the target is known (begin), so we emit a backward jump immediately.
// For foreach/for-with-post the target isn't known yet, so we emit a
// placeholder that the loop compiler patches later.
func (c *compiler) ContinueNode(_ *ast.ContinueStmt) {
	if len(c.loopStack) == 0 {
		panic("continue outside loop")
	}
	top := len(c.loopStack) - 1
	frame := &c.loopStack[top]
	if frame.isForeach || frame.hasPostExpr {
		// foreach / for-with-post: continue target is not yet known
		// (it's after the body). Emit a placeholder and patch later.
		ph := c.emit(OpJump, placeholder)
		frame.continueHoles = append(frame.continueHoles, ph)
	} else {
		// while / do-while: continue target is always frame.begin,
		// so we can emit the backward jump immediately.
		c.emit(OpJumpBackward, c.calcBackwardJump(frame.begin))
	}
}

// ReturnNode compiles the return value (if any) and emits OpReturn.
// A bare `return` with no value leaves whatever is on the stack.
func (c *compiler) ReturnNode(node *ast.ReturnStmt) {
	if node.Value != nil {
		c.compile(node.Value)
		c.derefInNeeded(node.Value)
	}
	c.emit(OpReturn)
}

// TryNode compiles try/catch/finally.
//
// try-finally (no catch):
//
//	OpTry catchIP  → <body> → OpTry 0 → OpJump end
//	catchIP: <save thrown val> → <finally> → OpLoadVar → OpThrow
//	end: <finally>
//
// try-catch[-finally]:
//
//	OpTry catchIP  → <body> → OpTry 0 → OpJump end
//	catchIP: [OpStore catchVar] → <catch body> → [<finally> → OpJump past]
//	end: [<finally>]
//
// OpTry arg encodes: catchIP = (ip_after_OpTry) + arg.
// OpTry 0 on the success path pops the try frame without jumping.
func (c *compiler) TryNode(node *ast.TryStmt) {
	if node.CatchBody == nil {
		// try-finally: must guarantee finally runs on both paths,
		// re-throwing the exception on the error path.
		tryInstr := c.emit(OpTry, placeholder)
		c.compile(node.Body)
		c.emit(OpTry, 0) // success: pop try frame
		endJump := c.emit(OpJump, placeholder)
		// Exception path: patch OpTry to land here.
		c.arguments[tryInstr-1] = len(c.bytecode) - tryInstr
		// Stash the thrown value so finally can't accidentally lose it.
		rethrowSlot := c.addVariable("$rethrow")
		c.emit(OpStore, rethrowSlot)
		c.compile(node.FinallyBody)
		c.emit(OpLoadVar, rethrowSlot)
		c.emit(OpThrow)
		// Success path: run finally normally.
		c.patchJump(endJump)
		c.compile(node.FinallyBody)
		return
	}

	tryInstr := c.emit(OpTry, placeholder)
	c.compile(node.Body)
	c.emit(OpTry, 0)                     // success: pop try frame
	endJump := c.emit(OpJump, placeholder) // skip past catch body

	// Patch OpTry to point here — the catch entry.
	// OpThrow pushes the thrown value onto the stack before jumping.
	c.arguments[tryInstr-1] = len(c.bytecode) - tryInstr

	if node.CatchVar != "" {
		// Bind the thrown value to the catch variable.
		varSlot := c.addVariable(node.CatchVar)
		c.emit(OpStore, varSlot)
		c.beginScope(node.CatchVar, varSlot)
		c.compile(node.CatchBody)
		c.endScope()
	} else {
		c.emit(OpPop) // thrown value not needed
		c.compile(node.CatchBody)
	}

	// Compile finally after the catch body (exception path).
	// Jump past the success-path finally copy to avoid running it twice.
	var catchExitJump int
	if node.FinallyBody != nil {
		c.compile(node.FinallyBody)
		catchExitJump = c.emit(OpJump, placeholder)
	}

	// Success path lands here (after catch body).
	c.patchJump(endJump)

	// Compile finally for the success path.
	if node.FinallyBody != nil {
		c.compile(node.FinallyBody)
		c.patchJump(catchExitJump) // catch path jumps past this copy
	}
}

// ThrowNode compiles the thrown value and emits OpThrow, which
// unwinds to the nearest OpTry frame or propagates as an error.
func (c *compiler) ThrowNode(node *ast.ThrowStmt) {
	c.compile(node.Value)
	c.derefInNeeded(node.Value)
	c.emit(OpThrow)
}

// isCaseBodyEmpty returns true when a case body is an empty BlockNode (fallthrough case).
func isCaseBodyEmpty(body ast.Node) bool {
	if block, ok := body.(*ast.BlockStmt); ok {
		return len(block.Statements) == 0
	}
	return false
}

// SwitchNode compiles a switch statement. The subject is evaluated once
// and stashed in a temp slot. Each case emits an equality test chain;
// empty-body cases fall through to the next body. break exits the switch.
func (c *compiler) SwitchNode(node *ast.SwitchStmt) {
	// Compile subject once and stash it in a temp variable slot.
	c.compile(node.Subject)
	c.derefInNeeded(node.Subject)
	tempSlot := c.addVariable("$switch")
	c.emit(OpStore, tempSlot)
	c.beginScope("$switch", tempSlot)

	// Push a loopFrame so that break inside a case body patches to the switch end.
	c.loopStack = append(c.loopStack, loopFrame{isForeach: false})

	// endHoles collects all OpJump placeholders that must land after the switch.
	var endHoles []int
	// fallthroughHoles are jumps from empty-body cases to the next case's body.
	var fallthroughHoles []int

	for _, cas := range node.Cases {
		// Patch any pending fallthrough jumps to land here (at the body of this case,
		// but we need to emit the test first, then body; so patch to body entry below).
		pendingFallthrough := fallthroughHoles
		fallthroughHoles = nil

		// Emit test chain for this case's values.
		var bodyHoles []int
		var skipHoles []int
		for _, val := range cas.Values {
			c.emit(OpLoadVar, tempSlot)
			c.compile(val)
			c.derefInNeeded(val)
			c.emit(OpEqual)
			bodyHoles = append(bodyHoles, c.emit(OpJumpIfTrue, placeholder))
			c.emit(OpPop)
		}
		// None of the values matched: jump to next case block.
		skipHoles = append(skipHoles, c.emit(OpJump, placeholder))

		// Body entry for condition-match path: patch bodyHoles here, emit OpPop.
		bodyEntry := len(c.bytecode)
		for _, ph := range bodyHoles {
			c.arguments[ph-1] = bodyEntry - ph
		}
		c.emit(OpPop) // pop the true comparison result

		// Body entry for fallthrough path: patch pending-fallthrough holes AFTER OpPop.
		bodyEntryAfterPop := len(c.bytecode)
		for _, ph := range pendingFallthrough {
			c.arguments[ph-1] = bodyEntryAfterPop - ph
		}

		if isCaseBodyEmpty(cas.Body) {
			// Fallthrough: emit a forward jump that will land at the NEXT case's body.
			fallthroughHoles = append(fallthroughHoles, c.emit(OpJump, placeholder))
		} else {
			c.compile(cas.Body)
			// Jump to end after executing this body.
			endHoles = append(endHoles, c.emit(OpJump, placeholder))
		}

		// Skip target: patch skipHoles to land here (start of next case tests).
		skipEntry := len(c.bytecode)
		for _, ph := range skipHoles {
			c.arguments[ph-1] = skipEntry - ph
		}
	}

	// Default body (if present).
	// Patch any remaining fallthrough holes to land at the default body.
	defaultEntry := len(c.bytecode)
	for _, ph := range fallthroughHoles {
		c.arguments[ph-1] = defaultEntry - ph
	}
	if node.Default != nil {
		c.compile(node.Default)
	} else {
		// No default: throw an error if no case matched.
		c.emit(OpPush, c.addConstant(fmt.Errorf("switch: no matching case")))
		c.emit(OpThrow)
	}

	// Patch all end-of-body jumps and break holes to land here.
	end := len(c.bytecode)
	for _, ph := range endHoles {
		c.arguments[ph-1] = end - ph
	}
	frame := c.loopStack[len(c.loopStack)-1]
	c.loopStack = c.loopStack[:len(c.loopStack)-1]
	for _, ph := range frame.breakHoles {
		c.arguments[ph-1] = end - ph
	}

	c.endScope()
}

// EmptyNode compiles the operand and emits OpEmpty, which tests
// whether the value is nil, "", 0, or an empty collection.
func (c *compiler) EmptyNode(node *ast.EmptyExpr) {
	c.compile(node.Value)
	c.derefInNeeded(node.Value)
	c.emit(OpEmpty)
}

// SizeNode compiles the operand and emits OpSize, which returns
// the length of a string, array, or map as int64.
func (c *compiler) SizeNode(node *ast.SizeExpr) {
	c.compile(node.Value)
	c.derefInNeeded(node.Value)
	c.emit(OpSize)
}

// SetNode compiles each element onto the stack and emits OpSet
// with the element count so the VM builds a HashSet.
func (c *compiler) SetNode(node *ast.SetLit) {
	for _, elem := range node.Elements {
		c.compile(elem)
		c.derefInNeeded(elem)
	}
	c.emit(OpSet, len(node.Elements))
}

// NewNode pushes the class name then each constructor arg, then emits
// OpNew. The VM looks up the class in the registry and calls its constructor.
func (c *compiler) NewNode(node *ast.NewExpr) {
	c.emitPush(node.ClassName)
	for _, arg := range node.Args {
		c.compile(arg)
		c.derefInNeeded(arg)
	}
	c.emit(OpNew, len(node.Args))
}

// NamespaceCallNode compiles a Namespace::method(...) call. The
// registry object is captured into a closure so the VM can dispatch
// via the standard function table without knowing about namespaces.
func (c *compiler) NamespaceCallNode(node *ast.NamespaceCallExpr) {
	ns := node.Namespace
	method := node.Method
	obj, _ := c.config.Registry.Lookup(ns)

	for _, arg := range node.Args {
		c.compile(arg)
		c.derefInNeeded(arg)
	}

	rfn := &runtime.Function{
		Name: ns + ":" + method,
		Func: func(args ...any) (any, error) { return obj.Call(method, args...) },
	}
	c.emitFunction(rfn, len(node.Args))
}

// derefInNeeded emits OpDeref when the node's type is a non-nil
// pointer (other than *big.Int, which the VM handles specially).
// This unwraps *T so subsequent operations work on T directly.
func (c *compiler) derefInNeeded(node ast.Node) {
	if node.Nature().Nil {
		return
	}
	t := node.Type()
	if t.Kind() == reflect.Ptr && t != bigIntPtrType {
		c.emit(OpDeref)
	}
}

// derefParam emits OpDeref when passing a *T argument to a func
// parameter that expects T (not *T). Skipped when the pointer type
// is already directly assignable to the param type.
func (c *compiler) derefParam(in reflect.Type, param ast.Node) {
	if param.Nature().Nil {
		return
	}
	if param.Type().AssignableTo(in) {
		return
	}
	if in.Kind() != reflect.Ptr && param.Type().Kind() == reflect.Ptr {
		c.emit(OpDeref)
	}
}

// optimize performs a single-pass peephole optimisation: it
// collapses chains of identical conditional jumps (e.g. two
// consecutive OpJumpIfTrue) into one jump that skips them all.
// This is safe because a chain of the same conditional opcode
// has the same condition, so all but the first are dead.
func (c *compiler) optimize() {
	for i, op := range c.bytecode {
		switch op {
		case OpJumpIfTrue, OpJumpIfFalse, OpJumpIfNil, OpJumpIfNotNil:
			target := i + c.arguments[i] + 1
			for target < len(c.bytecode) && c.bytecode[target] == op {
				target += c.arguments[target] + 1
			}
			c.arguments[i] = target - i - 1
		}
	}
}

// kind returns t.Kind(), treating nil as reflect.Invalid.
func kind(t reflect.Type) reflect.Kind {
	if t == nil {
		return reflect.Invalid
	}
	return t.Kind()
}

// extractClassCall peels a.b.c.Method() into ("a.b.c", "Method", true).
func extractClassCall(callee ast.Node) (className, method string, ok bool) {
	member, ok := callee.(*ast.MemberExpr)
	if !ok {
		return "", "", false
	}
	prop, ok := member.Property.(*ast.StringLit)
	if !ok {
		return "", "", false
	}
	method = prop.Value
	className, ok = dotChain(member.Node)
	return className, method, ok
}

// memberRoot walks a MemberExpr chain and returns the root identifier name.
func memberRoot(node ast.Node) string {
	switch n := node.(type) {
	case *ast.Ident:
		return n.Value
	case *ast.MemberExpr:
		return memberRoot(n.Node)
	default:
		return ""
	}
}

// dotChain flattens a chain of MemberExprs into a dotted string
// (e.g. "java.util.Arrays"), returning false on any non-string property.
func dotChain(node ast.Node) (string, bool) {
	switch n := node.(type) {
	case *ast.Ident:
		return n.Value, true
	case *ast.MemberExpr:
		prop, ok := n.Property.(*ast.StringLit)
		if !ok {
			return "", false
		}
		left, ok := dotChain(n.Node)
		if !ok {
			return "", false
		}
		return left + "." + prop.Value, true
	default:
		return "", false
	}
}
