package main

import (
	"fmt"

	"github.com/harness/go-jexl/jexl"
	"github.com/harness/go-jexl/jexl/classes/java/lang"
)

func main() {
	script := `
		var sb = new('java.lang.StringBuilder', 'hello')
		sb.append(' ').append('world').toString()
	`

	program, err := jexl.Compile(script,
		jexl.WithClass("java.lang.StringBuilder", lang.StringBuilderClass),
	)
	if err != nil {
		panic(err)
	}

	out, err := jexl.Run(program, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(out) // hello world
}
