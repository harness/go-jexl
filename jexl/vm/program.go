// Copyright (c) 2026 Harness Inc.
// Copyright (c) 2018 Anton Medvedev
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package vm

import (
	"github.com/harness/go-jexl/jexl/classes"
	"github.com/harness/go-jexl/jexl/token"
)

// Program represents a compiled expression.
type Program struct {
	Bytecode      []Opcode
	Arguments     []int
	Constants     []any
	Registry      *classes.Registry
	MaxIterations uint
	MaxMemory     uint

	source    token.File
	locations []token.Range
	variables int
	functions []Function
	debugInfo map[string]string
}

// NewProgram returns a new Program. It's used by the compiler.
func NewProgram(
	source token.File,
	locations []token.Range,
	variables int,
	constants []any,
	bytecode []Opcode,
	arguments []int,
	functions []Function,
	debugInfo map[string]string,
) *Program {
	return &Program{
		source:    source,
		locations: locations,
		variables: variables,
		Constants: constants,
		Bytecode:  bytecode,
		Arguments: arguments,
		functions: functions,
		debugInfo: debugInfo,
	}
}
