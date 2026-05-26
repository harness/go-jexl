// Copyright (c) 2026 Harness Inc.
// Copyright (c) 2018 Anton Medvedev
// Use of this source code is governed by the MIT license
// that can be found in the LICENSE file.

package jexl

import (
	"github.com/harness/go-jexl/jexl/checker"
	"github.com/harness/go-jexl/jexl/compiler"
	"github.com/harness/go-jexl/jexl/config"
	"github.com/harness/go-jexl/jexl/parser"
	"github.com/harness/go-jexl/jexl/pragma"
	"github.com/harness/go-jexl/jexl/vm"
)

// Compile parses and compiles given input expression to bytecode program.
func Compile(input string, ops ...Option) (*vm.Program, error) {
	config := config.New()
	for _, op := range ops {
		op(config)
	}

	input, directives := pragma.Parse(input)
	pragma.Apply(config, directives)

	tree, err := checker.ParseCheck(input, config)
	if err != nil {
		return nil, err
	}

	program, err := compiler.Compile(tree, config)
	if err != nil {
		return nil, err
	}

	return program, nil
}

// Run evaluates given bytecode program.
func Run(program *vm.Program, env any) (any, error) {
	return vm.Run(program, env)
}

// Eval parses, compiles and runs given input.
func Eval(input string, env any) (any, error) {
	tree, err := parser.Parse(input)
	if err != nil {
		return nil, err
	}

	program, err := compiler.Compile(tree, nil)
	if err != nil {
		return nil, err
	}

	output, err := Run(program, env)
	if err != nil {
		return nil, err
	}

	return output, nil
}
