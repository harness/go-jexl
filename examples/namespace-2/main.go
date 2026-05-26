package main

import (
	"fmt"
	"strings"

	"github.com/harness/go-jexl/jexl"
)

func main() {
	program, err := jexl.Compile(`Strings:reverse(word)`,
		jexl.WithNamespace("Strings", stringUtils{}),
		jexl.WithContext(map[string]any{"word": ""}),
	)
	if err != nil {
		panic(err)
	}

	out, err := jexl.Run(program, map[string]any{"word": "hello"})
	if err != nil {
		panic(err)
	}

	fmt.Println(out) // olleh
}

//
// Custom class implementation
//

type stringUtils struct{}

func (stringUtils) Call(method string, args ...any) (any, error) {
	switch method {
	case "reverse":
		s := fmt.Sprint(args[0])
		r := []rune(s)
		for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
			r[i], r[j] = r[j], r[i]
		}
		return string(r), nil
	case "shout":
		return strings.ToUpper(fmt.Sprint(args[0])) + "!", nil
	}
	return nil, fmt.Errorf("unknown method: %s", method)
}
