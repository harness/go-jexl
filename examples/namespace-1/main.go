package main

import (
	"fmt"

	"github.com/harness/go-jexl/jexl"
	"github.com/harness/go-jexl/jexl/classes/java/math"
)

func main() {
	program, err := jexl.Compile(`Math:sqrt(n)`,
		jexl.WithNamespace("Math", math.MathClass),
		jexl.WithContext(map[string]any{"n": 0}),
	)
	if err != nil {
		panic(err)
	}

	out, err := jexl.Run(program, map[string]any{"n": 144})
	if err != nil {
		panic(err)
	}

	fmt.Println(out) // 12
}
