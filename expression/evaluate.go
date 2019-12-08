package expression

import (
	"fmt"
	"math"
)

func evaluateShunted(input []token, locals map[string]float64, globals map[string]float64, functions map[string]func(float64) float64) (float64, error) {

	stack := []float64{}

	for _, t := range input {
		if t.kind == tokenOperator {

			operator := t.value.(rune)

			if len(stack) < 2 {
				return 0.0, SyntaxError{
					Description: "Not enough parameters for operator",
				}
			}

			var op1, op2 float64
			op2, stack = stack[len(stack)-1], stack[:len(stack)-1]
			op1, stack = stack[len(stack)-1], stack[:len(stack)-1]

			var result float64
			switch operator {
			case '+':
				result = op1 + op2
			case '-':
				result = op1 - op2
			case '*':
				result = op1 * op2
			case '/':
				result = op1 / op2
			case '^':
				result = math.Pow(op1, op2)
			}

			stack = append(stack, result)

		}

		if t.kind == tokenNumber {
			stack = append(stack, t.value.(float64))
		}

		if t.kind == tokenVariable {
			name := t.value.(string)
			var v float64
			vGlobal, okGlobal := globals[name]
			vLocal, okLocal := locals[name]

			if okGlobal {
				v = vGlobal
			} else if okLocal {
				v = vLocal
			} else {
				return 0.0, SyntaxError{
					Description: fmt.Sprintf("Variable with name %s doesn't exist", name),
				}
			}
			stack = append(stack, v)
		}

		if t.kind == tokenFunction {

			var f func(float64) float64

			if _, ok := t.value.(reservedFunction); ok {
				name := t.value.(reservedFunction)
				switch name {
				case negateFunction:
					f = func(v float64) float64 {
						return -v
					}
				default:
					panic("The internal function is somehow not a valid internal function, this should never happen")
				}
			} else {
				name := t.value.(string)
				var ok bool
				f, ok = functions[name]
				if !ok {
					return 0.0, SyntaxError{
						Description: fmt.Sprintf("Function with name %s doesn't exist", name),
					}
				}
			}

			if len(stack) < 1 {
				return 0.0, SyntaxError{
					Description: "Not enough parameters for function",
				}
			}
			var op1 float64
			op1, stack = stack[len(stack)-1], stack[:len(stack)-1]

			stack = append(stack, f(op1))

		}
	}

	if len(stack) == 0 {
		return 0.0, SyntaxError{
			Description: "There is no value left on the stack, check expression",
		}
	}

	return stack[0], nil

}

// EvalExpression evaluates an expression that has no functions or variables
func EvalExpression(input string) (result float64, err error) {

	tokens, err := tokenizer(input)
	if err != nil {
		return
	}
	shunted, err := shuntingYard(tokens)

	if err != nil {
		return
	}

	result, err = evaluateShunted(shunted, map[string]float64{}, map[string]float64{}, map[string]func(float64) float64{})
	return

}

// GetExpression creates an Expression
func GetExpression(input string) (*Expression, error) {

	tokens, err := tokenizer(input)
	if err != nil {
		return &Expression{}, err
	}
	shunted, err := shuntingYard(tokens)

	if err != nil {
		return &Expression{}, err
	}

	expr := &Expression{
		shunted:   shunted,
		variables: map[string]float64{},
		functions: map[string]func(float64) float64{},
		changedVariables: true,
	}
	expr.Variables()

	return expr, nil

}

// The Expression type contains a parsed expression and the variables and functions that can calculate a value
type Expression struct {
	shunted          []token
	variables        map[string]float64
	globVariables    map[string]float64
	functions        map[string]func(float64) float64
	changedVariables bool
}

// SetVariables sets the variables to be used in the Eval.
// These variables may be used elsewhere, eg. for differentiation.
func (e *Expression) SetVariables(variables map[string]float64) {
	e.variables = variables
	e.changedVariables = true
}

// SetGlobals sets the 'global' variables. This essentially means that these are considered to be 'fundamental constants'.
// These are not touched by differentiation.
// eg. These could be considered to be e and pi.
func (e *Expression) SetGlobals(variables map[string]float64) {
	e.globVariables = variables
	e.changedVariables = true
}

// SetFunctions sets the functions to be used in the Eval
// Functions may only accept a single value and return a single value.
// This limitation may be relaxed in the future, and may accept multiple inputs.
func (e *Expression) SetFunctions(functions map[string]func(float64) float64) {
	e.functions = functions
}

// Eval evaluates the expression with the provided variables and functions
func (e *Expression) Eval() (float64, error) {

	return evaluateShunted(e.shunted, e.variables, e.globVariables, e.functions)
}

// VariableNames returns a map with the keys set to the names of the variables present in the expression.
// It excludes the so called 'Global Variables' from the list.
func (e *Expression) VariableNames() map[string]struct{} {
	names := map[string]struct{}{}
	for _, t := range e.shunted {
		if t.kind == tokenVariable {
			if _, exists := e.globVariables[t.value.(string)]; !exists {
				names[t.value.(string)] = struct{}{}
			}
		}
	}
	return names
}

// Variables gets the variables and the current values, this gets all variables even if they haven't explicitly been set yet.
// Excludes so called 'Global Variables'.
func (e *Expression) Variables() map[string]float64 {

	if e.changedVariables {
		names := e.VariableNames()
		vars := map[string]float64{}

		for n := range names {
			vars[n] = e.variables[n]
		}

		e.variables = vars
		e.changedVariables = false
	}

	return e.variables
}

// FunctionNames returns a map with the keys set to the names of the functions present in the expression.
func (e *Expression) FunctionNames() map[string]struct{} {
	names := map[string]struct{}{}
	for _, t := range e.shunted {
		if t.kind == tokenFunction {
			if n, ok := t.value.(string); ok {
				names[n] = struct{}{}
			}
		}
	}
	return names
}
