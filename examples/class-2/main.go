package main

import (
	"fmt"

	"github.com/harness/go-jexl/jexl"
	"github.com/harness/go-jexl/jexl/classes/java/util"
)

func main() {
	script := `
		var m = new('java.util.HashMap')
		m.put('name', 'Alice')
		m.put('role', 'admin')
		m.get('name') + ' is ' + m.get('role')
	`

	program, err := jexl.Compile(script,
		jexl.WithClass("java.util.HashMap", util.HashMapClass),
	)
	if err != nil {
		panic(err)
	}

	out, err := jexl.Run(program, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(out) // Alice is admin
}
