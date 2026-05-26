package main

import (
	"fmt"

	"github.com/harness/go-jexl/jexl"
	javalang "github.com/harness/go-jexl/jexl/classes/java/lang"
)

func main() {
	// Call a static method on java.lang.Boolean using dot notation.
	program, err := jexl.Compile(`java.lang.Boolean.logicalAnd(true, false)`,
		jexl.WithClass("java.lang.Boolean", javalang.BooleanClass),
	)
	if err != nil {
		panic(err)
	}

	out, err := jexl.Run(program, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(out) // false
}
