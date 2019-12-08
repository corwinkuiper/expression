package expression

import (
	"strconv"
	"unicode"
)

func tokenizer(input string) ([]token, error) {

	var tokens []token

	var numberBuffer []rune
	var letterBuffer []rune

	interruptedToken := false

	for index, character := range input {

		if character == ' ' {
			if len(letterBuffer) != 0 || len(numberBuffer) != 0 {
				interruptedToken = true
			}
			continue
		}

		if interruptedToken && (isNumber(character) || isLetter(character)) {
			return nil, SyntaxError{
				Description: "Interrupted token",
				Position:    index,
			}
		}

		if isNumber(character) {
			if len(letterBuffer) != 0 {
				return nil, SyntaxError{
					Description: "Number following letter",
					Position:    index,
				}
			}

			numberBuffer = append(numberBuffer, character)

		} else if isLetter(character) {

			if len(numberBuffer) != 0 {
				return nil, SyntaxError{
					Description: "Letter following number",
					Position:    index,
				}
			}

			letterBuffer = append(letterBuffer, character)

		} else if isOperator(character) {
			interruptedToken = false
			// flush buffers
			if len(numberBuffer) != 0 {
				num, err := toNumber(numberBuffer)
				if err != nil {
					return nil, err
				}

				tokens = append(tokens, token{
					kind:  tokenNumber,
					value: num,
				})
				numberBuffer = []rune{}
			}
			if len(letterBuffer) != 0 {
				tokens = append(tokens, token{
					kind:  tokenVariable,
					value: toString(letterBuffer),
				})
				letterBuffer = []rune{}
			}

			// check if it's a negation rather than a minus sign
			if character == '-' {
				if len(tokens) == 0 {
					tokens = append(tokens, token{
						kind:  tokenFunction,
						value: negateFunction,
					})
					continue
				} else {
					leftNum := tokens[len(tokens)-1].kind == tokenNumber
					leftVar := tokens[len(tokens)-1].kind == tokenVariable
					leftRightBracket := tokens[len(tokens)-1].kind == tokenSeparator && tokens[len(tokens)-1].value.(rune) == ')'
					if !leftNum && !leftVar && !leftRightBracket {
						tokens = append(tokens, token{
							kind:  tokenFunction,
							value: negateFunction,
						})
						continue
					}
				}
			}

			tokens = append(tokens, token{
				kind:  tokenOperator,
				value: character,
			})

		} else if isSeparator(character) {
			interruptedToken = false
			if character == '(' {
				if len(numberBuffer) != 0 {
					return nil, SyntaxError{
						Description: "Number before opening bracket",
						Position:    index,
					}
				}

				if len(letterBuffer) != 0 {
					tokens = append(tokens, token{
						kind:  tokenFunction,
						value: string(letterBuffer),
					})
					letterBuffer = []rune{}
				}
			} else if character == ')' {
				if len(numberBuffer) != 0 {
					num, err := toNumber(numberBuffer)
					if err != nil {
						return nil, err
					}

					tokens = append(tokens, token{
						kind:  tokenNumber,
						value: num,
					})
					numberBuffer = []rune{}
				}
				if len(letterBuffer) != 0 {
					tokens = append(tokens, token{
						kind:  tokenVariable,
						value: toString(letterBuffer),
					})
					letterBuffer = []rune{}
				}
			}
			tokens = append(tokens, token{
				kind:  tokenSeparator,
				value: character,
			})
		}

	}

	if len(numberBuffer) != 0 {
		num, err := toNumber(numberBuffer)
		if err != nil {
			return nil, err
		}

		tokens = append(tokens, token{
			kind:  tokenNumber,
			value: num,
		})
		numberBuffer = []rune{}
	}
	if len(letterBuffer) != 0 {
		tokens = append(tokens, token{
			kind:  tokenVariable,
			value: toString(letterBuffer),
		})
		letterBuffer = []rune{}
	}

	return tokens, nil

}

func toNumber(chars []rune) (float64, error) {
	return strconv.ParseFloat(string(chars), 64)
}

func toString(chars []rune) string {
	return string(chars)
}

func isNumber(char rune) bool {
	return unicode.IsDigit(char) || char == '.'
}

func isLetter(char rune) bool {
	return unicode.IsLetter(char)
}

func isOperator(char rune) bool {
	return char == '+' || char == '-' || char == '*' || char == '/' || char == '^'
}

func isSeparator(char rune) bool {
	return char == '(' || char == ')'
}
