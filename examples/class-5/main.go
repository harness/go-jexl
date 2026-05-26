package main

import (
	"fmt"

	"github.com/harness/go-jexl/jexl"
	javalang "github.com/harness/go-jexl/jexl/classes/java/lang"
)

func main() {
	// Construct a java.lang.Integer instance with new and call instance methods.
	script := `
		var n = new('java.lang.Integer', 42)
		n.compareTo(100)
	`

	program, err := jexl.Compile(script,
		jexl.WithClass("java.lang.Integer", javalang.IntegerClass),
	)
	if err != nil {
		panic(err)
	}

	out, err := jexl.Run(program, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(out) // -1
}
