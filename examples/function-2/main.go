package main

import (
	"fmt"

	"github.com/harness/go-jexl/jexl"
)

func main() {
	// function defined via context
	env := map[string]any{
		"secrets": map[string]any{
			"getValue": func(args ...any) (any, error) {
				return "dummy-23e4567-e89b-12d3-a456-426614174000", nil
			},
		},
	}

	program, err := jexl.Compile(`secrets.getValue('account.token')`,
		jexl.WithContext(env),
	)
	if err != nil {
		panic(err)
	}

	out, err := jexl.Run(program, env)
	if err != nil {
		panic(err)
	}

	fmt.Println(out) // dummy-23e4567-e89b-12d3-a456-426614174000
}
