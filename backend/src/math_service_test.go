package math

import (
	"errors"
	"math/big"
	"strings"
	"testing"
)

func TestEvaluate(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		want       string
	}{
		{name: "addition", expression: "2+3", want: "5"},
		{name: "subtraction", expression: "2-5", want: "-3"},
		{name: "multiplication", expression: "-2*3.5", want: "-7"},
		{name: "integer literal division produces decimal", expression: "10/4", want: "2.5"},
		{name: "operator precedence", expression: "2+3*4", want: "14"},
		{name: "parentheses", expression: "(2+3)*4", want: "20"},
		{name: "nested parentheses", expression: "8/(1+(1+2))", want: "2"},
		{name: "left associative subtraction", expression: "10-3-2", want: "5"},
		{name: "left associative division", expression: "20/5/2", want: "2"},
		{name: "unary operators", expression: "-(2+3)++8", want: "3"},
		{name: "negative zero is normalized", expression: "-0", want: "0"},
		{name: "exact repeating fraction", expression: "1/3", want: "1/3"},
		{name: "no intermediate rounding", expression: "(1/3)*3000000", want: "1000000"},
		{name: "large integer addition", expression: "9007199254740992+1", want: "9007199254740993"},
		{name: "exact decimal addition", expression: "0.1+0.2", want: "0.3"},
		{name: "cancellation preserves fraction", expression: "9999999999999.9999-9999999999999", want: "0.9999"},
		{name: "leading zeros stay decimal", expression: "0012.50+000.50", want: "13"},
		{name: "maximum input length", expression: strings.Repeat("0", MaxExpressionLength), want: "0"},
	}

	service := NewMathService()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := service.Evaluate(tt.expression)
			if err != nil {
				t.Fatalf("Evaluate(%q) returned error: %v", tt.expression, err)
			}
			want, ok := new(big.Rat).SetString(tt.want)
			if !ok {
				t.Fatalf("Invalid expected rational %q", tt.want)
			}
			if got.Cmp(want) != 0 {
				t.Errorf("Evaluate(%q) = %v, want %v", tt.expression, got, tt.want)
			}
		})
	}
}

func TestEvaluateErrors(t *testing.T) {
	tests := []struct {
		name       string
		expression string
		want       error
	}{
		{name: "division by zero", expression: "1/0", want: ErrDivideByZero},
		{name: "division by calculated zero", expression: "1/(2-2)", want: ErrDivideByZero},
		{name: "division by exact decimal zero", expression: "1/(0.1+0.2-0.3)", want: ErrDivideByZero},
		{name: "too long", expression: strings.Repeat("0", MaxExpressionLength+1), want: ErrInvalidExpression},
		{name: "power syntax unsupported", expression: "10**32", want: ErrInvalidExpression},
		{name: "invalid negative operand", expression: "-()", want: ErrInvalidExpression},
		{name: "empty", expression: "", want: ErrInvalidExpression},
		{name: "whitespace only", expression: "   ", want: ErrInvalidExpression},
		{name: "missing operand", expression: "1+", want: ErrInvalidExpression},
		{name: "missing operator", expression: "1 2", want: ErrInvalidExpression},
		{name: "unclosed parenthesis", expression: "(1+2", want: ErrInvalidExpression},
		{name: "unexpected closing parenthesis", expression: "1+2)", want: ErrInvalidExpression},
		{name: "empty parentheses", expression: "()", want: ErrInvalidExpression},
		{name: "decimal missing integer digits", expression: ".5+1", want: ErrInvalidExpression},
		{name: "decimal missing fractional digits", expression: "5.+1", want: ErrInvalidExpression},
		{name: "scientific notation", expression: "1e2", want: ErrInvalidExpression},
		{name: "unsupported operator", expression: "2^3", want: ErrInvalidExpression},
		{name: "implicit multiplication", expression: "2(3)", want: ErrInvalidExpression},
		{name: "non ASCII whitespace", expression: "1\u00a0+2", want: ErrInvalidExpression},
	}

	service := NewMathService()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Evaluate(tt.expression)
			if !errors.Is(err, tt.want) {
				t.Errorf("Evaluate(%q) error = %v, want %v", tt.expression, err, tt.want)
			}
		})
	}
}
