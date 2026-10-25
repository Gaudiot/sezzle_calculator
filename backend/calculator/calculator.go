package calculator

import (
	"errors"
	"math"
	"math/big"
)

const (
	// MaxExpressionLength bounds parsing depth and exact-arithmetic work.
	MaxExpressionLength     = 64
	maxIntegerExponent      = 10_000
	maxPowerBits            = 65_536
	squareRootPrecisionBits = 256
	squareRootDegree        = 2
)

var (
	ErrUndefinedValue    = errors.New("Undefined value")
	ErrDivideByZero      = errors.New("Division by zero")
	ErrInvalidExpression = errors.New("Invalid expression")
	ErrCalculationLimit  = errors.New("Calculation exceeds supported precision or size")
)

type Calculator struct{}

func New() *Calculator {
	return &Calculator{}
}

func (calculator *Calculator) Evaluate(expression string) (*big.Rat, error) {
	if len(expression) == 0 || len(expression) > MaxExpressionLength {
		return nil, ErrInvalidExpression
	}
	parser := expressionParser{expression: expression}

	result, err := parser.evaluateExpression()
	if err != nil {
		return nil, err
	}

	if parser.position != len(parser.expression) {
		return nil, ErrInvalidExpression
	}

	return result, nil
}

func power(base, exponent *big.Rat) (*big.Rat, error) {
	if base.Sign() == 0 {
		if exponent.Sign() == 0 {
			return nil, ErrUndefinedValue
		}
		// Calculator convention: zero to any nonzero exponent is zero.
		return new(big.Rat), nil
	}
	if exponent.IsInt() {
		return integerPower(base, exponent.Num())
	}
	if base.Sign() < 0 {
		return nil, ErrInvalidExpression
	}
	if exponent.Cmp(big.NewRat(1, squareRootDegree)) == 0 {
		return squareRoot(base), nil
	}
	return fractionalPower(base, exponent)
}

func integerPower(base *big.Rat, exponent *big.Int) (*big.Rat, error) {
	// Bound allocations even for short inputs such as 9^9^9.
	if !exponent.IsInt64() {
		return nil, ErrCalculationLimit
	}
	n := exponent.Int64()
	if n < -maxIntegerExponent || n > maxIntegerExponent {
		return nil, ErrCalculationLimit
	}
	negative := n < 0
	if negative {
		n = -n
	}
	if int64(base.Num().BitLen()+base.Denom().BitLen())*n > maxPowerBits {
		return nil, ErrCalculationLimit
	}
	numerator := new(big.Int).Exp(base.Num(), big.NewInt(n), nil)
	denominator := new(big.Int).Exp(base.Denom(), big.NewInt(n), nil)
	if negative {
		numerator, denominator = denominator, numerator
	}
	return new(big.Rat).SetFrac(numerator, denominator), nil
}

// squareRoot retains high precision and is exact for rational squares.
// The caller must ensure the base is nonnegative.
func squareRoot(base *big.Rat) *big.Rat {
	numerator := new(big.Int).Sqrt(base.Num())
	denominator := new(big.Int).Sqrt(base.Denom())
	numeratorSquared := new(big.Int).Mul(numerator, numerator)
	denominatorSquared := new(big.Int).Mul(denominator, denominator)
	if numeratorSquared.Cmp(base.Num()) == 0 && denominatorSquared.Cmp(base.Denom()) == 0 {
		return new(big.Rat).SetFrac(numerator, denominator)
	}
	value := new(big.Float).SetPrec(squareRootPrecisionBits).SetRat(base)
	result, _ := value.Sqrt(value).Rat(nil)
	return result
}

func fractionalPower(base, exponent *big.Rat) (*big.Rat, error) {
	b, _ := base.Float64()
	e, _ := exponent.Float64()
	value := math.Pow(b, e)
	if math.IsNaN(value) {
		return nil, ErrUndefinedValue
	}
	if math.IsInf(value, 0) || (value == 0 && base.Sign() != 0) {
		return nil, ErrCalculationLimit
	}
	return new(big.Rat).SetFloat64(value), nil
}
