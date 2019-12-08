package expression

import (
	"math"
	"testing"
)

func TestEvaluateSuccess(t *testing.T) {

	expression := "2 * 3 + 2*(5 + 4 / 2) ^ 2"
	expectedResult := 104.0

	result, err := EvalExpression(expression)

	if err != nil {
		t.Fatalf("Unexpected Error\n%s", err)
	}

	if result != expectedResult {
		t.Fatalf("Unexpected answer, expected %f but got %f", expectedResult, result)
	}

}

func TestEvaluateFailNumberLetter(t *testing.T) {

	expression := "2x"
	_, err := EvalExpression(expression)
	if err == nil {
		t.Fatal("Expected an error but got nothing.")
	}

}

func TestEvaluateFailLetterNumber(t *testing.T) {

	expression := "x2"
	_, err := EvalExpression(expression)
	if err == nil {
		t.Fatal("Expected an error but got nothing.")
	}

}

func TestEvaluateFailNumberBracket(t *testing.T) {

	expression := "2("
	_, err := EvalExpression(expression)
	if err == nil {
		t.Fatal("Expected an error but got nothing.")
	}

}

func TestEvaluateFailMismatchBracket(t *testing.T) {

	expression := "2*("
	_, err := EvalExpression(expression)
	if err == nil {
		t.Fatal("Expected an error but got nothing.")
	}

}

func TestEvaluateFailMismatchBracketRight(t *testing.T) {

	expression := "2*)"
	_, err := EvalExpression(expression)
	if err == nil {
		t.Fatal("Expected an error but got nothing.")
	}

}

func TestEvaluateVariablesNoneSupplied(t *testing.T) {

	expression := "2 * variable"
	_, err := EvalExpression(expression)
	if err == nil {
		t.Fatal("Expected an error but got nothing.")
	}
}

func TestEvaluateWithVariables(t *testing.T) {
	variables := map[string]float64{
		"coefficient": 5.0,
		"constant":    3.0,
	}

	expression := "2 * coefficient + constant"
	expected := 13.0

	e, err := GetExpression(expression)
	if err != nil {
		t.Fatalf("Unexpected Error\n%s", err)
	}

	e.SetVariables(variables)

	res, err := e.Eval()

	if err != nil {
		t.Fatalf("Unexpected Error\n%s", err)
	}

	if res != expected {
		t.Fatalf("Expected %f but got %f", expected, res)
	}
}

func TestBadExpressionInterrupted(t *testing.T) {
	expression := "2 2"
	_, err := EvalExpression(expression)
	if err == nil {
		t.Fatal("Expected an error but got nothing.")
	}
}

func TestBadExpressionOperatorGalore(t *testing.T) {
	expression := "2 + + 2"
	_, err := EvalExpression(expression)
	if err == nil {
		t.Fatal("Expected an error but got nothing.")
	}
}

func TestNegativeNumber(t *testing.T) {
	expression := "-5"
	expected := -5.0

	result, err := EvalExpression(expression)

	if err != nil {
		t.Fatalf("Unexpected Error\n%s", err)
	}

	if result != expected {
		t.Fatalf("Unexpected answer, expected %f but got %f", expected, result)
	}
}

func TestExpressionWithNegativeNumbersEverywhere(t *testing.T) {
	expression := "-5 + (-6 * -7) / (-5 * 3) - 10/2"
	expected := -12.8

	result, err := EvalExpression(expression)

	if err != nil {
		t.Fatalf("Unexpected Error\n%s", err)
	}

	if result != expected {
		t.Fatalf("Unexpected answer, expected %f but got %f", expected, result)
	}
}

func TestExpressionWithNegativeBrackets(t *testing.T) {
	expression := "-(5+2)"
	expected := -7.0

	result, err := EvalExpression(expression)

	if err != nil {
		t.Fatalf("Unexpected Error\n%s", err)
	}

	if result != expected {
		t.Fatalf("Unexpected answer, expected %f but got %f", expected, result)
	}
}

func TestExpressionWithNegativeFunction(t *testing.T) {
	expression := "-double(x) + - x"
	expected := -7.5

	expr, err := GetExpression(expression)

	if err != nil {
		t.Fatalf("Unexpected Error\n%s", err)
	}

	expr.SetFunctions(map[string]func(float64) float64{
		"double": func(v float64) float64 {
			return v * 2
		},
	})
	expr.SetVariables(map[string]float64{
		"x": 2.5,
	})

	result, err := expr.Eval()

	if result != expected {
		t.Fatalf("Unexpected answer, expected %f but got %f", expected, result)
	}
}

func TestDecimal(t *testing.T) {
	expression := "2.5645"
	expected := 2.5645

	result, err := EvalExpression(expression)

	if err != nil {
		t.Fatalf("Unexpected Error\n%s", err)
	}

	if result != expected {
		t.Fatalf("Unexpected answer, expected %f but got %f", expected, result)
	}

}

func TestVariableNames(t *testing.T) {
	expressionString := "A*x^2 + B^x + Constant"

	expr, e := GetExpression(expressionString)
	if e != nil {
		t.Fatalf("Got unexpected error %s", e)
	}

	names := expr.VariableNames()

	if len(names) != 4 {
		t.Errorf("There are some missing variables: have %d", len(names))
	}
	for key, _ := range names {
		if !(key == "A" || key == "x" || key == "B" || key == "Constant") {
			t.Errorf("Expected that the variables have correct names")
		}
	}

}

func TestFunctionNames(t *testing.T) {
	expressionString := "sin(20) + ln(30)"

	expr, e := GetExpression(expressionString)
	if e != nil {
		t.Fatalf("Got unexpected error %s", e)
	}

	names := expr.FunctionNames()

	if len(names) != 2 {
		t.Errorf("There are some missing functions: have %d", len(names))
	}

	for key, _ := range names {
		if !(key == "sin" || key == "ln") {
			t.Errorf("Expected that the functions have correct names")
		}
	}
}

func TestVariableGlobals(t *testing.T) {

	expressionString := "sin(y) + e^(2*pi*x) + c"

	globals := map[string]float64{
		"e":  math.E,
		"pi": math.Pi,
	}

	expr, e := GetExpression(expressionString)
	if e != nil {
		t.Fatalf("Got unexpected error %s", e)
	}

	expr.SetGlobals(globals)
	names := expr.VariableNames()

	if len(names) != 3 {
		t.Errorf("There are some missing variables: have %d", len(names))
	}

	for key, _ := range names {
		if !(key == "y" || key == "x" || key == "c") {
			t.Errorf("Expected that the variables have correct names")
		}
	}
}
