package expression

import "testing"

func TestTokeniser(t *testing.T) {
	expression := "2 * 3 * (2 + 4)"
	_, err := tokenizer(expression)

	if err != nil {
		t.Fatalf("Tokeniser errored with message %s", err)
	}
}
