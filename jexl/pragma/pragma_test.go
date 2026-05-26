// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package pragma

import (
	"testing"

	"github.com/harness/go-jexl/jexl/config"
)

// Ensure Parse returns empty slice and unchanged input
// when no pragma present.
func TestParse_noDirective(t *testing.T) {
	out, directives := Parse("foo + bar")
	if len(directives) != 0 {
		t.Fatalf("expected no directives, got %d", len(directives))
	}
	if out != "foo + bar" {
		t.Fatalf("expected input unchanged, got %q", out)
	}
}

// Ensure Parse extracts a single-quoted value.
func TestParse_singleQuotedValue(t *testing.T) {
	_, directives := Parse("#pragma jexl.namespace.Math 'Math'")
	if len(directives) != 1 {
		t.Fatalf("expected 1 directive, got %d", len(directives))
	}
	if directives[0].Key != "jexl.namespace.Math" {
		t.Fatalf("unexpected key: %q", directives[0].Key)
	}
	if directives[0].Val != "Math" {
		t.Fatalf("unexpected val: %q", directives[0].Val)
	}
}

// Ensure Parse extracts a double-quoted value.
func TestParse_doubleQuotedValue(t *testing.T) {
	_, directives := Parse(`#pragma jexl.import "mylib"`)
	if len(directives) != 1 {
		t.Fatalf("expected 1 directive, got %d", len(directives))
	}
	if directives[0].Val != "mylib" {
		t.Fatalf("unexpected val: %q", directives[0].Val)
	}
}

// Ensure Parse handles a directive with no value.
func TestParse_noValue(t *testing.T) {
	_, directives := Parse("#pragma jexl.options")
	if len(directives) != 1 {
		t.Fatalf("expected 1 directive, got %d", len(directives))
	}
	if directives[0].Key != "jexl.options" {
		t.Fatalf("unexpected key: %q", directives[0].Key)
	}
	if directives[0].Val != "" {
		t.Fatalf("expected empty val, got %q", directives[0].Val)
	}
}

// Ensure Parse strips the directive from the returned string,
// preserving offsets. The trailing space after the key is
// consumed as part of the directive scan.
func TestParse_stripsDirective(t *testing.T) {
	input := "#pragma jexl.options foo + bar"
	out, _ := Parse(input)
	if len(out) != len(input) {
		t.Fatalf("expected length %d, got %d", len(input), len(out))
	}
	if out[len("#pragma jexl.options "):] != "foo + bar" {
		t.Fatalf("unexpected stripped output: %q", out)
	}
}

// Ensure Apply sets strict=true for "strict" value.
func TestApply_optionsStrict(t *testing.T) {
	conf := config.New()
	Apply(conf, []Directive{{Key: "jexl.options", Val: "strict"}})
	if !conf.Strict {
		t.Fatal("expected Strict to be true")
	}
}

// Ensure Apply sets strict=true for "+strict" value.
func TestApply_optionsPlusStrict(t *testing.T) {
	conf := config.New()
	Apply(conf, []Directive{{Key: "jexl.options", Val: "+strict"}})
	if !conf.Strict {
		t.Fatal("expected Strict to be true")
	}
}

// Ensure Apply sets strict=false for "-strict" value.
func TestApply_optionsMinusStrict(t *testing.T) {
	conf := config.New()
	conf.Strict = true
	Apply(conf, []Directive{{Key: "jexl.options", Val: "-strict"}})
	if conf.Strict {
		t.Fatal("expected Strict=false")
	}
}

// Ensure Apply ignores unknown jexl.options values.
func TestApply_optionsUnknown(t *testing.T) {
	conf := config.New()
	Apply(conf, []Directive{{Key: "jexl.options", Val: "unknown"}})
	if conf.Strict {
		t.Fatal("expected StrictExplicit to remain false")
	}
}

// Ensure Apply registers a namespace alias via jexl.namespace.*.
func TestApply_namespace(t *testing.T) {
	conf := config.New()
	Apply(conf, []Directive{{Key: "jexl.import", Val: "java.lang.String"}})
	Apply(conf, []Directive{{Key: "jexl.namespace.Str", Val: "java.lang.String"}})
	if _, ok := conf.Registry.Lookup("Str"); !ok {
		t.Fatal("expected Str to be registered")
	}
}

// Ensure Apply ignores jexl.namespace.* when the target is not in the registry.
func TestApply_namespaceMissing(t *testing.T) {
	conf := config.New()
	Apply(conf, []Directive{{Key: "jexl.namespace.Foo", Val: "no.such.Class"}})
	if _, ok := conf.Registry.Lookup("Foo"); ok {
		t.Fatal("expected Foo to not be registered")
	}
}

// Ensure Apply ignores jexl.namespace.* when alias is empty.
func TestApply_namespaceEmptyAlias(t *testing.T) {
	conf := config.New()
	Apply(conf, []Directive{{Key: "jexl.namespace.", Val: "java.lang.String"}})
}

// Ensure importPackage registers all java.lang classes.
func TestImportPackage_javaLang(t *testing.T) {
	conf := config.New()
	Apply(conf, []Directive{{Key: "jexl.import", Val: "java.lang"}})
	for _, name := range []string{
		"java.lang.Boolean", "java.lang.Byte", "java.lang.Character",
		"java.lang.Double", "java.lang.Float", "java.lang.Integer",
		"java.lang.Long", "java.lang.Short", "java.lang.String",
		"java.lang.StringBuilder", "java.lang.StringBuffer",
	} {
		if _, ok := conf.Registry.Lookup(name); !ok {
			t.Fatalf("expected %s to be registered", name)
		}
	}
}

// Ensure importPackage registers individual java.lang classes.
func TestImportPackage_javaLangIndividual(t *testing.T) {
	cases := []string{
		"java.lang.Boolean", "java.lang.Byte", "java.lang.Character",
		"java.lang.Double", "java.lang.Float", "java.lang.Integer",
		"java.lang.Long", "java.lang.Short", "java.lang.String",
		"java.lang.StringBuilder", "java.lang.StringBuffer",
	}
	for _, name := range cases {
		conf := config.New()
		Apply(conf, []Directive{{Key: "jexl.import", Val: name}})
		if _, ok := conf.Registry.Lookup(name); !ok {
			t.Fatalf("expected %s to be registered", name)
		}
	}
}

// Ensure importPackage registers all java.math classes.
func TestImportPackage_javaMath(t *testing.T) {
	conf := config.New()
	Apply(conf, []Directive{{Key: "jexl.import", Val: "java.math"}})
	for _, name := range []string{"java.lang.Math", "java.math.BigDecimal", "java.math.BigInteger"} {
		if _, ok := conf.Registry.Lookup(name); !ok {
			t.Fatalf("expected %s to be registered", name)
		}
	}
}

// Ensure importPackage registers individual java.math classes.
func TestImportPackage_javaMathIndividual(t *testing.T) {
	cases := []string{"java.lang.Math", "java.math.BigDecimal", "java.math.BigInteger"}
	for _, name := range cases {
		conf := config.New()
		Apply(conf, []Directive{{Key: "jexl.import", Val: name}})
		if _, ok := conf.Registry.Lookup(name); !ok {
			t.Fatalf("expected %s to be registered", name)
		}
	}
}

// Ensure importPackage registers all java.util classes.
func TestImportPackage_javaUtil(t *testing.T) {
	conf := config.New()
	Apply(conf, []Directive{{Key: "jexl.import", Val: "java.util"}})
	for _, name := range []string{
		"java.util.ArrayList", "java.util.LinkedList", "java.util.HashMap",
		"java.util.TreeMap", "java.util.HashSet", "java.util.LinkedHashSet",
		"java.util.Stack",
	} {
		if _, ok := conf.Registry.Lookup(name); !ok {
			t.Fatalf("expected %s to be registered", name)
		}
	}
}

// Ensure importPackage registers individual java.util classes.
func TestImportPackage_javaUtilIndividual(t *testing.T) {
	cases := []string{
		"java.util.ArrayList", "java.util.LinkedList", "java.util.HashMap",
		"java.util.TreeMap", "java.util.HashSet", "java.util.LinkedHashSet",
		"java.util.Stack",
	}
	for _, name := range cases {
		conf := config.New()
		Apply(conf, []Directive{{Key: "jexl.import", Val: name}})
		if _, ok := conf.Registry.Lookup(name); !ok {
			t.Fatalf("expected %s to be registered", name)
		}
	}
}

// Ensure importPackage is a no-op for unknown names.
func TestImportPackage_unknown(t *testing.T) {
	conf := config.New()
	Apply(conf, []Directive{{Key: "jexl.import", Val: "com.example.Unknown"}})
}

// Ensure Parse handles multiple directives in one input.
func TestParse_multipleDirectives(t *testing.T) {
	input := "#pragma jexl.options\n#pragma jexl.import 'java.lang.String'\nx + y"
	out, directives := Parse(input)
	if len(directives) != 2 {
		t.Fatalf("expected 2 directives, got %d", len(directives))
	}
	if directives[0].Key != "jexl.options" {
		t.Fatalf("unexpected first key: %q", directives[0].Key)
	}
	if directives[1].Key != "jexl.import" {
		t.Fatalf("unexpected second key: %q", directives[1].Key)
	}
	if directives[1].Val != "java.lang.String" {
		t.Fatalf("unexpected second val: %q", directives[1].Val)
	}
	if len(out) != len(input) {
		t.Fatalf("expected length %d, got %d", len(input), len(out))
	}
}
