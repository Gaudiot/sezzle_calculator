package httpapi

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"math/big"
	"net/http"
	"strings"

	"github.com/gaudiot/sezzle-challenge/calculator"
	"github.com/go-chi/chi/v5"
)

const maxResultIntegerDigits = 33

var errResultTooLarge = errors.New("Result cannot exceed 33 integer digits")

type calculatorResponse struct {
	Result string `json:"result"`
}

type CalculatorHandler struct{}

func NewCalculatorHandler() *CalculatorHandler {
	return &CalculatorHandler{}
}

func (handler *CalculatorHandler) Router(r chi.Router) {
	r.HandleFunc("/", handler.handleCalculate)
}

func (handler *CalculatorHandler) handleCalculate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.Header().Set("Allow", http.MethodGet)
		writeError(w, "Only GET is supported", http.StatusMethodNotAllowed)
		return
	}
	queryParams := r.URL.Query()
	encodedQueries := queryParams["query"]
	if len(encodedQueries) != 1 {
		writeError(w, "Exactly one 'query' parameter is required", http.StatusBadRequest)
		return
	}
	encodedQuery := encodedQueries[0]

	if encodedQuery == "" {
		writeError(w, "Missing required 'query' parameter", http.StatusBadRequest)
		return
	}

	decodedQuery, err := base64.StdEncoding.DecodeString(encodedQuery)
	if err != nil {
		writeError(w, "Invalid Base64 query", http.StatusBadRequest)
		return
	}

	query := string(decodedQuery)

	if len(query) > calculator.MaxExpressionLength {
		writeError(w, "Query parameter cannot be larger than 64 characters", http.StatusBadRequest)
		return
	}

	engine := calculator.New()

	result, err := engine.Evaluate(query)
	if err != nil {
		status := http.StatusBadRequest
		if errors.Is(err, calculator.ErrUndefinedValue) || errors.Is(err, calculator.ErrDivideByZero) || errors.Is(err, calculator.ErrCalculationLimit) {
			status = http.StatusUnprocessableEntity
		}
		writeError(w, err.Error(), status)
		return
	}

	formatted, err := formatResult(result)
	if err != nil {
		writeError(w, err.Error(), http.StatusUnprocessableEntity)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&calculatorResponse{
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

func writeError(w http.ResponseWriter, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]string{"error": message})
}
