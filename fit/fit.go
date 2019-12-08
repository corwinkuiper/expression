package fit

import (
	"github.com/corwinkuiper/expression/expression"
	"math"
)

const dependent = "x"
const iterationCount = 5000

type DataElement struct {
	X float64
	Y float64
}

func LeastSquares(data []DataElement, expr *expression.Expression) (map[string]float64, error) {

	VariableValues := expr.Variables()

	// want to differentiate leastSquaresResidualFunction on every variable except "dependent".

	for i := 0; i < iterationCount; i++ {
		differentials := map[string]float64{}
		for key := range VariableValues {
			if key == dependent {
				continue
			}
			r, e := differentiateLeastSquaresFunction(data, expr, key)
			if e != nil {
				return nil, e
			}

			differentials[key] = r
		}

		for key, diff := range differentials {
			//fmt.Println(key, VariableValues[key], diff)
			VariableValues[key] = VariableValues[key] - diff / 100
		}
	}

	fitVars := map[string]float64{}

	for key, value := range VariableValues {
		if key == dependent {
			continue
		}
		fitVars[key] = value
	}

	return fitVars, nil
}

func differentiateLeastSquaresFunction(data []DataElement, expr *expression.Expression, variableName string) (float64, error) {

	vars := expr.Variables()
	epsilon := math.Nextafter(1.0, 2.0) - 1.0

	x := vars[variableName]
	h := 0.000001

	if h == 0 {
		h = epsilon
	}

	xForward := x + h
	xBackward := x - h
	dx := xForward - xBackward

	vars[variableName] = xForward
	result1, err := leastSquaresResidualFunction(data, expr)
	if err != nil {
		return 0, err
	}

	vars[variableName] = xBackward
	result2, err := leastSquaresResidualFunction(data, expr)
	if err != nil {
		return 0, err
	}

	vars[variableName] = x

	return (result2 - result1) / dx, nil

}

func leastSquaresResidualFunction(data []DataElement, expr *expression.Expression) (float64, error) {

	v := expr.Variables()

	x := v[dependent]

	sum := 0.0

	for _, d := range data {
		v[dependent] = d.X
		y, err := expr.Eval()
		if err != nil {
			return 0, err
		}
		sum += math.Pow(d.Y-y, 2)
	}

	v[dependent] = x

	return -sum, nil
}
