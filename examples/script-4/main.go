package main

import (
	"errors"
	"fmt"

	"github.com/harness/go-jexl/jexl"
)

func main() {
	env := map[string]any{
		"price":    49,
		"quantity": 3,
	}

	// this is a complex expression that is not
	// eligible for fast-path evaluation
	expr := "price * quantity"

	// we attempt to execute via fast path
	out, err := jexl.EvalPath(expr, env)

	// if we detect the sentinel error we should eval
	// using the standard jexl compiler / evaluator
	if errors.Is(err, jexl.ErrNotPropertyPath) {
		out, err = jexl.Eval(expr, env)
	}

	fmt.Println(out) // 147
}
