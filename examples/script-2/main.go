package main

import (
	"fmt"

	"github.com/harness/go-jexl/jexl"
)

func main() {
	env := map[string]any{
		"items": []any{1, 2, 3, 4, 5},
	}

	script := `
		var total = 0
		for (item : items) {
			total += item
		}
		return total
	`

	out, err := jexl.Eval(script, env)
	if err != nil {
		panic(err)
	}

	fmt.Println(out) // 15
}
