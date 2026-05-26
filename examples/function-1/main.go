package main

import (
	"fmt"
	"strings"

	"github.com/harness/go-jexl/jexl"
)

func main() {
	shout := func(args ...any) (any, error) {
		s := fmt.Sprint(args[0])
		return strings.ToUpper(s) + "!", nil
	}

	program, err := jexl.Compile(`shout(name)`,
		jexl.WithFunction("shout", shout, func(string) string { return "" }),
		jexl.WithContext(map[string]any{"name": "alice"}),
	)
	if err != nil {
		panic(err)
	}

	out, err := jexl.Run(program, map[string]any{"name": "alice"})
	if err != nil {
		panic(err)
	}

	fmt.Println(out) // ALICE!
}
