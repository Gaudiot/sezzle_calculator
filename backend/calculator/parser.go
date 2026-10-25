package calculator

import "math/big"

type expressionParser struct {
	expression string
	position   int
}

// evaluateExpression handles addition and subtraction.
func (parser *expressionParser) evaluateExpression() (*big.Rat, error) {
	result, err := parser.evaluateTerm()
	if err != nil {
		return nil, err
	}

	for {
		operator := parser.nextOperator('+', '-')
		if operator == 0 {
			return result, nil
		}

		right, err := parser.evaluateTerm()
		if err != nil {
			return nil, err
		}

		if operator == '+' {
			result.Add(result, right)
		} else {
			result.Sub(result, right)
		}
	}
}

// evaluateTerm handles multiplication and division.
func (parser *expressionParser) evaluateTerm() (*big.Rat, error) {
	result, err := parser.evaluateUnary()
	if err != nil {
		return nil, err
	}

	for {
		operator := parser.nextOperator('*', '/')
		if operator == 0 {
			return result, nil
		}

		right, err := parser.evaluateUnary()
		if err != nil {
			return nil, err
		}
		if operator == '/' && right.Sign() == 0 {
			if result.Sign() == 0 {
				return nil, ErrUndefinedValue
			}
			return nil, ErrDivideByZero
		}

		if operator == '*' {
			result.Mul(result, right)
		} else {
			result.Quo(result, right)
		}
	}
}

func (parser *expressionParser) evaluateUnary() (*big.Rat, error) {
	if parser.consume('+') {
		return parser.evaluateUnary()
	}
	if parser.consume('-') {
		value, err := parser.evaluateUnary()
		if err != nil {
			return nil, err
		}
		return value.Neg(value), nil
	}

	return parser.evaluatePower()
}

// Powers associate right to left and bind more tightly than unary signs.
func (parser *expressionParser) evaluatePower() (*big.Rat, error) {
	base, err := parser.evaluateValue()
	if err != nil {
		return nil, err
	}
	for parser.consume('%') {
		base.Quo(base, big.NewRat(100, 1))
	}
	if !parser.consume('^') {
		return base, nil
	}
	exponent, err := parser.evaluateUnary()
	if err != nil {
		return nil, err
	}
	return power(base, exponent)
}

func (parser *expressionParser) evaluateValue() (*big.Rat, error) {
	if parser.consume('(') {
		value, err := parser.evaluateExpression()
		if err != nil {
			return nil, err
		}

		if !parser.consume(')') {
			return nil, ErrInvalidExpression
		}

		return value, nil
	}

	return parser.getNumber()
}

func (parser *expressionParser) getNumber() (*big.Rat, error) {
	start := parser.position

	for parser.position < len(parser.expression) && isDigit(parser.expression[parser.position]) {
		parser.position++
	}
	if parser.position == start {
		return nil, ErrInvalidExpression
	}

	if parser.consume('.') {
		fractionStart := parser.position
		for parser.position < len(parser.expression) && isDigit(parser.expression[parser.position]) {
			parser.position++
		}
		if parser.position == fractionStart {
			return nil, ErrInvalidExpression
		}
	}

	// The parser restricts tokens to decimal notation before SetString is called.
	value, ok := new(big.Rat).SetString(parser.expression[start:parser.position])
	if !ok {
		return nil, ErrInvalidExpression
	}

	return value, nil
}

func (parser *expressionParser) nextOperator(first, second byte) byte {
	if parser.position >= len(parser.expression) {
		return 0
	}

	operator := parser.expression[parser.position]
	if operator != first && operator != second {
		return 0
	}

	parser.position++
	return operator
}

func (parser *expressionParser) consume(character byte) bool {
	if parser.position >= len(parser.expression) || parser.expression[parser.position] != character {
		return false
	}

	parser.position++
	return true
}

func isDigit(character byte) bool {
	return character >= '0' && character <= '9'
}
