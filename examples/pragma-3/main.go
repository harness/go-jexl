package main

import (
	"fmt"

	"github.com/harness/go-jexl/jexl"
)

func main() {
	// #pragma jexl.import 'java.util' registers ArrayList, HashMap,
	// HashSet, Stack, etc. without any Go-side setup.
	script := `
		#pragma jexl.import 'java.util'
		var list = new('java.util.ArrayList')
		list.add('a')
		list.add('b')
		list.add('c')
		list.size()
	`

	program, err := jexl.Compile(script)
	if err != nil {
		panic(err)
	}

	out, err := jexl.Run(program, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(out) // 3
}
