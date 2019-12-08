package expression

import "testing"

func BenchmarkExpression(b *testing.B) {

	expresionString := "2 * 10 + 20 * (4 + 5) * 100/20 + 20/40"
	expected := 920.5

	expression, _ := GetExpression(expresionString)

	for i := 0; i < b.N; i++ {
		r, _ := expression.Eval()
		if r != expected {
			b.Fatalf("Expected %f but got %f", expected, r)
		}
	}
}