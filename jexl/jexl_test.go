// Copyright (c) 2026 Harness Inc.
// Copyright (c) 2018 Anton Medvedev
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package jexl

import (
	"strings"
	"testing"

	"github.com/harness/go-jexl/jexl/coerce"
	"github.com/harness/go-jexl/jexl/vm"
)

func TestSyntax(t *testing.T) {
	ctx := map[string]any{
		"a":    10,
		"b":    3,
		"name": "Alice",
		"age":  30,
		"user": map[string]any{
			"name": "Bob",
		},
		"arr":  []any{1, 2, 3},
		"m":    map[string]any{"key": "val"},
		"flag": true,
		"n":    2,
		"prop": "name",
	}

	tests := []struct {
		expr string
		want any
	}{
		// Literals
		{`42`, 42},
		{`-7`, -7},
		{`3.14`, 3.14},
		{`true`, true},
		{`false`, false},
		{`null`, nil},
		{`"hello"`, "hello"},
		{`'world'`, "world"},
		{`[1, 2, 3]`, []any{1, 2, 3}},
		{`{"key": "value"}`, map[string]any{"key": "value"}},

		// Arithmetic
		{`1 + 2`, 3},
		{`10 - 4`, 6},
		{`3 * 4`, 12},
		{`7 / 2`, 3},
		{`7 div 2`, 3},
		{`7 % 3`, 1},
		{`5 mod 2`, 1},
		{`14 mod 5`, 4},
		{`2 ** 10`, float64(1024)},
		{`-a`, -10},

		// Comparison
		{`1 == 1`, true},
		{`1 eq 1`, true},
		{`1 != 2`, true},
		{`1 ne 2`, true},
		{`1 === 1`, true},
		{`1 !== 2`, true},
		{`1 < 2`, true},
		{`1 lt 2`, true},
		{`3 lt 2`, false},
		{`1 <= 1`, true},
		{`1 le 1`, true},
		{`2 le 1`, false},
		{`2 > 1`, true},
		{`1 > 2`, false},
		{`2 gt 1`, true},
		{`1 gt 2`, false},
		{`2 >= 2`, true},
		{`2 ge 2`, true},
		{`1 ge 2`, false},

		// Logical
		{`true && false`, false},
		{`true and false`, false},
		{`false || true`, true},
		{`false or true`, true},
		{`!true`, false},
		{`not true`, false},
		{`null ?? "x"`, "x"},
		{`false ?? "x"`, false},
		{`null ?: "x"`, "x"},
		{`false ?: "x"`, "x"},

		// Bitwise
		{`6 | 3`, 7},
		{`6 ^ 3`, 5},
		{`6 & 3`, 2},
		{`~0`, -1},
		{`1 << 3`, 8},
		{`8 >> 1`, 4},
		{`8 >>> 1`, 4},

		// Ternary
		{`true ? "yes" : "no"`, "yes"},
		{`false ? "yes" : "no"`, "no"},

		// Member access
		{`user.name`, "Bob"},
		{`arr[0]`, 1},
		{`user['name']`, "Bob"},
		{`user.'name'`, "Bob"},
		{`m["key"]`, "val"},

		// Template literals
		{"` Hello ${name}`", " Hello Alice"},
		{"`${a} + ${b}`", "10 + 3"},

		// String operators
		{`"hello" =^ "he"`, true},
		{`"hello" =$ "lo"`, true},
		{`"hello" =~ "hel.*"`, true},
		{`"hello" !^ "he"`, false},
		{`"hello" !$ "lo"`, false},
		{`"hello" !~ "xyz"`, true},

		// Membership
		{`"admin" in ["admin", "mod"]`, true},
		{`"admin" in {"admin": true}`, true},
		{`1 =~ [1, 2, 3]`, true},

		// Range
		{`1..5`, []int{1, 2, 3, 4, 5}},

		// instanceof
		{`a instanceof Integer`, true},
		{`a !instanceof String`, true},

		// Builtins: empty / size
		{`empty(null)`, true},
		{`empty("")`, true},
		{`empty(0)`, true},
		{`empty([])`, true},
		{`empty name`, false}, // name="Alice", not empty
		{`empty dummy`, true}, // dummy does not exist
		{`size("hello")`, int64(5)},
		{`size([1,2,3])`, int64(3)},
		{`size arr`, int64(3)}, // arr=[1,2,3]

		// Variables
		{`var x = 10; x`, 10},
		{`let y = 5; y`, 5},
		{`const PI = 3; PI`, 3},

		// Assignment operators
		{`var x = 10; x += 1; x`, 11},
		{`var x = 10; x -= 1; x`, 9},
		{`var x = 10; x *= 2; x`, 20},
		{`var x = 10; x /= 2; x`, 5.0},
		{`var x = 10; x %= 3; x`, 1},
		{`var x = 10; x++; x`, 11},
		{`var x = 10; x--; x`, 9},
		{`var x = 10; ++x`, 11},
		{`var x = 10; --x`, 9},

		// Loops
		{`var s = 0; for (var i : 1..3) { s += i }; s`, 6},
		{`var s = 0; for (var i = 0; i < 3; i++) { s += i }; s`, 3},
		{`var s = 0; var i = 0; while (i < 3) { s += i; i++ }; s`, 3},
		{`var s = 0; var i = 0; do { s += i; i++ } while (i < 3); s`, 3},

		// Control flow
		{`var s = 0; for (var i : 1..5) { if (i == 3) { break }; s += i }; s`, 3},
		{`var s = 0; for (var i : 1..5) { if (i == 3) { continue }; s += i }; s`, 12},
		{`function f() { return 42 }; f()`, 42},

		// if / else if / else
		{`if (a > 5) { "big" } else { "small" }`, "big"},
		{`if (a < 5) { "small" } else if (a == 10) { "ten" } else { "other" }`, "ten"},

		// try / catch
		{`try { throw "oops" } catch (let e) { e }`, "oops"},

		// switch — arrow form
		{`var r = switch (n) { case 1 -> "one"; case 2, 3 -> "two or three"; default -> "other" }; r`, "two or three"},

		// switch — colon form
		{`var r = ""; switch (n) { case 1: r = "one"; break; case 2: r = "two"; break; default: r = "other" }; r`, "two"},

		// Comments
		{`1 + 1 // add`, 2},
		{`1 + 1 ## add`, 2},
		{`1 + /* block */ 1`, 2},

		// Functions
		{`function add(x, y) { x + y }; add(3, 4)`, 7},

		// Lambdas
		{`var double = (x) -> x * 2; double(5)`, 10},
		{`var add = (x, y) -> x + y; add(3, 4)`, 7},
		{`var ops = {"double": (x) -> x * 2}; ops["double"](4)`, 8},

		// Closures capture value at definition
		{`var t = 20; var s = function(x, y) { x + y + t }; t = 54; s(15, 7)`, 42},

		// String method calls
		{`"hello".toUpperCase()`, "HELLO"},
		{`"  hello  ".trim()`, "hello"},
		{`"hello world".indexOf("world")`, 6},
		{`"hello".lastIndexOf("l")`, 3},
		{`"hello".substring(1, 3)`, "el"},
		{`"hello".contains("ell")`, true},
		{`"hello".startsWith("he")`, true},
		{`"hello".endsWith("lo")`, true},
		{`"hello".replace("l", "r")`, "herro"},
		{`"a,b,c".split(",")`, []any{"a", "b", "c"}},
		{`"hello".repeat(2)`, "hellohello"},
		{`"hello".matches("hel.*")`, true},
		{`"hello".equals("hello")`, true},
		{`"hello".equalsIgnoreCase("HELLO")`, true},
		{`"hello".compareTo("hello")`, 0},
		{`"hello".isEmpty()`, false},
		{`"hello".isBlank()`, false},
		{`"hello".concat(" world")`, "hello world"},
		{`"hello".length()`, 5},

		// Array/list method calls
		{`arr.size()`, 3},
		{`arr.isEmpty()`, false},
		{`arr.contains(2)`, true},
		{`arr.get(0)`, 1},
		{`arr.indexOf(2)`, 1},

		// Map method calls
		{`m.size()`, 1},
		{`m.isEmpty()`, false},
		{`m.containsKey("key")`, true},
		{`m.get("key")`, "val"},

		// safe navigation
		{`user?.name`, "Bob"},
		{`dummy?.name`, nil},

		// dynamic property via backtick interpolation
		{"user.`${prop}`", "Bob"},
		{"dummy?.`${prop}`", nil},
	}

	for _, tc := range tests {
		out, err := Eval(tc.expr, ctx)
		if err != nil {
			t.Errorf("expr %q: %v", tc.expr, err)
			continue
		}
		if !coerce.DeepEqual(out, tc.want) {
			t.Errorf("expr %q: got %v (%T), want %v (%T)", tc.expr, out, out, tc.want, tc.want)
		}
	}
}

//
// Limits
//

// EnsuVerifiesres that a loop exceeding the limit returns an error.
func TestMaxIterations(t *testing.T) {
	program, err := Compile(`var s = 0; for (var i = 0; i < 100; i++) { s = s + i }; s`,
		WithMaxIterations(10),
	)
	if err != nil {
		t.Fatal(err)
	}
	_, err = Run(program, nil)
	if err == nil {
		t.Fatal("expected error for exceeded max iterations, got nil")
	}
	if !strings.Contains(err.Error(), "maximum iterations") {
		t.Fatalf("expected 'maximum iterations' in error, got: %v", err)
	}
}

// Verifies allocating beyond the memory budget returns an error.
func TestMaxMemory(t *testing.T) {
	// A range of 1000 elements will call memGrow(1000), exceeding a budget of 10.
	program, err := Compile(`1..1000`,
		WithMaxMemory(10),
	)
	if err != nil {
		t.Fatal(err)
	}
	_, err = Run(program, nil)
	if err == nil {
		t.Fatal("expected error for exceeded max memory, got nil")
	}
	if !strings.Contains(err.Error(), "memory budget exceeded") {
		t.Fatalf("expected 'memory budget exceeded' in error, got: %v", err)
	}
}

// Verifies that a deeply nested expression exceeding the node limit returns an error.
func TestMaxNodes(t *testing.T) {
	// Build an expression with many nodes: 1+1+1+... repeated 50 times.
	parts := make([]string, 50)
	for i := range parts {
		parts[i] = "1"
	}
	input := strings.Join(parts, "+")

	_, err := Compile(input, WithMaxNodes(5))
	if err == nil {
		t.Fatal("expected error for exceeded max nodes, got nil")
	}
	if !strings.Contains(err.Error(), "exceeds maximum allowed nodes") {
		t.Fatalf("expected 'exceeds maximum allowed nodes' in error, got: %v", err)
	}
}

//
// Safe mode
//

// TestSafeMode_contextAssignBlocked verifies that assigning to a context
// variable is rejected when WithSafe is enabled.
func TestSafeMode_contextAssignBlocked(t *testing.T) {
	ctx := map[string]any{"x": 10}
	_, err := Compile(`x = 99`, WithSafe())
	if err == nil {
		t.Fatal("expected error assigning to context variable in safe mode, got nil")
	}
	if ctx["x"] != 10 {
		t.Fatalf("context was mutated: x = %v", ctx["x"])
	}
}

// TestSafeMode_contextMemberAssignBlocked verifies that assigning to a nested
// context field is rejected when WithSafe is enabled.
func TestSafeMode_contextMemberAssignBlocked(t *testing.T) {
	ctx := map[string]any{"user": map[string]any{"name": "Alice"}}
	_, err := Compile(`user.name = "Eve"`, WithSafe())
	if err == nil {
		t.Fatal("expected error assigning to context member in safe mode, got nil")
	}
	user := ctx["user"].(map[string]any)
	if user["name"] != "Alice" {
		t.Fatalf("context was mutated: user.name = %v", user["name"])
	}
}

// TestSafeMode_contextCompoundAssignBlocked verifies that compound assignment
// to a context variable is rejected when WithSafe is enabled.
func TestSafeMode_contextCompoundAssignBlocked(t *testing.T) {
	_, err := Compile(`x += 1`, WithSafe())
	if err == nil {
		t.Fatal("expected error for compound assignment to context variable in safe mode, got nil")
	}
}

// TestSafeMode_localVarAllowed verifies that declaring and assigning local
// variables still works in safe mode.
func TestSafeMode_localVarAllowed(t *testing.T) {
	program, err := Compile(`var n = 5; n += 3; n`, WithSafe())
	if err != nil {
		t.Fatal(err)
	}
	out, err := Run(program, nil)
	if err != nil {
		t.Fatal(err)
	}
	if out != 8 {
		t.Fatalf("expected 8, got %v", out)
	}
}

//
// Benchmarks
//

func BenchmarkRun(b *testing.B) {
	params := make(map[string]any)
	params["Origin"] = "MOW"
	params["Country"] = "RU"
	params["Adults"] = 1
	params["Value"] = 100

	program, err := Compile(`(Origin == "MOW" || Country == "RU") && (Value >= 100 || Adults == 1)`, WithContext(params))
	if err != nil {
		b.Fatal(err)
	}

	var out any

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, err = vm.Run(program, params)
	}
	b.StopTimer()

	if err != nil {
		b.Fatal(err)
	}
	if out.(bool) != true {
		b.Fatalf("expected true, got %v", out)
	}
}

func BenchmarkEval(b *testing.B) {
	params := make(map[string]any)
	params["Origin"] = "MOW"
	params["Country"] = "RU"
	params["Adults"] = 1
	params["Value"] = 100

	var out any
	var err error

	b.ResetTimer()
	for n := 0; n < b.N; n++ {
		out, err = Eval(`(Origin == "MOW" || Country == "RU") && (Value >= 100 || Adults == 1)`, params)
	}
	b.StopTimer()

	if err != nil {
		b.Fatal(err)
	}
	if out.(bool) != true {
		b.Fatalf("expected true, got %v", out)
	}
}
