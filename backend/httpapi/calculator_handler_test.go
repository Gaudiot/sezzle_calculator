package httpapi

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gaudiot/sezzle-challenge/calculator"
	"github.com/go-chi/chi/v5"
)

func TestMathResponseFormatting(t *testing.T) {
	tests := []struct{ expression, want string }{
		{"1.234", "1.234"},
		{"1.200", "1.2"},
		{"1.0", "1"},
		{"1/3", "0.3333"},
		{"2/3", "0.6667"},
		{"-2/3", "-0.6667"},
		{"1.23455", "1.2346"},
		{"0.1+0.2", "0.3"},
		{"9.99999", "10"},
		{"-0.00001", "0"},
		{"0.00001", "0"},
		{"-0", "0"},
		{"100", "100"},
		{"0.0001", "0.0001"},
		{strings.Repeat("0", calculator.MaxExpressionLength), "0"},
		{"1.23445", "1.2345"},
		{"-1.23445", "-1.2345"},
		{"0.00005", "0.0001"},
		{"-0.00005", "-0.0001"},
		{"1.2344499", "1.2344"},
		{"(1/3)*3000000", "1000000"},
		{"9007199254740992+1", "9007199254740993"},
		{"9999999999999.9999-9999999999999", "0.9999"},
		{"1/0.00000000000000000001", "100000000000000000000"},
		{"100000000000000000000*100000000000000000000/100000000", "100000000000000000000000000000000"},
		{strings.Repeat("9", maxResultIntegerDigits) + ".99994", strings.Repeat("9", maxResultIntegerDigits) + ".9999"},
		{"-" + strings.Repeat("9", maxResultIntegerDigits), "-" + strings.Repeat("9", maxResultIntegerDigits)},
	}
	for _, tt := range tests {
		t.Run(tt.expression, func(t *testing.T) {
			response := mathRequest(tt.expression)
			if response.Code != http.StatusOK {
				t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
			}
			if got := response.Header().Get("Content-Type"); got != "application/json" {
				t.Errorf("Content-Type = %q", got)
			}
			var body calculatorResponse
			if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if body.Result != tt.want {
				t.Errorf("result = %q, want %q", body.Result, tt.want)
			}
		})
	}
}

// These are decimal literals for 10**25 and 10**32, not exponentiation syntax.
func TestMathPowersOfTenReturnUnchanged(t *testing.T) {
	for _, exponent := range []int{25, 32} {
		literal := "1" + strings.Repeat("0", exponent)
		t.Run(literal, func(t *testing.T) {
			response := mathRequest(literal)
			if response.Code != http.StatusOK {
				t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
			}
			var body calculatorResponse
			if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if body.Result != literal {
				t.Errorf("result = %q, want unchanged %q", body.Result, literal)
			}
		})
	}
}

func TestMathErrors(t *testing.T) {
	tests := []struct {
		name, expression string
		status           int
	}{
		{"division by zero", "1/0", http.StatusUnprocessableEntity},
		{"division by exact zero", "1/(0.1+0.2-0.3)", http.StatusUnprocessableEntity},
		{"output too large", "1" + strings.Repeat("0", maxResultIntegerDigits), http.StatusUnprocessableEntity},
		{"rounding over output limit", strings.Repeat("9", maxResultIntegerDigits) + ".99995", http.StatusUnprocessableEntity},
		{"negative rounding over output limit", "-" + strings.Repeat("9", maxResultIntegerDigits) + ".99995", http.StatusUnprocessableEntity},
		{"input too long", strings.Repeat("0", calculator.MaxExpressionLength+1), http.StatusBadRequest},
		{"invalid expression", "1+", http.StatusBadRequest},
		{"power syntax", "10**32", http.StatusBadRequest},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := mathRequest(tt.expression)
			if response.Code != tt.status {
				t.Fatalf("status = %d, want %d; body = %s", response.Code, tt.status, response.Body.String())
			}
			if response.Header().Get("Content-Type") != "application/json" {
				t.Error("error response must be JSON")
			}
			var body struct {
				Error string `json:"error"`
			}
			if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil || body.Error == "" {
				t.Fatalf("expected useful JSON error, got %s", response.Body.String())
			}
		})
	}
}

func TestMathRejectsWhitespace(t *testing.T) {
	expressions := []string{" 1+2", "1+2 ", "1 +2", "1\t+2", "1\n+2", "1\r+2", "1\v+2", "1\f+2", "1\u00a0+2", "1\u2003+2", "\t", "1 2"}
	for _, expression := range expressions {
		t.Run(expression, func(t *testing.T) {
			if _, err := calculator.New().Evaluate(expression); err != calculator.ErrInvalidExpression {
				t.Errorf("Evaluate error = %v", err)
			}
			response := mathRequest(expression)
			if response.Code != http.StatusBadRequest {
				t.Errorf("status = %d, want 400; body = %s", response.Code, response.Body.String())
			}
		})
	}
}

func TestMathRejectsInvalidQueryParameters(t *testing.T) {
	for _, query := range []string{"", "query=", "query=!!!", "query=MQ%3D%3D&query=Mg%3D%3D", "query=%2Fw%3D%3D"} {
		t.Run(query, func(t *testing.T) {
			router := chi.NewRouter()
			router.Route("/math", NewCalculatorHandler().Router)
			response := httptest.NewRecorder()
			router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/math?"+query, nil))
			if response.Code != http.StatusBadRequest {
				t.Fatalf("status = %d, body = %s", response.Code, response.Body.String())
			}
			var body map[string]string
			if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil || body["error"] == "" {
				t.Fatalf("expected JSON error, got %s", response.Body.String())
			}
		})
	}
}

func mathRequest(expression string) *httptest.ResponseRecorder {
	router := chi.NewRouter()
	router.Route("/math", NewCalculatorHandler().Router)
	query := url.Values{"query": {base64.StdEncoding.EncodeToString([]byte(expression))}}
	request := httptest.NewRequest(http.MethodGet, "/math?"+query.Encode(), nil)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, request)
	return response
}

func TestAdvancedHTTPResults(t *testing.T) {
	for expression, want := range map[string]string{
		"0^-2": "0", "0^0.5": "0", "0^(-1/2)": "0", "0^3": "0", "2^3^2": "512", "9^(1/2)": "3", "2^(1/2)": "1.4142", "8^(1/3)": "2",
		"(20+30)%": "0.5", "200*10%": "20", "10^32": "100000000000000000000000000000000",
	} {
		t.Run(expression, func(t *testing.T) {
			response := mathRequest(expression)
			var body calculatorResponse
			if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if response.Code != http.StatusOK || body.Result != want {
				t.Fatalf("status %d, body %s; want %s", response.Code, response.Body.String(), want)
			}
		})
	}
	for _, expression := range []string{"2^10001", "9^9^9", "10^33"} {
		if response := mathRequest(expression); response.Code != http.StatusUnprocessableEntity {
			t.Errorf("%s: status %d, body %s", expression, response.Code, response.Body.String())
		}
	}
}

func TestUnsupportedMethod(t *testing.T) {
	router := chi.NewRouter()
	router.Route("/math", NewCalculatorHandler().Router)
	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodPost, "/math", nil))
	var body map[string]string
	if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
		t.Fatal(err)
	}
	if response.Code != http.StatusMethodNotAllowed || response.Header().Get("Allow") != "GET" || body["error"] == "" {
		t.Fatalf("unexpected response: %v %s", response, response.Body.String())
	}
}

func TestUndefinedValues(t *testing.T) {
	for _, expression := range []string{"0^0", "0/0", "(1-1)^(2-2)", "(1-1)/(2-2)"} {
		t.Run(expression, func(t *testing.T) {
			response := mathRequest(expression)
			var body map[string]string
			if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if response.Code != http.StatusUnprocessableEntity || body["error"] != calculator.ErrUndefinedValue.Error() {
				t.Fatalf("status %d, body %s", response.Code, response.Body.String())
			}
		})
	}
}

func TestNegativeRootsAreInvalidExpressions(t *testing.T) {
	for _, expression := range []string{"(-4)^(1/2)", "(-9)^0.5", "(-8)^(1/3)", "(-4)^(-1/2)"} {
		t.Run(expression, func(t *testing.T) {
			response := mathRequest(expression)
			var body map[string]string
			if err := json.Unmarshal(response.Body.Bytes(), &body); err != nil {
				t.Fatal(err)
			}
			if response.Code != http.StatusBadRequest || body["error"] != calculator.ErrInvalidExpression.Error() {
				t.Fatalf("status %d, body %s", response.Code, response.Body.String())
			}
		})
	}
}
