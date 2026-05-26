package main

import (
	"fmt"

	"github.com/harness/go-jexl/jexl"
)

func main() {
	// #pragma jexl.import 'java.lang' registers the java.lang classes
	// (String, StringBuilder, Integer, etc.) without any Go-side setup.
	script := `
		#pragma jexl.import 'java.lang'
		var sb = new('java.lang.StringBuilder', 'hello')
		sb.append(' world').toString()
	`

	program, err := jexl.Compile(script)
	if err != nil {
		panic(err)
	}

	out, err := jexl.Run(program, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(out) // hello world
}
