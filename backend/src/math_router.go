package math

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

const maxResultIntegerDigits = 33

var errResultTooLarge = errors.New("Result cannot exceed 33 integer digits")

type handleMathReturn struct {
	Result string `json:"result"`
}

type mathRouter struct{}

func NewMathRouter() *mathRouter {
	return &mathRouter{}
}

func (mr *mathRouter) Router(r chi.Router) {
	r.Get("/", mr.handleMath)
}

func (mr *mathRouter) handleMath(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	encodedQueries := queryParams["query"]
	if len(encodedQueries) != 1 {
		writeMathError(w, "Exactly one 'query' parameter is required", http.StatusBadRequest)
		return
	}
	encodedQuery := encodedQueries[0]

	if encodedQuery == "" {
		writeMathError(w, "Missing required 'query' parameter", http.StatusBadRequest)
		return
	}

	decodedQuery, err := base64.StdEncoding.DecodeString(encodedQuery)
	if err != nil {
		writeMathError(w, "Invalid Base64 query", http.StatusBadRequest)
		return
	}

	query := string(decodedQuery)

	if len(query) > MaxExpressionLength {
		writeMathError(w, "Query parameter cannot be larger than 64 characters", http.StatusBadRequest)
		return
	}

	mathService := NewMathService()

	result, err := mathService.Evaluate(query)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, ErrDivideByZero) {
			status = http.StatusUnprocessableEntity
		}
		writeMathError(w, err.Error(), status)
		return
	}

	formatted, err := formatResult(result)
	if err != nil {
		writeMathError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&handleMathReturn{
		Result: formatted,
	})
}

func formatResult(result *big.Rat) (string, error) {
	// Round only the final exact value, with halves rounded away from zero.
	formatted := result.FloatString(4)
	formatted = strings.TrimRight(formatted, "0")
	formatted = strings.TrimRight(formatted, ".")

	if formatted == "-0" {
		formatted = "0"
	}
	integerPart, _, _ := strings.Cut(strings.TrimPrefix(formatted, "-"), ".")
	if len(integerPart) > maxResultIntegerDigits {
		return "", errResultTooLarge
	}
	return formatted, nil
}

func writeMathError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
