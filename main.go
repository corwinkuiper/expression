package main

import (
	"fmt"
	"github.com/corwinkuiper/expression/expression"
	"math"
	"os"
)

func main() {
	exprStr := os.Args[1]

	expr, err := expression.GetExpression(exprStr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	expr.SetFunctions(map[string]func(float64) float64{
		"sin":  math.Sin,
		"cos":  math.Cos,
		"sqrt": math.Sqrt,
		"tan":  math.Tan,
		"ln":   math.Log,
	})

	expr.SetVariables(map[string]float64{
		"pi": math.Pi,
		"e":  math.E,
	})

	result, err := expr.Eval()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Printf("%s = %f", exprStr, result)
}
