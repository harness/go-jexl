package main

import (
	"fmt"

	"github.com/harness/go-jexl/jexl"
)

func main() {
	// #pragma jexl.options 'strict' causes the compiler to error
	// on undefined variables instead of treating them as null.
	script := `
		#pragma jexl.options 'strict'
		name + " is " + age
	`

	_, err := jexl.Compile(script)
	if err != nil {
		fmt.Println("compile error:", err)
		// compile error: undefined: name (1:3)
		return
	}
}
