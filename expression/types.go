package expression

import "fmt"

//SyntaxError is the error type if the expression contains a syntactical error.
type SyntaxError struct {
	Position    int
	Description string
}

func (e SyntaxError) Error() string {
	if e.Position != 0 {
		return fmt.Sprintf(`"%s" in position %d`, e.Description, e.Position)
	}
	return e.Description
}

type token struct {
	kind  tokenKind
	value interface{}
}

type tokenKind int

const (
	tokenVariable tokenKind = iota
	tokenFunction
	tokenOperator
	tokenSeparator
	tokenNumber
)

type reservedFunction int

const (
	negateFunction reservedFunction = iota
)
