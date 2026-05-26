package main

import (
	"fmt"
	"math"

	"github.com/harness/go-jexl/jexl"
)

func main() {
	program, err := jexl.Compile(`new('Point', 3.0, 4.0).distance()`,
		jexl.WithClass("Point", pointClass{}),
	)
	if err != nil {
		panic(err)
	}

	out, err := jexl.Run(program, nil)
	if err != nil {
		panic(err)
	}

	fmt.Println(out) // 5
}

//
// Custom class
//

type pointClass struct{}

type point struct{ x, y float64 }

func (p point) Call(method string, args ...any) (any, error) {
	switch method {
	case "distance":
		return math.Sqrt(p.x*p.x + p.y*p.y), nil
	}
	return nil, fmt.Errorf("unknown method: %s", method)
}

func (pointClass) Call(method string, args ...any) (any, error) {
	switch method {
	case "new":
		x, _ := args[0].(float64)
		y, _ := args[1].(float64)
		return &point{x, y}, nil
	}
	return nil, fmt.Errorf("unknown method: %s", method)
}
