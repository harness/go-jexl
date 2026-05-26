package main

import (
	"fmt"

	"github.com/harness/go-jexl/jexl"
)

func main() {
	// function defined via function namespace
	program, err := jexl.Compile(`secrets.getValue('account.token')`,
		jexl.WithFunctionNamespace("secrets", "getValue", func(args ...any) (any, error) {
			return "dummy-23e4567-e89b-12d3-a456-426614174000", nil
		}),
	)
	if err != nil {
		panic(err)
	}

	out, err := jexl.Run(program, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(out) // dummy-23e4567-e89b-12d3-a456-426614174000
}
