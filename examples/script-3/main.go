package main

import (
	"fmt"

	"github.com/harness/go-jexl/jexl"
)

func main() {
	env := map[string]any{
		"trigger": map[string]any{
			"payload": map[string]any{
				"crNumber": "CR-12345",
			},
		},
	}

	// zero-allocation fast evaluation for path
	// expressions.
	out, err := jexl.EvalPath("trigger.payload.crNumber", env)
	if err != nil {
		panic(err)
	}

	fmt.Println(out) // CR-12345
}
