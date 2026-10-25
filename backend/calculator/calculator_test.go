package calculator

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

	service := New()
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
		{name: "unsupported operator", expression: "2&3", want: ErrInvalidExpression},
		{name: "implicit multiplication", expression: "2(3)", want: ErrInvalidExpression},
		{name: "non ASCII whitespace", expression: "1\u00a0+2", want: ErrInvalidExpression},
	}

	service := New()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := service.Evaluate(tt.expression)
			if !errors.Is(err, tt.want) {
				t.Errorf("Evaluate(%q) error = %v, want %v", tt.expression, err, tt.want)
			}
		})
	}
}

func TestAdvancedOperations(t *testing.T) {
	cases := []struct{ expression, want string }{
		{"2^3", "8"}, {"2^3^2", "512"}, {"-2^2", "-4"}, {"(-2)^2", "4"},
		{"2^-3", "1/8"}, {"0^2", "0"}, {"0^-1", "0"}, {"0^(1/2)", "0"}, {"0^(-1/2)", "0"}, {"(2-2)^(-3)", "0"}, {"0^9999999999999999999999999", "0"}, {"(-2)^3", "-8"}, {"5^0", "1"}, {"9^(1/2)", "3"}, {"0.25^0.5", "1/2"},
		{"50%", "1/2"}, {"(20+30)%", "1/2"}, {"200*10%", "20"}, {"50%%", "1/200"},
		{"4^50%", "2"}, {"100000000000000000000000000000000^1", "100000000000000000000000000000000"},
	}
	for _, tt := range cases {
		t.Run(tt.expression, func(t *testing.T) {
			got, err := New().Evaluate(tt.expression)
			if err != nil {
				t.Fatal(err)
			}
			want, _ := new(big.Rat).SetString(tt.want)
			if got.Cmp(want) != 0 {
				t.Fatalf("got %v, want %v", got, want)
			}
		})
	}
}

func TestAdvancedErrors(t *testing.T) {
	cases := []struct {
		expression string
		want       error
	}{
		{"0^0", ErrUndefinedValue}, {"0/0", ErrUndefinedValue}, {"(2-2)/(3-3)", ErrUndefinedValue}, {"0^(1-1)", ErrUndefinedValue}, {"0/(-0)", ErrUndefinedValue}, {"(-4)^0.5", ErrInvalidExpression}, {"2^10001", ErrCalculationLimit},
		{"9^9^9", ErrCalculationLimit}, {"1000^(-10000)", ErrCalculationLimit},
		{"(-9)^(1/2)", ErrInvalidExpression}, {"(-8)^(1/3)", ErrInvalidExpression}, {"(-4)^(-1/2)", ErrInvalidExpression}, {"0^(1/0)", ErrDivideByZero}, {"2^", ErrInvalidExpression}, {"%2", ErrInvalidExpression}, {"2%3", ErrInvalidExpression},
		{"2^(1/0)", ErrDivideByZero}, {"2^9999999999999999999999999", ErrCalculationLimit},
		{"2^0.00000000000000000000001", nil},
	}
	for _, tt := range cases {
		t.Run(tt.expression, func(t *testing.T) {
			_, err := New().Evaluate(tt.expression)
			if !errors.Is(err, tt.want) {
				t.Fatalf("got %v, want %v", err, tt.want)
			}
		})
	}
}

func TestPowerBoundaries(t *testing.T) {
	tests := []struct {
		expression string
		want       string
		wantErr    error
	}{
		{expression: "1^10000", want: "1.0000"},
		{expression: "1^(-10000)", want: "1.0000"},
		{expression: "1^10001", wantErr: ErrCalculationLimit},
		{expression: "1^(-10001)", wantErr: ErrCalculationLimit},
		{expression: "2^1024.5", wantErr: ErrCalculationLimit},
		{expression: "2^(-1075.5)", wantErr: ErrCalculationLimit},
		{expression: "2^(1/2)", want: "1.4142"},
		{expression: "0.5^(1/2)", want: "0.7071"},
		{expression: "8^(1/3)", want: "2.0000"},
		{expression: "4^(-0.5)", want: "0.5000"},
	}
	for _, tt := range tests {
		t.Run(tt.expression, func(t *testing.T) {
			got, err := New().Evaluate(tt.expression)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("error = %v, want %v", err, tt.wantErr)
			}
			if tt.wantErr == nil && got.FloatString(4) != tt.want {
				t.Fatalf("result = %s, want %s", got.FloatString(4), tt.want)
			}
		})
	}
}
