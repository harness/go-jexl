// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package pragma

import (
	"strings"

	"github.com/harness/go-jexl/jexl/classes/java/lang"
	"github.com/harness/go-jexl/jexl/classes/java/math"
	"github.com/harness/go-jexl/jexl/classes/java/util"
	"github.com/harness/go-jexl/jexl/config"
)

// Directive represents a single parsed #pragma directive.
type Directive struct {
	Key string
	Val string
}

// Parse scans input for #pragma directives and returns
// them along with the input with each directive replaced
// by spaces (preserving offsets).
func Parse(input string) (string, []Directive) {
	var directives []Directive
	for {
		idx := strings.Index(input, "#pragma")
		if idx < 0 {
			break
		}
		rest := input[idx+len("#pragma"):]
		i := 0
		for i < len(rest) && (rest[i] == ' ' || rest[i] == '\t') {
			i++
		}
		keyStart := i
		for i < len(rest) && rest[i] != ' ' && rest[i] != '\t' && rest[i] != '\n' {
			i++
		}
		key := rest[keyStart:i]
		for i < len(rest) && (rest[i] == ' ' || rest[i] == '\t') {
			i++
		}
		val := ""
		if i < len(rest) && (rest[i] == '\'' || rest[i] == '"') {
			quote := rest[i]
			i++
			valStart := i
			for i < len(rest) && rest[i] != quote {
				i++
			}
			val = rest[valStart:i]
			if i < len(rest) {
				i++
			}
		}
		end := idx + len("#pragma") + i
		directives = append(directives, Directive{Key: key, Val: val})
		replaced := strings.Repeat(" ", end-idx)
		input = input[:idx] + replaced + input[end:]
	}
	return input, directives
}

// Apply applies directives to conf.
func Apply(conf *config.Config, directives []Directive) {
	for _, d := range directives {
		switch {
		case d.Key == "jexl.options":
			switch d.Val {
			case "strict", "+strict":
				conf.Strict = true
			case "-strict":
				conf.Strict = false
			}
		case d.Key == "jexl.import":
			importPackage(conf, d.Val)
		case strings.HasPrefix(d.Key, "jexl.namespace."):
			alias := d.Key[len("jexl.namespace."):]
			if alias != "" {
				if obj, ok := conf.Registry.Lookup(d.Val); ok {
					conf.Registry.Register(alias, obj)
				}
			}
		}
	}
}

// helper function to implement well-known packages.
func importPackage(conf *config.Config, name string) {
	switch name {
	case "java.lang":
		conf.Registry.Register("java.lang.Boolean", lang.BooleanClass)
		conf.Registry.Register("java.lang.Byte", lang.ByteClass)
		conf.Registry.Register("java.lang.Character", lang.CharacterClass)
		conf.Registry.Register("java.lang.Double", lang.DoubleClass)
		conf.Registry.Register("java.lang.Float", lang.FloatClass)
		conf.Registry.Register("java.lang.Integer", lang.IntegerClass)
		conf.Registry.Register("java.lang.Long", lang.LongClass)
		conf.Registry.Register("java.lang.Short", lang.ShortClass)
		conf.Registry.Register("java.lang.String", lang.StringClass)
		conf.Registry.Register("java.lang.StringBuilder", lang.StringBuilderClass)
		conf.Registry.Register("java.lang.StringBuffer", lang.StringBufferClass)
	case "java.lang.Boolean":
		conf.Registry.Register("java.lang.Boolean", lang.BooleanClass)
	case "java.lang.Byte":
		conf.Registry.Register("java.lang.Byte", lang.ByteClass)
	case "java.lang.Character":
		conf.Registry.Register("java.lang.Character", lang.CharacterClass)
	case "java.lang.Double":
		conf.Registry.Register("java.lang.Double", lang.DoubleClass)
	case "java.lang.Float":
		conf.Registry.Register("java.lang.Float", lang.FloatClass)
	case "java.lang.Integer":
		conf.Registry.Register("java.lang.Integer", lang.IntegerClass)
	case "java.lang.Long":
		conf.Registry.Register("java.lang.Long", lang.LongClass)
	case "java.lang.Short":
		conf.Registry.Register("java.lang.Short", lang.ShortClass)
	case "java.lang.String":
		conf.Registry.Register("java.lang.String", lang.StringClass)
	case "java.lang.StringBuilder":
		conf.Registry.Register("java.lang.StringBuilder", lang.StringBuilderClass)
	case "java.lang.StringBuffer":
		conf.Registry.Register("java.lang.StringBuffer", lang.StringBufferClass)
	case "java.math":
		conf.Registry.Register("java.lang.Math", math.MathClass)
		conf.Registry.Register("java.math.BigDecimal", math.BigDecimalClass)
		conf.Registry.Register("java.math.BigInteger", math.BigIntegerClass)
	case "java.lang.Math":
		conf.Registry.Register("java.lang.Math", math.MathClass)
	case "java.math.BigDecimal":
		conf.Registry.Register("java.math.BigDecimal", math.BigDecimalClass)
	case "java.math.BigInteger":
		conf.Registry.Register("java.math.BigInteger", math.BigIntegerClass)
	case "java.util":
		conf.Registry.Register("java.util.ArrayList", util.ArrayListClass)
		conf.Registry.Register("java.util.LinkedList", util.LinkedListClass)
		conf.Registry.Register("java.util.HashMap", util.HashMapClass)
		conf.Registry.Register("java.util.TreeMap", util.TreeMapClass)
		conf.Registry.Register("java.util.HashSet", util.HashSetClass)
		conf.Registry.Register("java.util.LinkedHashSet", util.LinkedHashSetClass)
		conf.Registry.Register("java.util.Stack", util.StackClass)
	case "java.util.ArrayList":
		conf.Registry.Register("java.util.ArrayList", util.ArrayListClass)
	case "java.util.LinkedList":
		conf.Registry.Register("java.util.LinkedList", util.LinkedListClass)
	case "java.util.HashMap":
		conf.Registry.Register("java.util.HashMap", util.HashMapClass)
	case "java.util.TreeMap":
		conf.Registry.Register("java.util.TreeMap", util.TreeMapClass)
	case "java.util.HashSet":
		conf.Registry.Register("java.util.HashSet", util.HashSetClass)
	case "java.util.LinkedHashSet":
		conf.Registry.Register("java.util.LinkedHashSet", util.LinkedHashSetClass)
	case "java.util.Stack":
		conf.Registry.Register("java.util.Stack", util.StackClass)
	}
}
