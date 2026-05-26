package main

import (
	"fmt"

	"github.com/harness/go-jexl/jexl"
)

func main() {
	env := map[string]any{
		"name": "Alice",
		"age":  30,
	}

	out, err := jexl.Eval(`name + " is " + age`, env)
	if err != nil {
		panic(err)
	}

	fmt.Println(out) // Alice is 30
}
