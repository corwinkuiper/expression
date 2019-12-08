package expression

import (
	"math"
)

// Differentiates expression given with the variables given over the variable given.
// Not thread safe, in the sense that you shouldn't use the same expression object for other things whilst this is happening.
// It modifies the maps but at the end the map is placed back into it's original form.
// Must have set variables to contain all variables and to the values you want to differentiate at.
func (e *Expression) Differentiate(variable string) (float64, error) {

	epsilon := math.Nextafter(1.0, 2.0) - 1.0

	x := e.variables[variable]
	h := math.Sqrt(epsilon) * x
	if h == 0 {
		h = epsilon
	}

	xForward := x + h
	xBackward := x - h

	dx := xForward - xBackward

	e.variables[variable] = xForward
	result1, err := e.Eval()
	if err != nil {
		return 0, err
	}

	e.variables[variable] = xBackward
	result2, err := e.Eval()
	if err != nil {
		return 0, err
	}

	e.variables[variable] = x
	return (result1 - result2) / dx, nil

}

// DifferentiateAll differentiates all variables in the expression.
// Must have set variables to contain all variables and to the values you want to differentiate at.
func (e *Expression) DifferentiateAll() (map[string]float64, error) {

	resultants := map[string]float64{}

	for key, _ := range e.variables {
		res, err := e.Differentiate(key)
		if err != nil {
			return nil, err
		}
		resultants[key] = res
	}

	return resultants, nil
}
