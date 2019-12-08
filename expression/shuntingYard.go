package expression

func shuntingYard(tokens []token) ([]token, error) {

	var output []token
	var operator []token

	for index, t := range tokens {

		switch t.kind {
		case tokenNumber:
			output = append(output, t)
		case tokenVariable:
			output = append(output, t)
		case tokenFunction:
			operator = append(operator, t)
		case tokenOperator:
			for {
				if len(operator) == 0 {
					break
				}
				topOperator := operator[len(operator)-1]

				fn := topOperator.kind == tokenFunction
				op, equal := false, false

				if topOperator.kind == tokenOperator {
					op = operatorHasPrecedence(topOperator.value.(rune), t.value.(rune))
					equal = operatorEqualPrecedence(topOperator.value.(rune), t.value.(rune))
				}

				left := topOperator.kind == tokenSeparator && topOperator.value.(rune) != '('

				if !(!left && (fn || op || equal)) {
					break
				}

				var to token
				to, operator = operator[len(operator)-1], operator[:len(operator)-1]
				output = append(output, to)

			}

			operator = append(operator, t)
		case tokenSeparator:
			if t.value.(rune) == '(' {
				operator = append(operator, t)
			} else {
				for {
					if len(operator) == 0 {
						return nil, SyntaxError{
							Description: "Mismatched brackets",
							Position:    index,
						}
					}

					var to token
					to, operator = operator[len(operator)-1], operator[:len(operator)-1]

					if to.kind == tokenSeparator {
						break
					}
					output = append(output, to)

				}
			}

		}

	}

	for len(operator) != 0 {
		var to token
		to, operator = operator[len(operator)-1], operator[:len(operator)-1]

		if to.kind == tokenSeparator {
			return nil, SyntaxError{
				Description: "Mismatched brackets, guessing position",
				Position:    len(output),
			}
		}

		output = append(output, to)
	}

	return output, nil

}

func operatorHasPrecedence(op1 rune, op2 rune) bool {

	if op1 == '+' || op1 == '-' {
		return false
	}

	if op1 == '*' || op1 == '/' {
		return op2 == '+' || op2 == '-'
	}

	if op1 == '^' {
		return op2 != '^'
	}

	panic("The operator was an unrecognised type, this is a bug in the tokenizer.")
}

func operatorEqualPrecedence(op1 rune, op2 rune) bool {

	if op1 == '+' || op1 == '-' {
		return op2 == '+' || op2 == '-'
	}

	if op1 == '*' || op1 == '/' {
		return op2 == '*' || op2 == '/'
	}

	if op1 == '^' {
		return op2 == '^'
	}
	panic("The operator was an unrecognised type, this is a bug in the tokenizer.")

}
