// Copyright (c) 2026 Harness Inc.
// Copyright (c) 2018 Anton Medvedev
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package checker

import (
	"fmt"
	"reflect"
	"regexp"
	"time"

	"math/big"

	"github.com/harness/go-jexl/jexl/ast"
	. "github.com/harness/go-jexl/jexl/checker/nature"
	"github.com/harness/go-jexl/jexl/classes/java/util"
	conf "github.com/harness/go-jexl/jexl/config"
	"github.com/harness/go-jexl/jexl/internal/decimal"
	"github.com/harness/go-jexl/jexl/parser"
	"github.com/harness/go-jexl/jexl/token"
	"github.com/harness/go-jexl/jexl/vm/runtime"
)

var (
	anyType        = reflect.TypeOf(new(any)).Elem()
	boolType       = reflect.TypeOf(true)
	intType        = reflect.TypeOf(0)
	int64Type      = reflect.TypeOf(int64(0))
	floatType      = reflect.TypeOf(float64(0))
	stringType     = reflect.TypeOf("")
	arrayType      = reflect.TypeOf([]any{})
	mapType        = reflect.TypeOf(map[string]any{})
	anyMapType     = reflect.TypeOf(map[any]any{})
	setType        = reflect.TypeOf(&util.HashSet{})
	setValueType   = reflect.TypeOf(util.HashSet{})
	regexpType     = reflect.TypeOf((*regexp.Regexp)(nil))
	regexpStructT  = reflect.TypeOf(regexp.Regexp{})
	timeType       = reflect.TypeOf(time.Time{})
	durationType   = reflect.TypeOf(time.Duration(0))
	javaBigIntType = reflect.TypeOf(big.Int{})
	javaBigDecType = reflect.TypeOf(decimal.Decimal{})

	anyTypeSlice = []reflect.Type{anyType}
)

// ParseCheck parses input expression and checks its types.
// In case of error, it returns error with a tree.
func ParseCheck(input string, config *conf.Config) (*ast.Tree, error) {
	tree, err := parser.ParseWithConfig(input, config)

	if err != nil {
		return tree, err
	}

	_, err = new(Checker).Check(tree, config)
	if err != nil {
		return tree, err
	}

	return tree, nil
}

// Checker walks an AST and assigns a Nature to every node,
// recording the first type error it encounters.
type Checker struct {
	config          *conf.Config
	predicateScopes []predicateScope
	varScopes       []varScope // lexical variable stack (innermost last)
	loopDepth       int        // nesting level; > 0 inside a loop
	switchDepth     int        // nesting level; > 0 inside a switch
	err             *token.Error
	needsReset      bool
}

// predicateScope holds the collection and loop-var bindings
// introduced by a predicate / filter expression.
type predicateScope struct {
	collection Nature
	vars       []varScope
}

// varScope is one entry on the variable stack — a name bound
// to a Nature inside a let / lambda / for / catch block.
type varScope struct {
	name   string
	nature Nature
}

// Check checks types of the expression tree. It returns type of the expression
// and error if any. If config is nil, then default configuration will be used.
func (v *Checker) Check(tree *ast.Tree, config *conf.Config) (reflect.Type, error) {
	v.reset(config)
	return v.check(tree)
}

// check visits the root node, normalizes nil type to any,
// and binds the source to any recorded error before returning.
func (v *Checker) check(tree *ast.Tree) (reflect.Type, error) {
	nt := v.visit(tree.Node)

	// nil Type means unknown; expose as any for callers
	// that predate the Nature type.
	t := nt.Type
	if t == nil {
		t = anyType
	}

	if v.err != nil {
		return t, v.err.Bind(tree.Source)
	}

	return t, nil
}

// reset prepares the Checker for a new Check call,
// clearing scopes and errors from any prior run.
func (v *Checker) reset(config *conf.Config) {
	if v.needsReset {
		// Zero out slice elements before truncating so that
		// Nature values (which hold reflect.Type pointers) are
		// not retained and can be GC'd.
		clearSlice(v.predicateScopes)
		clearSlice(v.varScopes)
		v.predicateScopes = v.predicateScopes[:0]
		v.varScopes = v.varScopes[:0]
		v.loopDepth = 0
		v.err = nil
	}
	v.needsReset = true

	if config == nil {
		config = conf.New()
	}
	v.config = config
}

// clearSlice zeros each element so GC can reclaim any
// pointer-holding values (e.g. Nature.Type) before reslicing.
func clearSlice[S ~[]E, E any](s S) {
	var zero E
	for i := range s {
		s[i] = zero
	}
}

// visit dispatches to the correct typed handler, sets the
// resulting Nature on the node, and returns it.
func (v *Checker) visit(node ast.Node) Nature {
	var nt Nature
	switch n := node.(type) {
	case *ast.NilLit:
		nt = v.config.Cache.NatureOf(nil)
	case *ast.Ident:
		nt = v.identifierNode(n)
	case *ast.IntLit:
		nt = v.config.Cache.FromType(intType)
	case *ast.FloatLit:
		nt = v.config.Cache.FromType(floatType)
	case *ast.BoolLit:
		nt = v.config.Cache.FromType(boolType)
	case *ast.StringLit:
		nt = v.config.Cache.FromType(stringType)
	case *ast.ConstantExpr:
		nt = v.config.Cache.FromType(reflect.TypeOf(n.Value))
	case *ast.UnaryExpr:
		nt = v.unaryNode(n)
	case *ast.BinaryExpr:
		nt = v.binaryNode(n)
	case *ast.ChainExpr:
		nt = v.chainNode(n)
	case *ast.MemberExpr:
		nt = v.memberNode(n)
	case *ast.CallExpr:
		nt = v.callNode(n)
	case *ast.PredicateExpr:
		nt = v.predicateNode(n)
	case *ast.ConditionalExpr:
		nt = v.conditionalNode(n)
	case *ast.ArrayLit:
		nt = v.arrayNode(n)
	case *ast.MapLit:
		nt = v.mapNode(n)
	case *ast.SetLit:
		nt = v.setNode(n)
	case *ast.KeyValueExpr:
		nt = v.pairNode(n)
	case *ast.BitExpr:
		nt = v.bitOpNode(n)
	case *ast.StrictEqualExpr:
		nt = v.strictEqualNode(n)
	case *ast.LambdaExpr:
		nt = v.lambdaNode(n)
	case *ast.BlockStmt:
		nt = v.blockNode(n)
	case *ast.VarDecl:
		nt = v.varNode(n)
	case *ast.AssignStmt:
		nt = v.assignNode(n)
	case *ast.IncDecStmt:
		nt = v.incrDecrNode(n)
	case *ast.ForeachStmt:
		nt = v.foreachNode(n)
	case *ast.WhileStmt:
		nt = v.whileNode(n)
	case *ast.ForStmt:
		nt = v.forNode(n)
	case *ast.DoWhileStmt:
		nt = v.doWhileNode(n)
	case *ast.BreakStmt:
		nt = v.breakNode(n)
	case *ast.ContinueStmt:
		nt = v.continueNode(n)
	case *ast.ReturnStmt:
		nt = v.returnNode(n)
	case *ast.TryStmt:
		nt = v.tryNode(n)
	case *ast.ThrowStmt:
		nt = v.throwNode(n)
	case *ast.SwitchStmt:
		nt = v.switchNode(n)
	case *ast.EmptyExpr:
		nt = v.emptyNode(n)
	case *ast.SizeExpr:
		nt = v.sizeNode(n)
	case *ast.RegexLit:
		nt = v.regexNode(n)
	case *ast.NewExpr:
		nt = v.newNode(n)
	case *ast.NamespaceCallExpr:
		nt = v.namespaceCallNode(n)
	default:
		panic(fmt.Sprintf("undefined node type (%T)", node))
	}
	node.SetNature(nt)
	return nt
}

// isBigNum returns true if n is a BigInteger or BigDecimal type (pointer or struct).
func isBigNum(n Nature) bool {
	if n.Type == javaBigIntType || n.Type == javaBigDecType {
		return true
	}
	if n.Kind == reflect.Ptr && n.Type != nil {
		return n.Type.Elem() == javaBigIntType
	}
	return false
}

// error records the first type error encountered and returns
// an empty Nature so callers can propagate without branching.
func (v *Checker) error(node ast.Node, format string, args ...any) Nature {
	if v.err == nil { // show first error
		v.err = &token.Error{
			Range:   node.Location(),
			Message: fmt.Sprintf(format, args...),
		}
	}
	return Nature{}
}

// identifierNode resolves a bare name: first checks var
// scopes (innermost first), then the environment/functions.
func (v *Checker) identifierNode(node *ast.Ident) Nature {
	for i := len(v.varScopes) - 1; i >= 0; i-- {
		if v.varScopes[i].name == node.Value {
			return v.varScopes[i].nature
		}
	}
	if node.Value == "$env" {
		return Nature{}
	}

	return v.ident(node, node.Value, v.config.Strict, true)
}

// ident resolves name against the env context first, then
// (when builtins=true) against declared functions.
// strict=true turns an unresolved name into a compile error.
// memberNode passes builtins=false so $env.foo never resolves
// to a global function.
func (v *Checker) ident(node ast.Node, name string, strict, builtins bool) Nature {
	if nt, ok := v.config.Context.Get(&v.config.Cache, name); ok {
		return nt
	}
	if builtins {
		if fn, ok := v.config.Functions[name]; ok {
			nt := v.config.Cache.FromType(fn.Type())
			// TypeData carries the *runtime.Function so callNode
			// can later route to checkFunction instead of the
			// generic checkArguments path.
			if nt.TypeData == nil {
				nt.TypeData = new(TypeData)
			}
			nt.TypeData.Func = fn
			return nt
		}
	}
	if v.config.Strict && strict {
		return v.error(node, "unknown name %s", name)
	}
	// Unknown in non-strict mode — treat as any.
	return Nature{}
}

// unaryNode checks !, not, +, - operators. Pointer operands
// are unwrapped first; BigNum and unknown pass through to runtime.
func (v *Checker) unaryNode(node *ast.UnaryExpr) Nature {
	nt := v.visit(node.Node)
	nt = nt.Deref(&v.config.Cache) // unwrap *T before type tests

	switch node.Operator {

	case "!", "not":
		// Both bool and unknown operands yield bool.
		// Unknown is allowed so that `!anyVar` compiles
		// in non-strict mode.
		if nt.IsBool() {
			return v.config.Cache.FromType(boolType)
		}
		if nt.IsUnknown(&v.config.Cache) {
			return v.config.Cache.FromType(boolType)
		}

	case "+", "-":
		if nt.IsNumber() {
			return nt // preserve int vs float
		}
		// JEXL: -bool = !bool (for -)  or stays bool (for +)
		if nt.IsBool() {
			return v.config.Cache.FromType(boolType)
		}
		// unknown or BigInt/BigDec: pass through, runtime handles it
		if nt.IsUnknown(&v.config.Cache) {
			return Nature{}
		}
		if isBigNum(nt) {
			return Nature{}
		}

	default:
		return v.error(node, "unknown operator (%s)", node.Operator)
	}

	return v.error(node, `invalid operation: %s (mismatched type %s)`, node.Operator, nt.String())
}

// binaryNode checks all infix operators. Each case returns the
// result type; falling through to the bottom emits a type error.
func (v *Checker) binaryNode(node *ast.BinaryExpr) Nature {
	// instanceof only needs the LHS visited; the RHS is a
	// type name, not an expression.
	if node.Operator == "instanceof" || node.Operator == "!instanceof" {
		v.visit(node.Left)
		return v.config.Cache.FromType(boolType)
	}

	l := v.visit(node.Left)
	r := v.visit(node.Right)

	// Unwrap pointers before type tests below.
	l = l.Deref(&v.config.Cache)
	r = r.Deref(&v.config.Cache)

	switch node.Operator {
	case "==", "!=":
		if l.ComparableTo(&v.config.Cache, r) {
			return v.config.Cache.FromType(boolType)
		}
		// JEXL loose equality: number == string coerces via toString.
		if (l.IsNumber() && r.IsString()) || (l.IsString() && r.IsNumber()) {
			return v.config.Cache.FromType(boolType)
		}

	case "or", "||", "and", "&&":
		// Accept mixed bool/unknown — JEXL is truthy, so
		// `x && y` is valid even when x is not a bool.
		if l.IsBool() && r.IsBool() {
			return v.config.Cache.FromType(boolType)
		}
		if l.IsBool() || r.IsBool() {
			return v.config.Cache.FromType(boolType)
		}
		if l.MaybeCompatible(&v.config.Cache, r, BoolCheck) {
			return v.config.Cache.FromType(boolType)
		}

	case "<", ">", ">=", "<=":
		if l.IsNumber() && r.IsNumber() {
			return v.config.Cache.FromType(boolType)
		}
		if l.IsString() && r.IsString() {
			return v.config.Cache.FromType(boolType)
		}
		if l.IsTime() && r.IsTime() {
			return v.config.Cache.FromType(boolType)
		}
		if l.IsDuration() && r.IsDuration() {
			return v.config.Cache.FromType(boolType)
		}
		// BigInt/BigDec comparisons are handled at runtime.
		if isBigNum(l) || isBigNum(r) {
			return v.config.Cache.FromType(boolType)
		}
		if l.MaybeCompatible(&v.config.Cache, r, NumberCheck, StringCheck, TimeCheck, DurationCheck) {
			return v.config.Cache.FromType(boolType)
		}

	case "-":
		if l.IsNumber() && r.IsNumber() {
			return l.PromoteNumericNature(&v.config.Cache, r)
		}
		// time - time = duration; time - duration = time
		if l.IsTime() && r.IsTime() {
			return v.config.Cache.FromType(durationType)
		}
		if l.IsTime() && r.IsDuration() {
			return v.config.Cache.FromType(timeType)
		}
		if l.IsDuration() && r.IsDuration() {
			return v.config.Cache.FromType(durationType)
		}
		// JEXL coerces bool/string operands to numbers;
		// result type is unknown until runtime.
		if (l.IsBool() || l.IsString()) && (r.IsNumber() || r.IsBool() || r.IsString()) {
			return Nature{}
		}
		if l.IsNumber() && (r.IsBool() || r.IsString()) {
			return Nature{}
		}
		if isBigNum(l) || isBigNum(r) {
			return Nature{}
		}
		if l.MaybeCompatible(&v.config.Cache, r, NumberCheck, TimeCheck, DurationCheck) {
			return Nature{}
		}

	case "*":
		if l.IsNumber() && r.IsNumber() {
			return l.PromoteNumericNature(&v.config.Cache, r)
		}
		// Scaling a duration by a number yields a duration.
		if l.IsNumber() && r.IsDuration() {
			return v.config.Cache.FromType(durationType)
		}
		if l.IsDuration() && r.IsNumber() {
			return v.config.Cache.FromType(durationType)
		}
		if l.IsDuration() && r.IsDuration() {
			return v.config.Cache.FromType(durationType)
		}
		// JEXL coerces bool/string to number for multiplication.
		if (l.IsBool() || l.IsString()) && (r.IsNumber() || r.IsBool() || r.IsString()) {
			return Nature{}
		}
		if l.IsNumber() && (r.IsBool() || r.IsString()) {
			return Nature{}
		}
		if isBigNum(l) || isBigNum(r) {
			return Nature{}
		}
		if l.MaybeCompatible(&v.config.Cache, r, NumberCheck, DurationCheck) {
			return Nature{}
		}

	case "/":
		// Division always produces float (like Java / JEXL).
		if l.IsNumber() && r.IsNumber() {
			return v.config.Cache.FromType(floatType)
		}
		if isBigNum(l) || isBigNum(r) {
			return Nature{}
		}
		if l.MaybeCompatible(&v.config.Cache, r, NumberCheck) {
			return v.config.Cache.FromType(floatType)
		}

	case "**", "^":
		// Exponentiation always produces float.
		if l.IsNumber() && r.IsNumber() {
			return v.config.Cache.FromType(floatType)
		}
		if isBigNum(l) || isBigNum(r) {
			return Nature{}
		}
		if l.MaybeCompatible(&v.config.Cache, r, NumberCheck) {
			return v.config.Cache.FromType(floatType)
		}

	case "%":
		// Modulo produces int (truncated).
		if l.IsNumber() && r.IsNumber() {
			return v.config.Cache.FromType(intType)
		}
		if isBigNum(l) || isBigNum(r) {
			return Nature{}
		}
		if l.MaybeCompatible(&v.config.Cache, r, NumberCheck) {
			return v.config.Cache.FromType(intType)
		}

	case "+":
		if l.IsNumber() && r.IsNumber() {
			return l.PromoteNumericNature(&v.config.Cache, r)
		}
		if l.IsString() && r.IsString() {
			return v.config.Cache.FromType(stringType)
		}
		// Any operand being a string forces string concatenation.
		if l.IsString() || r.IsString() {
			return v.config.Cache.FromType(stringType)
		}
		// JEXL coerces bool/string to number for addition.
		if (l.IsBool() || l.IsString()) && (r.IsNumber() || r.IsBool() || r.IsString()) {
			return Nature{}
		}
		if l.IsNumber() && (r.IsBool() || r.IsString()) {
			return Nature{}
		}
		// time + duration = time; duration + duration = duration
		if l.IsTime() && r.IsDuration() {
			return v.config.Cache.FromType(timeType)
		}
		if l.IsDuration() && r.IsTime() {
			return v.config.Cache.FromType(timeType)
		}
		if l.IsDuration() && r.IsDuration() {
			return v.config.Cache.FromType(durationType)
		}
		if isBigNum(l) || isBigNum(r) {
			return Nature{}
		}
		if l.MaybeCompatible(&v.config.Cache, r, NumberCheck, StringCheck, TimeCheck, DurationCheck) {
			return Nature{}
		}

	case "in", "=~", "!~":
		// Membership test against a Set is always valid.
		if r.Type == setType || r.Type == setValueType {
			return v.config.Cache.FromType(boolType)
		}
		// =~ / !~ with a string or regexp RHS is a regex match.
		if node.Operator != "in" && (l.IsString() || l.IsByteSlice() || l.IsUnknown(&v.config.Cache)) && (r.IsString() || r.Type == regexpType || r.IsUnknown(&v.config.Cache)) {
			return v.config.Cache.FromType(boolType)
		}
		// string `in` struct delegates to a Contains-style method.
		if (l.IsString() || l.IsUnknown(&v.config.Cache)) && r.IsStruct() {
			return v.config.Cache.FromType(boolType)
		}
		if r.IsMap() {
			// Verify the key type is assignable to the map's key type.
			rKey := r.Key(&v.config.Cache)
			if !l.IsUnknown(&v.config.Cache) && !l.AssignableTo(rKey) {
				return v.error(node, "cannot use %s as type %s in map key", l.String(), rKey.String())
			}
			return v.config.Cache.FromType(boolType)
		}
		if r.IsArray() {
			// Verify the element is comparable to the array's elem type.
			rElem := r.Elem(&v.config.Cache)
			if !l.ComparableTo(&v.config.Cache, rElem) {
				return v.error(node, "cannot use %s as type %s in array", l.String(), rElem.String())
			}
			return v.config.Cache.FromType(boolType)
		}
		// Unknown LHS against a known collection kind is fine.
		if l.IsUnknown(&v.config.Cache) && r.IsAnyOf(StringCheck, ArrayCheck, MapCheck) {
			return v.config.Cache.FromType(boolType)
		}
		// Unknown RHS — can't rule it out at compile time.
		if r.IsUnknown(&v.config.Cache) {
			return v.config.Cache.FromType(boolType)
		}

	case "=^", "!^", "=$", "!$":
		// startsWith / endsWith operators — both sides must be strings.
		if l.IsString() && r.IsString() {
			return v.config.Cache.FromType(boolType)
		}
		if l.MaybeCompatible(&v.config.Cache, r, StringCheck) {
			return v.config.Cache.FromType(boolType)
		}

	case "..":
		// Range operator — produces []int.
		if l.IsInteger && r.IsInteger || l.MaybeCompatible(&v.config.Cache, r, IntegerCheck) {
			return ArrayFromType(&v.config.Cache, intType)
		}

	case "??":
		// Null-coalescing: return whichever side is non-nil.
		// If both sides have compatible types, use the LHS type.
		if l.Nil && !r.Nil {
			return r
		}
		if !l.Nil && r.Nil {
			return l
		}
		if l.Nil && r.Nil {
			return v.config.Cache.NatureOf(nil)
		}
		if r.AssignableTo(l) {
			return l
		}
		return Nature{} // incompatible types — unknown at compile time

	case "instanceof", "!instanceof":
		return v.config.Cache.FromType(boolType)

	default:
		return v.error(node, "unknown operator (%s)", node.Operator)

	}

	// If we reach here every typed branch above fell through.
	// BigNum operands are a last-resort escape hatch — runtime
	// handles the actual arithmetic via reflection.
	if isBigNum(l) || isBigNum(r) {
		return Nature{}
	}

	return v.error(node, `invalid operation: %s (mismatched types %s and %s)`, node.Operator, l.String(), r.String())
}

// chainNode is a transparent wrapper — it delegates straight to
// the inner node. The ChainExpr exists only to carry optional-chain
// semantics in the AST; the checker doesn't need to act on it.
func (v *Checker) chainNode(node *ast.ChainExpr) Nature {
	return v.visit(node.Node)
}

// memberNode resolves field, index, and method access.
// It handles $env shortcuts, optional chaining (?.), pointer
// unwrapping, and Java-compat method allow-lists per kind.
func (v *Checker) memberNode(node *ast.MemberExpr) Nature {
	// $env.foo is a direct environment lookup — skip the normal
	// base/property visit and go straight to ident resolution.
	// builtins=false so $env.size doesn't accidentally bind to
	// a global function named "size".
	if an, ok := node.Node.(*ast.Ident); ok && an.Value == "$env" {
		if name, ok := node.Property.(*ast.StringLit); ok {
			strict := v.config.Strict
			// Optional flag ($env?.foo) suppresses "unknown name"
			// errors — user is intentionally guarding the access.
			if node.Optional {
				strict = false
			}
			return v.ident(node, name.Value, strict, false)
		}
		return Nature{}
	}

	// Snapshot the current error so that if visiting the base
	// produces a new error and the access is optional (?.), we
	// can roll back and return nil instead of propagating
	// the error.
	savedErr := v.err
	base := v.visit(node.Node)
	if node.Optional && v.err != nil && v.err != savedErr {
		v.err = savedErr
		return v.config.Cache.NatureOf(nil)
	}
	prop := v.visit(node.Property)

	// Unknown base means we can't validate the access; pass through.
	if base.IsUnknown(&v.config.Cache) {
		return Nature{}
	}

	if name, ok := node.Property.(*ast.StringLit); ok {
		if base.Nil {
			if node.Optional {
				return v.config.Cache.NatureOf(nil)
			}
			// Strict mode: nil dereference is a compile error.
			// Non-strict: let it blow up at runtime (try/catch can handle it).
			if v.config.Strict {
				return v.error(node, "type nil has no field %s", name.Value)
			}
			return Nature{}
		}

		// Check methods on the *original* (non-dereferenced) type first.
		// A pointer receiver method must be found before we unwrap.
		if m, ok := base.MethodByName(&v.config.Cache, name.Value); ok {
			return m
		}
	}

	// Now unwrap one pointer level for field/index access.
	base = base.Deref(&v.config.Cache)

	if base.Nil && node.Optional {
		return v.config.Cache.NatureOf(nil)
	}

	switch base.Kind {
	case reflect.Map:
		// Java-style map method calls — return unknown because
		// the exact return type depends on the call arguments
		// which we don't track here.
		if name, ok := node.Property.(*ast.StringLit); ok && node.Method {
			switch name.Value {
			case "size", "isEmpty", "containsKey", "containsValue", "get", "put",
				"remove", "keySet", "values", "entrySet", "clear", "putAll",
				"putIfAbsent", "getOrDefault", "forEach", "merge":
				return Nature{}
			}
		}
		// Try to match prop to the map's key type. Deref the
		// property if necessary (e.g. *string key on string map).
		if !prop.AssignableTo(base.Key(&v.config.Cache)) {
			propDeref := prop.Deref(&v.config.Cache)
			if propDeref.AssignableTo(base.Key(&v.config.Cache)) {
				prop = propDeref
			}
		}
		if !prop.AssignableTo(base.Key(&v.config.Cache)) && !prop.IsUnknown(&v.config.Cache) {
			return v.error(node.Property, "cannot use %s to get an element from %s", prop.String(), base.String())
		}
		// For typed/strict maps, verify the string key exists.
		if prop, ok := node.Property.(*ast.StringLit); ok && base.TypeData != nil {
			if field, ok := base.Fields[prop.Value]; ok {
				return field
			} else if base.Strict {
				return v.error(node.Property, "unknown field %s", prop.Value)
			}
		}
		return base.Elem(&v.config.Cache)

	case reflect.Array, reflect.Slice:
		// Java-style collection method calls.
		if name, ok := node.Property.(*ast.StringLit); ok && node.Method {
			switch name.Value {
			case "size", "length", "isEmpty", "contains", "get", "add", "remove",
				"indexOf", "lastIndexOf", "subList", "toArray", "iterator",
				"stream", "forEach", "sort", "reverse", "clear", "set",
				"addAll", "removeAll", "containsAll", "retainAll":
				return Nature{}
			}
		}
		prop = prop.Deref(&v.config.Cache)
		// Only numeric subscripts are valid on arrays/slices.
		if !prop.IsInteger && !prop.IsFloat && !prop.IsUnknown(&v.config.Cache) {
			return v.error(node.Property, "array elements can only be selected using an integer (got %s)", prop.String())
		}
		return base.Elem(&v.config.Cache)

	case reflect.Struct:
		if name, ok := node.Property.(*ast.StringLit); ok {
			propertyName := name.Value
			if field, ok := base.FieldByName(&v.config.Cache, propertyName); ok {
				return v.config.Cache.FromType(field.Type)
			}
			if node.Method {
				return v.error(node, "type %v has no method %v", base.String(), propertyName)
			}
			return v.error(node, "type %v has no field %v", base.String(), propertyName)
		}
	}

	// Fallback: allow well-known Java primitive / string methods on
	// any type that reached here. Return unknown — the VM resolves
	// these via reflection at runtime.
	if name, ok := node.Property.(*ast.StringLit); ok && node.Method {
		switch name.Value {
		case "booleanValue", "byteValue", "charAt", "codePointAt",
			"compareTo", "compareToIgnoreCase", "concat", "contains",
			"doubleValue", "endsWith", "equals", "equalsIgnoreCase",
			"floatValue", "formatted", "indexOf", "isBlank",
			"isEmpty", "isInfinite", "isNaN", "intValue",
			"lastIndexOf", "length", "longValue", "matches",
			"repeat", "replace", "replaceAll", "replaceFirst",
			"shortValue", "split", "startsWith", "strip",
			"stripLeading", "stripTrailing", "substring", "toCharArray",
			"toLowerCase", "toString", "toUpperCase", "trim",
			"trimLeft", "trimRight", "valueOf":
			return Nature{}
		}
	}

	if name, ok := node.Property.(*ast.StringLit); ok {
		if node.Method {
			return v.error(node, "type %v has no method %v", base.String(), name.Value)
		}
		return v.error(node, "type %v has no field %v", base.String(), name.Value)
	}
	return v.error(node, "type %v[%v] is undefined", base.String(), prop.String())
}

// callNode type-checks a function call. It handles cached types,
// registry class calls, declared functions with overloads, and
// plain reflect.Func callees via checkArguments.
func (v *Checker) callNode(node *ast.CallExpr) Nature {
	// A patcher may have already stamped a concrete type onto
	// the node during a prior pass. Re-use it to avoid
	// overwriting correct type info with anyType.
	// We skip this shortcut when the type is anyType itself —
	// that signals "not yet resolved" from an earlier error
	// path, and a second pass may have better information.
	if typ := node.Type(); typ != nil && typ != anyType {
		return *node.Nature()
	}

	// $env is the environment object — it's not a function.
	if id, ok := node.Callee.(*ast.Ident); ok && id.Value == "$env" {
		return v.error(node, "%s is not callable", v.config.Context.String())
	}

	// Registry-backed class calls (e.g. MyClass.staticMethod())
	// are opaque to the checker — just type-check the args and
	// return unknown; the registry validates at runtime.
	if v.config.Registry != nil {
		if className, _, ok := extractClassCall(node.Callee); ok {
			if _, found := v.config.Registry.LookupClass(className); found {
				for _, arg := range node.Arguments {
					v.visit(arg)
				}
				return Nature{}
			}
		}
	}

	nt := v.visit(node.Callee)
	// Unknown callee (e.g. non-strict unknown var) — pass through.
	if nt.IsUnknown(&v.config.Cache) {
		return Nature{}
	}

	// Declared functions carry a *runtime.Function with possible
	// overloads and a custom Validate hook — use the richer path.
	if nt.TypeData != nil && nt.TypeData.Func != nil {
		return v.checkFunction(nt.TypeData.Func, node, node.Arguments)
	}

	// Best-effort name for error messages.
	fnName := "function"
	if identifier, ok := node.Callee.(*ast.Ident); ok {
		fnName = identifier.Value
	}
	if member, ok := node.Callee.(*ast.MemberExpr); ok {
		if name, ok := member.Property.(*ast.StringLit); ok {
			fnName = name.Value
		}
	}

	if nt.Nil {
		// foo?.bar() — safe-nav call on nil returns nil.
		if member, ok := node.Callee.(*ast.MemberExpr); ok && member.Optional {
			return v.config.Cache.NatureOf(nil)
		}
		// Non-strict: let the nil call fail at runtime.
		if !v.config.Strict {
			return Nature{}
		}
		return v.error(node, "%v is nil; cannot call nil as function", fnName)
	}

	if nt.Kind == reflect.Func {
		outType, err := v.checkArguments(fnName, nt, node.Arguments, node)
		if err != nil {
			if v.err == nil {
				v.err = err
			}
			return Nature{}
		}
		return outType
	}
	return v.error(node, "%s is not callable", nt.String())
}

// checkFunction resolves the correct overload of f, using
// f.Validate when present, otherwise trying each type in f.Types.
func (v *Checker) checkFunction(f *runtime.Function, node ast.Node, arguments []ast.Node) Nature {
	if f.Validate != nil {
		// Custom validator — collect arg types (substituting any
		// for unknowns) and let the function decide its return type.
		args := make([]reflect.Type, len(arguments))
		for i, arg := range arguments {
			argNature := v.visit(arg)
			if argNature.IsUnknown(&v.config.Cache) {
				args[i] = anyType
			} else {
				args[i] = argNature.Type
			}
		}
		t, err := f.Validate(args)
		if err != nil {
			return v.error(node, "%v", err)
		}
		return v.config.Cache.FromType(t)
	} else if len(f.Types) == 0 {
		// No overloads declared — validate against the single
		// reflect.Type returned by f.Type().
		nt, err := v.checkArguments(f.Name, v.config.Cache.FromType(f.Type()), arguments, node)
		if err != nil {
			if v.err == nil {
				v.err = err
			}
			return Nature{}
		}
		return nt
	}
	// Multiple overloads: try each in order and use the first
	// that accepts the given arguments.
	var lastErr *token.Error
	for _, t := range f.Types {
		outNature, err := v.checkArguments(f.Name, v.config.Cache.FromType(t), arguments, node)
		if err != nil {
			lastErr = err
			continue
		}

		// Stamp the matched overload type onto the callee node so
		// the compiler can emit the correct OpDeref when needed.
		if callNode, ok := node.(*ast.CallExpr); ok {
			callNode.Callee.SetType(t)
		}

		return outNature
	}
	if lastErr != nil {
		if v.err == nil {
			v.err = lastErr
		}
		return Nature{}
	}

	return v.error(node, "no matching overload for %v", f.Name)
}

// checkArguments validates argument count and types for a call
// to fn, coercing integer/float literals when possible.
func (v *Checker) checkArguments(
	name string,
	fn Nature,
	arguments []ast.Node,
	node ast.Node,
) (Nature, *token.Error) {
	if fn.IsUnknown(&v.config.Cache) {
		return Nature{}, nil
	}

	// Functions must return 1 or 2 values (value, or value+error).
	numOut := fn.NumOut()
	if numOut == 0 {
		return Nature{}, &token.Error{
			Range:   node.Location(),
			Message: fmt.Sprintf("func %v doesn't return value", name),
		}
	}
	if numOut > 2 {
		return Nature{}, &token.Error{
			Range:   node.Location(),
			Message: fmt.Sprintf("func %v returns more then two values", name),
		}
	}

	// Method receivers show up as the first In() parameter.
	// Subtract it from the count and skip it when indexing
	// so callers don't need to account for it.
	fnNumIn := fn.NumIn()
	if fn.Method { // TODO: Move subtraction to the Nature.NumIn() and Nature.In() methods.
		fnNumIn--
	}
	fnInOffset := 0
	if fn.Method {
		fnInOffset = 1
	}

	// Arity check before per-arg validation.
	var err *token.Error
	isVariadic := fn.IsVariadic()
	if isVariadic {
		// Variadic: need at least fnNumIn-1 fixed args.
		if len(arguments) < fnNumIn-1 {
			err = &token.Error{
				Range:   node.Location(),
				Message: fmt.Sprintf("not enough arguments to call %v", name),
			}
		}
	} else {
		if len(arguments) > fnNumIn {
			err = &token.Error{
				Range:   node.Location(),
				Message: fmt.Sprintf("too many arguments to call %v", name),
			}
		}
		if len(arguments) < fnNumIn {
			err = &token.Error{
				Range:   node.Location(),
				Message: fmt.Sprintf("not enough arguments to call %v", name),
			}
		}
	}

	if err != nil {
		// Still visit args so a later patcher pass can fix natures
		// on their nodes even though the call itself is wrong.
		for _, arg := range arguments {
			_ = v.visit(arg)
		}
		return fn.Out(&v.config.Cache, 0), err
	}

	for i, arg := range arguments {
		argNature := v.visit(arg)

		var in Nature
		if isVariadic && i >= fnNumIn-1 {
			// Variadic tail: fn(xs ...int) stores []int in reflect;
			// use the element type (int) for per-arg comparison.
			in = fn.InElem(&v.config.Cache, fnNumIn-1+fnInOffset)
		} else {
			in = fn.In(&v.config.Cache, i+fnInOffset)
		}

		// Promote integer literal to float when the param is float.
		// This rewrites the AST node in place so the compiler emits
		// a FloatLit instead of an IntLit.
		if in.IsFloat && argNature.IsInteger {
			traverseAndReplaceIntegerNodesWithFloatNodes(&arguments[i], in)
			continue
		}

		// Narrow integer literal to the exact int kind the param expects
		// (e.g. int → int32) so the VM doesn't have to convert at runtime.
		if in.IsInteger && argNature.IsInteger && argNature.Kind != in.Kind {
			traverseAndReplaceIntegerNodesWithIntegerNodes(&arguments[i], in)
			continue
		}

		if argNature.Nil {
			// nil is valid for pointer or interface params.
			if in.Kind == reflect.Ptr || in.Kind == reflect.Interface {
				continue
			}
			return Nature{}, &token.Error{
				Range:   arg.Location(),
				Message: fmt.Sprintf("cannot use nil as argument (type %s) to call %v", in.String(), name),
			}
		}

		// Check assignability against the *original* arg type first
		// (func may accept *time.Time, and we should not deref that away).
		assignable := argNature.AssignableTo(in)

		// If the arg is a pointer and the func wants the base type,
		// that's also fine — the compiler will emit OpDeref.
		if !assignable && argNature.IsPointer() {
			nt := argNature.Deref(&v.config.Cache)
			assignable = nt.AssignableTo(in)
		}

		if !assignable && !argNature.IsUnknown(&v.config.Cache) {
			return Nature{}, &token.Error{
				Range:   arg.Location(),
				Message: fmt.Sprintf("cannot use %s as argument (type %s) to call %v ", argNature.String(), in.String(), name),
			}
		}
	}

	return fn.Out(&v.config.Cache, 0), nil
}

// traverseAndReplaceIntegerNodesWithFloatNodes rewrites
// integer literals and their enclosing arithmetic to float
// so they match a func parameter that expects float.
func traverseAndReplaceIntegerNodesWithFloatNodes(node *ast.Node, newNature Nature) {
	switch (*node).(type) {
	case *ast.IntLit:
		*node = &ast.FloatLit{Value: float64((*node).(*ast.IntLit).Value)}
		(*node).SetType(newNature.Type)
	case *ast.UnaryExpr:
		unaryNode := (*node).(*ast.UnaryExpr)
		traverseAndReplaceIntegerNodesWithFloatNodes(&unaryNode.Node, newNature)
	case *ast.BinaryExpr:
		binaryNode := (*node).(*ast.BinaryExpr)
		switch binaryNode.Operator {
		case "+", "-", "*":
			traverseAndReplaceIntegerNodesWithFloatNodes(&binaryNode.Left, newNature)
			traverseAndReplaceIntegerNodesWithFloatNodes(&binaryNode.Right, newNature)
		}
	}
}

// traverseAndReplaceIntegerNodesWithIntegerNodes narrows an
// integer literal's type to match a specific int kind (e.g. int32).
func traverseAndReplaceIntegerNodesWithIntegerNodes(node *ast.Node, newNature Nature) {
	switch (*node).(type) {
	case *ast.IntLit:
		(*node).SetType(newNature.Type)
	case *ast.UnaryExpr:
		(*node).SetType(newNature.Type)
		unaryNode := (*node).(*ast.UnaryExpr)
		traverseAndReplaceIntegerNodesWithIntegerNodes(&unaryNode.Node, newNature)
	case *ast.BinaryExpr:
		// TODO: Binary node return type is dependent on the type of the operands. We can't just change the type of the node.
		binaryNode := (*node).(*ast.BinaryExpr)
		switch binaryNode.Operator {
		case "+", "-", "*":
			traverseAndReplaceIntegerNodesWithIntegerNodes(&binaryNode.Left, newNature)
			traverseAndReplaceIntegerNodesWithIntegerNodes(&binaryNode.Right, newNature)
		}
	}
}

// predicateNode wraps the inner expression as a func(any)T
// so it can be passed to filter/map builtins.
func (v *Checker) predicateNode(node *ast.PredicateExpr) Nature {
	nt := v.visit(node.Node)
	// Build the output type list for reflect.FuncOf.
	// Unknown body → any; nil body → no output (void predicate).
	var out []reflect.Type
	if nt.IsUnknown(&v.config.Cache) {
		out = append(out, anyType)
	} else if !nt.Nil {
		out = append(out, nt.Type)
	}
	// Ref carries the inner Nature so callers (e.g. filter builtins)
	// can inspect the element type without unpacking the func signature.
	n := v.config.Cache.FromType(reflect.FuncOf(anyTypeSlice, out, false))
	n.Ref = &nt
	return n
}

// conditionalNode infers the result type of a ternary or elvis
// expression by comparing the two branch types.
func (v *Checker) conditionalNode(node *ast.ConditionalExpr) Nature {
	c := v.visit(node.Cond)
	c = c.Deref(&v.config.Cache)
	// JEXL is truthy — any non-nil/non-false value is accepted.
	// We visit cond for side-effects (setting natures) but don't
	// enforce that it must be bool.
	_ = c

	t1 := v.visit(node.Exp1)
	var t2 Nature
	if node.Exp2 != nil {
		t2 = v.visit(node.Exp2)
	} else {
		// Elvis operator (x ?: y) has no false-branch in the AST;
		// treat the missing branch as nil.
		t2 = v.config.Cache.NatureOf(nil)
	}

	// Return whichever branch is non-nil so the type propagates.
	if t1.Nil && !t2.Nil {
		return t2
	}
	if !t1.Nil && t2.Nil {
		return t1
	}
	if t1.Nil && t2.Nil {
		return v.config.Cache.NatureOf(nil)
	}
	if t1.AssignableTo(t2) {
		// When both branches are arrays, only keep the typed array
		// nature if the element types are mutually assignable;
		// otherwise widen to []any to avoid a false type mismatch.
		if t1.IsArray() && t2.IsArray() {
			e1 := t1.Elem(&v.config.Cache)
			e2 := t2.Elem(&v.config.Cache)
			if !e1.AssignableTo(e2) || !e2.AssignableTo(e1) {
				return v.config.Cache.FromType(arrayType)
			}
		}
		return t1
	}
	// Branches have incompatible types — unknown at compile time.
	return Nature{}
}

// arrayNode infers the element type of an array literal.
// Homogeneous elements yield a typed slice; mixed kinds widen to []any.
func (v *Checker) arrayNode(node *ast.ArrayLit) Nature {
	// If all elements share the same Kind, produce a typed array
	// (e.g. []int) so downstream operations can infer elem type.
	// Mixed kinds widen to []any.
	var prev Nature
	allElementsAreSameType := true
	for i, node := range node.Nodes {
		curr := v.visit(node)
		if i > 0 {
			if curr.Kind != prev.Kind {
				allElementsAreSameType = false
			}
		}
		prev = curr
	}
	if allElementsAreSameType {
		return prev.MakeArrayOf(&v.config.Cache)
	}
	return v.config.Cache.FromType(arrayType)
}

// mapNode determines whether the literal uses only string keys
// (map[string]any) or mixed keys (map[any]any).
func (v *Checker) mapNode(node *ast.MapLit) Nature {
	// Distinguish map[string]any (common fast path) from
	// map[any]any (non-string keys, e.g. {1: "a"}).
	hasNonStringKey := false
	for _, pair := range node.Pairs {
		v.visit(pair)
		if p, ok := pair.(*ast.KeyValueExpr); ok {
			switch p.Key.(type) {
			case *ast.StringLit:
				// string key — stays on fast path
			default:
				hasNonStringKey = true
			}
		}
	}
	if hasNonStringKey {
		return v.config.Cache.FromType(anyMapType)
	}
	return v.config.Cache.FromType(mapType)
}

// setNode visits elements for side-effects; always returns
// *util.HashSet regardless of element types (sets are untyped).
func (v *Checker) setNode(node *ast.SetLit) Nature {
	for _, elem := range node.Elements {
		v.visit(elem)
	}
	return v.config.Cache.FromType(setType)
}

// bitOpNode checks bitwise and shift operators (~, &, |, ^, <<, >>, >>>).
// All operands are coerced to int64; shift counts must be non-negative.
func (v *Checker) bitOpNode(node *ast.BitExpr) Nature {
	l := v.visit(node.Left)
	l = l.Deref(&v.config.Cache)
	if node.Right == nil {
		// unary ~: allow any numeric or unknown type (JEXL coerces via toLong)
		if l.IsInteger || l.IsFloat || l.IsBool() || l.IsUnknown(&v.config.Cache) {
			return v.config.Cache.FromType(int64Type)
		}
		return v.error(node, "invalid operation: %s (non-integer type %s)", node.Operator, l.String())
	}
	r := v.visit(node.Right)
	r = r.Deref(&v.config.Cache)
	// JEXL coerces booleans and strings to long for bitwise operations.
	isCoercible := func(n Nature) bool {
		return n.IsInteger || n.IsFloat || n.IsBool() || n.IsString() || n.IsUnknown(&v.config.Cache)
	}
	if !(isCoercible(l) && isCoercible(r)) {
		return v.error(node, "invalid operation: %s (non-integer operands %s and %s)", node.Operator, l.String(), r.String())
	}
	// For shift operators, the right operand (shift count) must be non-negative when known at compile time.
	if node.Operator == "<<" || node.Operator == ">>" || node.Operator == ">>>" {
		if rn, ok := node.Right.(*ast.IntLit); ok && rn.Value < 0 {
			return v.error(node, "invalid operation: %s (negative shift count %d)", node.Operator, rn.Value)
		}
	}
	return v.config.Cache.FromType(int64Type)
}

// strictEqualNode handles === / !==. Unlike ==, type mismatches are
// not errors — they simply mean unequal, so we always return bool.
func (v *Checker) strictEqualNode(node *ast.StrictEqualExpr) Nature {
	v.visit(node.Left)
	v.visit(node.Right)
	// Strict equality never errors on type mismatch — different types simply means not equal.
	return v.config.Cache.FromType(boolType)
}

// lambdaNode builds a func Nature from the param list and body.
// Params are pushed as unknown-typed vars while the body is visited.
func (v *Checker) lambdaNode(node *ast.LambdaExpr) Nature {
	// Push each parameter as an unknown-typed variable scope.
	for _, p := range node.Params {
		v.varScopes = append(v.varScopes, varScope{p, Nature{}})
	}
	bodyNature := v.visit(node.Body)
	// Pop the param scopes.
	v.varScopes = v.varScopes[:len(v.varScopes)-len(node.Params)]

	// Build input types: all params are any.
	inTypes := make([]reflect.Type, len(node.Params))
	for i := range inTypes {
		inTypes[i] = anyType
	}
	// Output type from the body. Always include at least one output type
	// so the lambda can be called as an expression (nil body returns null).
	var outTypes []reflect.Type
	if bodyNature.IsUnknown(&v.config.Cache) || bodyNature.Nil {
		outTypes = []reflect.Type{anyType}
	} else {
		outTypes = []reflect.Type{bodyNature.Type}
	}
	nt := v.config.Cache.FromType(reflect.FuncOf(inTypes, outTypes, false))
	nt.Ref = &bodyNature
	return nt
}

// pairNode visits key and value for side-effects (setting natures)
// but returns nil — the parent mapNode consumes the pair directly.
func (v *Checker) pairNode(node *ast.KeyValueExpr) Nature {
	v.visit(node.Key)
	v.visit(node.Value)
	return v.config.Cache.NatureOf(nil)
}

// blockNode visits each statement and returns the nature of
// the last one — blocks evaluate to their final expression.
func (v *Checker) blockNode(node *ast.BlockStmt) Nature {
	var last Nature
	for _, stmt := range node.Statements {
		last = v.visit(stmt)
	}
	return last
}

// varNode pushes a new var scope entry, then visits the
// initializer to determine (and record) the variable's type.
func (v *Checker) varNode(node *ast.VarDecl) Nature {
	// Reserve the slot first (with unknown nature) so that the
	// initializer can reference the var itself (e.g. recursive
	// lambda) without binding to a global with the same name.
	idx := len(v.varScopes)
	v.varScopes = append(v.varScopes, varScope{node.Name, Nature{}})
	valueNature := v.visit(node.Value)
	// `var x = nil` — initializer is nil, so type is unknown.
	// Later assignments determine the actual type.
	if _, isNil := node.Value.(*ast.NilLit); isNil {
		valueNature = Nature{}
	}
	v.varScopes[idx].nature = valueNature
	return valueNature
}

// incrDecrNode validates that ++ / -- targets are lvalues
// (identifiers or member accesses), then returns the target's type.
func (v *Checker) incrDecrNode(node *ast.IncDecStmt) Nature {
	// Only lvalues (bare names or member accesses) can be
	// incremented/decremented; literals and calls cannot.
	switch node.Target.(type) {
	case *ast.Ident, *ast.MemberExpr:
		// valid lvalue
	default:
		return v.error(node, "invalid %s operand: must be an identifier or member access", node.Op)
	}
	return v.visit(node.Target)
}

// foreachNode checks the collection type and binds the loop
// variable to the element type for the duration of the body.
func (v *Checker) foreachNode(node *ast.ForeachStmt) Nature {
	collection := v.visit(node.Collection)
	collection = collection.Deref(&v.config.Cache)

	var elemNature Nature
	if collection.IsArray() {
		elemNature = collection.Elem(&v.config.Cache)
	} else if collection.IsMap() {
		// Iterating over a map yields values
		elemNature = collection.Elem(&v.config.Cache)
	} else if collection.Type == setType || collection.Type == setValueType {
		// Iterating over a set yields elements
		elemNature = Nature{}
	} else if !collection.IsUnknown(&v.config.Cache) {
		return v.error(node.Collection, "for-each loop requires an array (got %v)", collection.String())
	}

	v.varScopes = append(v.varScopes, varScope{node.Var, elemNature})
	v.loopDepth++
	v.visit(node.Body)
	v.loopDepth--
	v.varScopes = v.varScopes[:len(v.varScopes)-1]

	return Nature{}
}

// whileNode visits the condition and body; loopDepth gates
// break/continue validation inside the body.
func (v *Checker) whileNode(node *ast.WhileStmt) Nature {
	cond := v.visit(node.Cond)
	cond = cond.Deref(&v.config.Cache)
	// JEXL is truthy — any value is a valid condition.
	// Visit cond only to assign natures to its sub-nodes.
	_ = cond
	v.loopDepth++ // enables break/continue inside the body
	v.visit(node.Body)
	v.loopDepth--
	return Nature{}
}

// forNode visits init, cond, body, and post in execution order.
// All three header clauses are optional.
func (v *Checker) forNode(node *ast.ForStmt) Nature {
	// All three clauses are optional (for ;; {} is valid).
	if node.Init != nil {
		v.visit(node.Init)
	}
	if node.Cond != nil {
		v.visit(node.Cond)
	}
	v.loopDepth++
	v.visit(node.Body)
	// Post runs after body — visit after so break/continue
	// inside the body don't affect the post expression.
	if node.Post != nil {
		v.visit(node.Post)
	}
	v.loopDepth--
	return Nature{}
}

// doWhileNode visits body before condition, matching execution order
// so natures are set on nodes in the order the VM will see them.
func (v *Checker) doWhileNode(node *ast.DoWhileStmt) Nature {
	// Body runs first; condition is evaluated after — mirror
	// that order here so natures are set in execution order.
	v.loopDepth++
	v.visit(node.Body)
	v.loopDepth--
	v.visit(node.Cond)
	return Nature{}
}

// breakNode errors if break appears outside a loop or switch.
func (v *Checker) breakNode(node *ast.BreakStmt) Nature {
	if v.loopDepth == 0 && v.switchDepth == 0 {
		return v.error(node, "break outside of loop or switch")
	}
	return Nature{}
}

// continueNode errors if continue appears outside a loop.
func (v *Checker) continueNode(node *ast.ContinueStmt) Nature {
	if v.loopDepth == 0 {
		return v.error(node, "continue outside of loop")
	}
	return Nature{}
}

// returnNode yields the returned value's nature so callers
// (e.g. lambdaNode) can infer the enclosing function's return type.
func (v *Checker) returnNode(node *ast.ReturnStmt) Nature {
	if node.Value != nil {
		return v.visit(node.Value)
	}
	return Nature{}
}

// tryNode visits the try body, then the catch body (with the catch
// variable scoped as unknown), then the finally body.
func (v *Checker) tryNode(node *ast.TryStmt) Nature {
	v.visit(node.Body)
	if node.CatchBody != nil {
		// Bind the catch variable as unknown-typed — the thrown
		// value can be anything; the body determines how it's used.
		if node.CatchVar != "" {
			v.varScopes = append(v.varScopes, varScope{node.CatchVar, Nature{}})
		}
		v.visit(node.CatchBody)
		if node.CatchVar != "" {
			v.varScopes = v.varScopes[:len(v.varScopes)-1]
		}
	}
	if node.FinallyBody != nil {
		v.visit(node.FinallyBody)
	}
	return Nature{}
}

// throwNode visits the thrown value for side-effects; returns
// unknown — throw terminates the current expression path.
func (v *Checker) throwNode(node *ast.ThrowStmt) Nature {
	v.visit(node.Value)
	return Nature{}
}

// switchNode visits the subject and all case values/bodies.
// switchDepth is incremented so break inside a case is valid.
func (v *Checker) switchNode(node *ast.SwitchStmt) Nature {
	v.visit(node.Subject)
	v.switchDepth++ // allows break inside case bodies
	for _, c := range node.Cases {
		// Visit each case value (the expressions after `case`).
		for _, val := range c.Values {
			v.visit(val)
		}
		v.visit(c.Body)
	}
	if node.Default != nil {
		v.visit(node.Default)
	}
	v.switchDepth--
	return Nature{}
}

// emptyNode handles the empty() built-in. The operand is visited
// for side-effects; the result is always bool.
func (v *Checker) emptyNode(node *ast.EmptyExpr) Nature {
	// empty(x) — tests whether x is nil, "", [], or 0.
	// Always returns bool regardless of the operand type.
	v.visit(node.Value)
	return v.config.Cache.FromType(boolType)
}

// sizeNode handles the size() built-in. The operand is visited
// for side-effects; the result is always int64 (Java long).
func (v *Checker) sizeNode(node *ast.SizeExpr) Nature {
	// size(x) — length of a string, array, or map.
	// Returns int64 to match Java's long.
	v.visit(node.Value)
	return v.config.Cache.FromType(int64Type)
}

// regexNode returns *regexp.Regexp regardless of the pattern;
// syntax errors are caught by the parser, not the checker.
func (v *Checker) regexNode(node *ast.RegexLit) Nature {
	return v.config.Cache.FromType(regexpType)
}

// newNode handles `new ClassName(...)`. If a registry is configured,
// it verifies the class exists; args are visited for side-effects.
// The compile-time type is map[string]any — the registry resolves
// the concrete type at runtime.
func (v *Checker) newNode(node *ast.NewExpr) Nature {
	// new ClassName(...) — verify the class is registered
	// when a registry is present.  The constructed object
	// is treated as map[string]any at compile time; the
	// actual type is resolved by the registry at runtime.
	if v.config.Registry != nil {
		if _, ok := v.config.Registry.Lookup(node.ClassName); !ok {
			return v.error(node, "unknown class %q (not in registry)", node.ClassName)
		}
	}
	for _, arg := range node.Args {
		v.visit(arg)
	}
	return v.config.Cache.FromType(mapType)
}

// namespaceCallNode handles Namespace::func() calls. If a registry
// is configured, it verifies the namespace exists; returns unknown
// because the concrete return type is resolved by the registry at runtime.
func (v *Checker) namespaceCallNode(node *ast.NamespaceCallExpr) Nature {
	// Namespace::function() calls — validate the namespace
	// exists in the registry but return unknown; the actual
	// return type depends on the function looked up at runtime.
	if v.config.Registry != nil {
		if _, ok := v.config.Registry.Lookup(node.Namespace); !ok {
			return v.error(node, "unknown namespace %q", node.Namespace)
		}
	}
	for _, arg := range node.Args {
		v.visit(arg)
	}
	return Nature{}
}

// assignNode validates the LHS is an lvalue and visits both sides.
// Type compatibility is not enforced — JEXL is dynamically typed.
// The result type is the RHS type, matching compound-assignment semantics.
func (v *Checker) assignNode(node *ast.AssignStmt) Nature {
	// Only lvalues are valid on the left side.
	var targetNature Nature
	switch t := node.Target.(type) {
	case *ast.Ident:
		targetNature = v.visit(t)
	case *ast.MemberExpr:
		targetNature = v.visit(t)
	default:
		return v.error(node, "invalid assignment target: left-hand side must be an identifier or member access")
	}

	valueNature := v.visit(node.Value)
	// targetNature is visited for side-effects (natures on sub-nodes)
	// but we don't enforce type compatibility here — JEXL is dynamic.
	_ = targetNature

	// For both = and compound assignments (+=, -=, …) propagate
	// the RHS type as the expression result.
	return valueNature
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

// dotChain flattens a chain of MemberExprs into a dotted
// string (e.g. "java.util.Arrays"), returning false on failure.
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
