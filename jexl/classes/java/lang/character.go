// Copyright (c) 2026 Harness Inc.
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package lang

import (
	"fmt"
	"unicode"

	"github.com/harness/go-jexl/jexl/coerce"
)

// CharacterClass is the java.lang.Character class object.
var CharacterClass characterClass

type characterClass struct{}

func (characterClass) Call(method string, args ...any) (any, error) {
	switch method {
	case "new":
		if len(args) == 0 {
			return NewCharacter(0), nil
		}
		return NewCharacterFrom(args[0]), nil
	case "valueOf":
		if len(args) != 1 {
			return nil, fmt.Errorf("Character.valueOf: expected 1 argument")
		}
		return NewCharacterFrom(args[0]), nil
	case "isDigit":
		if len(args) != 1 {
			return nil, fmt.Errorf("Character.isDigit: expected 1 argument")
		}
		return NewCharacterFrom(args[0]).IsDigit(), nil
	case "isLetter":
		if len(args) != 1 {
			return nil, fmt.Errorf("Character.isLetter: expected 1 argument")
		}
		return NewCharacterFrom(args[0]).IsLetter(), nil
	case "isLetterOrDigit":
		if len(args) != 1 {
			return nil, fmt.Errorf("Character.isLetterOrDigit: expected 1 argument")
		}
		return NewCharacterFrom(args[0]).IsLetterOrDigit(), nil
	case "isUpperCase":
		if len(args) != 1 {
			return nil, fmt.Errorf("Character.isUpperCase: expected 1 argument")
		}
		return NewCharacterFrom(args[0]).IsUpperCase(), nil
	case "isLowerCase":
		if len(args) != 1 {
			return nil, fmt.Errorf("Character.isLowerCase: expected 1 argument")
		}
		return NewCharacterFrom(args[0]).IsLowerCase(), nil
	case "isWhitespace":
		if len(args) != 1 {
			return nil, fmt.Errorf("Character.isWhitespace: expected 1 argument")
		}
		return NewCharacterFrom(args[0]).IsWhitespace(), nil
	case "toUpperCase":
		if len(args) != 1 {
			return nil, fmt.Errorf("Character.toUpperCase: expected 1 argument")
		}
		return NewCharacterFrom(args[0]).ToUpperCase(), nil
	case "toLowerCase":
		if len(args) != 1 {
			return nil, fmt.Errorf("Character.toLowerCase: expected 1 argument")
		}
		return NewCharacterFrom(args[0]).ToLowerCase(), nil
	case "toString":
		if len(args) != 1 {
			return nil, fmt.Errorf("Character.toString: expected 1 argument")
		}
		return NewCharacterFrom(args[0]).ToString(), nil
	case "compare":
		if len(args) != 2 {
			return nil, fmt.Errorf("Character.compare: expected 2 arguments")
		}
		return NewCharacterFrom(args[0]).CompareTo(rune(coerce.ToInt32(args[1]))), nil
	case "MAX_VALUE":
		return Character(unicode.MaxRune), nil
	case "MIN_VALUE":
		return Character(0), nil
	}
	return nil, fmt.Errorf("Character.%s: undefined", method)
}

// Character mirrors java.lang.Character.
type Character rune

// NewCharacter wraps a Go rune as a Character.
func NewCharacter(v rune) Character {
	return Character(v)
}

// NewCharacterFrom coerces any value to a Character.
func NewCharacterFrom(v any) Character {
	if s, ok := v.(string); ok {
		runes := []rune(s)
		if len(runes) > 0 {
			return Character(runes[0])
		}
		return Character(0)
	}
	return Character(rune(coerce.ToInt32(v)))
}

// CharValue returns the rune value.
func (c Character) CharValue() rune {
	return rune(c)
}

// ByteValue returns the code point as int8.
func (c Character) ByteValue() int8 {
	return int8(c)
}

// ShortValue returns the code point as int16.
func (c Character) ShortValue() int16 {
	return int16(c)
}

// IntValue returns the code point as int32.
func (c Character) IntValue() int32 {
	return int32(c)
}

// LongValue returns the code point as int64.
func (c Character) LongValue() int64 {
	return int64(c)
}

// FloatValue returns the code point as float32.
func (c Character) FloatValue() float32 {
	return float32(c)
}

// DoubleValue returns the code point as float64.
func (c Character) DoubleValue() float64 {
	return float64(c)
}

// BooleanValue returns false for the null character, true otherwise.
func (c Character) BooleanValue() bool {
	return c != 0
}

// IsDigit reports whether the character is a digit.
func (c Character) IsDigit() bool {
	return unicode.IsDigit(rune(c))
}

// IsLetter reports whether the character is a letter.
func (c Character) IsLetter() bool {
	return unicode.IsLetter(rune(c))
}

// IsLetterOrDigit reports whether the character is a letter or digit.
func (c Character) IsLetterOrDigit() bool {
	return unicode.IsLetter(rune(c)) || unicode.IsDigit(rune(c))
}

// IsUpperCase reports whether the character is uppercase.
func (c Character) IsUpperCase() bool {
	return unicode.IsUpper(rune(c))
}

// IsLowerCase reports whether the character is lowercase.
func (c Character) IsLowerCase() bool {
	return unicode.IsLower(rune(c))
}

// IsWhitespace reports whether the character is whitespace.
func (c Character) IsWhitespace() bool {
	return unicode.IsSpace(rune(c))
}

// ToUpperCase returns the uppercase version of the character.
func (c Character) ToUpperCase() Character {
	return Character(unicode.ToUpper(rune(c)))
}

// ToLowerCase returns the lowercase version of the character.
func (c Character) ToLowerCase() Character {
	return Character(unicode.ToLower(rune(c)))
}

// CompareTo returns -1, 0, or 1 comparing c to other.
func (c Character) CompareTo(other rune) int {
	switch {
	case rune(c) < other:
		return -1
	case rune(c) > other:
		return 1
	default:
		return 0
	}
}

// Equals reports whether c equals other.
func (c Character) Equals(other rune) bool {
	return rune(c) == other
}

// ToString returns the string representation.
func (c Character) ToString() string {
	return string(rune(c))
}

// Call dispatches instance methods.
func (c Character) Call(method string, args ...any) (any, error) {
	switch method {
	case "charValue":
		return c.CharValue(), nil
	case "byteValue", "toByte":
		return c.ByteValue(), nil
	case "shortValue", "toShort":
		return c.ShortValue(), nil
	case "intValue", "toInteger":
		return c.IntValue(), nil
	case "longValue", "toLong":
		return c.LongValue(), nil
	case "floatValue", "toFloat":
		return c.FloatValue(), nil
	case "doubleValue", "toDouble":
		return c.DoubleValue(), nil
	case "booleanValue", "toBoolean":
		return c.BooleanValue(), nil
	case "isDigit":
		return c.IsDigit(), nil
	case "isLetter":
		return c.IsLetter(), nil
	case "isLetterOrDigit":
		return c.IsLetterOrDigit(), nil
	case "isUpperCase":
		return c.IsUpperCase(), nil
	case "isLowerCase":
		return c.IsLowerCase(), nil
	case "isWhitespace":
		return c.IsWhitespace(), nil
	case "toUpperCase":
		return c.ToUpperCase(), nil
	case "toLowerCase":
		return c.ToLowerCase(), nil
	case "toString":
		return c.ToString(), nil
	case "compareTo":
		if len(args) != 1 {
			return nil, fmt.Errorf("compareTo: expected 1 argument")
		}
		return c.CompareTo(rune(coerce.ToInt32(args[0]))), nil
	case "equals":
		if len(args) != 1 {
			return nil, fmt.Errorf("equals: expected 1 argument")
		}
		return c.Equals(rune(coerce.ToInt32(args[0]))), nil
	case "default":
		return c, nil
	}
	return nil, fmt.Errorf("Character instance: undefined method %q", method)
}
