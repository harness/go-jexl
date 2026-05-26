package main

import (
	"fmt"

	"github.com/harness/go-jexl/jexl"
	"github.com/harness/go-jexl/jexl/template"
)

func main() {
	env := map[string]any{
		"name": "Alice",
		"age":  30,
	}

	tmpl := template.New(func(expr string) (any, error) {
		return jexl.Eval(expr, env)
	}).Delim("${{", "}}")

	out, err := tmpl.ExecString("hello ${{ name }}, you are ${{ age }} years old")
	if err != nil {
		panic(err)
	}

	fmt.Println(out) // hello Alice, you are 30 years old
}
