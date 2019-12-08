package fit

import (
	"github.com/corwinkuiper/expression/expression"
	"math"
	"testing"
)

func TestLeastSquares(t *testing.T) {

	const accuracy = 10000

	data := []DataElement{
		{
			X: 0,
			Y: 1,
		},
		{
			X: 1,
			Y: 3,
		}, {
			X: 2,
			Y: 5,
		},
		{
			X: 3,
			Y: 7,
		}, {
			X: 4,
			Y: 9,
		},
	}

	expr, e := expression.GetExpression("A*x + B")

	if e != nil {
		t.Fatalf("Could not get expression, %s", e)
	}

	a, e := LeastSquares(data, expr)
	if e != nil {
		t.Fatalf("Fit errored, %s", e)
	}

	if math.Round(a["A"]*accuracy)/accuracy != 2 {
		t.Fatalf("Incorrect value of A calculated, expected 2 but got %f", a["A"])
	}
	if math.Round(a["B"]*accuracy)/accuracy != 1 {
		t.Fatalf("Incorrect value of B calculated, expected 1 but got %f", a["A"])
	}
}
