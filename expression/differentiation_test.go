package expression

import (
	"math"
	"testing"
)

func TestDifferentiate(t *testing.T) {
	expr, e := GetExpression("2*x^2")
	if e != nil {
		t.Fatalf("Problem setting up expression, %s", e)
	}
	variables := map[string]float64{
		"x": 10.0,
	}

	expr.SetVariables(variables)

	num, e := expr.Differentiate("x")

	if e != nil {
		t.Fatalf("Differentiation failed, %s", e)
	}

	if num != 40.0 {
		t.Fatalf("Differentiation incorrect, got %f", num)
	}

}

func TestDifferentiateAll(t *testing.T) {

	expr, e := GetExpression("x^2 + 2*y^2 + e^z")
	if e != nil {
		t.Fatalf("Problem setting up expression, %s", e)
	}

	variables := map[string]float64{
		"x": 10.0,
		"y": 1.0,
		"z": 5.0,
	}

	globals := map[string]float64{
		"e": math.E,
	}

	expr.SetGlobals(globals)
	expr.SetVariables(variables)

	expected := map[string]float64{
		"x": 20,
		"y": 4,
		"z": math.Pow(math.E, 5), // this doesn't exactly equal due to numeric differentiation not exact.
	}

	diff, e := expr.DifferentiateAll()
	if e != nil {
		t.Fatalf("Problem setting up expression, %s", e)
	}

	const tolerance float64 = 1000000 // specify a tolerance here because numeric differentiation won't be exact

	for key, value := range diff {
		if math.Round(value*tolerance) != math.Round(expected[key]*tolerance) {
			t.Errorf("For the variable '%s' expected %f but got %f", key, expected[key], value)
		}
	}

}
