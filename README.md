# Sezzle Calculator

A full-stack calculator built for the [challenge](challenge.md), with a React/TypeScript frontend using Next.js and a Go REST API using Chi.

Supports addition, subtraction, multiplication, division, exponentiation (`^`), fractional roots, percentage (`%`), unary signs, and parentheses. The responsive interface accepts keyboard or keypad input, keeps a session history, and lets you continue calculating from a result. Whitespace is rejected in both layers.

## Setup and local development

Prerequisites:

- Git.
- Go 1.27.1 or newer, as declared in `backend/go.mod`.
- Node.js 24.x and npm. The frontend test tooling also supports 22.x from 22.12, or 26 and newer.

Clone the [repository](https://github.com/Gaudiot/sezzle_calculator):

```sh
git clone https://github.com/Gaudiot/sezzle_calculator.git
cd sezzle_calculator
```

From the repository root, start the backend in one terminal:

```sh
cd backend
go mod download
go run .
```

The API listens at `http://localhost:8000`. Its port is currently fixed in `backend/main.go`; restart the process after backend changes.

Start the frontend in another terminal, also beginning at the repository root:

```sh
cd frontend
npm ci
npm run dev
```

Open `http://localhost:3000`. The frontend sends requests to `/api/math`, which Next.js proxies to the backend's `/math` endpoint.

To use a different backend origin, create `frontend/.env.local` before starting or building the frontend:

```dotenv
BACKEND_URL=http://localhost:8000
```

Use an origin without a trailing slash. Restart the frontend after changing this configuration.

For a production frontend build, run these commands from `frontend/` while keeping the backend running separately:

```sh
npm run build
npm start
```

## API usage

`GET /math?query=<base64-expression>` evaluates an expression. Encode it with standard Base64, then URL-encode the parameter. Results are JSON objects containing a decimal **string**, preserving large numbers when consumed by JavaScript.

The examples use `curl`; in Windows PowerShell, use `curl.exe`.

```sh
curl "http://localhost:8000/math?query=Misz"
```

Evaluates `2+3`:

```json
{"result":"5"}
```

Use `--data-urlencode` to handle Base64 padding automatically. This evaluates `(2+3)*4`:

```sh
curl --get --data-urlencode "query=KDIrMykqNA==" "http://localhost:8000/math"
```

```json
{"result":"20"}
```

Division by zero (`1/0`) returns HTTP `422`:

```sh
curl -i "http://localhost:8000/math?query=MS8w"
```

```json
{"error":"Division by zero"}
```

See [endpoints.md](endpoints.md) for the complete syntax, result limits, and response statuses.

## Design decisions and assumptions

- **Exact arithmetic:** the Go service parses decimal tokens into standard-library `math/big.Rat` values. Basic arithmetic, percentages, and integer powers remain exact, avoiding floating-point errors for large integers and decimal cancellation. Square roots use exact rational results when possible, otherwise 256-bit precision. Other fractional powers use floating-point approximation and can lose precision for large results or near rounding boundaries. The final result is rounded, to at most four decimal places, with halves rounded away from zero. Trailing zeros and negative zero are normalized.
- **64-character input limit:** this is a deliberate UX choice for a toy project focused on simple calculations. It leaves enough room for everyday expressions and for continuing from a large result, while keeping entry, editing, and the display manageable. It also bounds parsing and calculation work.
- **33-digit integer output limit:** the original motivation was to define a supported numeric range while pursuing precise results and avoiding floating-point errors. The limit applies after rounding and excludes the sign and fractional digits; it includes decimal literals for both `10^25` and `10^32`. Precision comes from exact rational arithmetic, rather than from the digit cap itself: `float64` cannot guarantee 33-digit precision. The cap keeps the supported output range explicit. Integer powers have bounded computation cost; see the API contract for limits.
- **Explicit parsing:** a recursive-descent parser implements arithmetic precedence directly, without `eval`. The `httpapi` handler handles Base64 decoding, validation responses, and result formatting; parsing and arithmetic stay in the independent `calculator` package.
- **Frontend expression validation:** typing, keypad entry, and paste are checked for valid expression structure. Leading `+`, consecutive operators, repeated decimal points within a number, unmatched closing parentheses, and implicit multiplication are rejected. Negative numbers are allowed at the start or after `(`; write `1*(-2)` rather than `1*-2`. Unfinished prefixes such as `1+`, `1.`, and `(` remain editable, but cannot be submitted until complete. The backend remains responsible for arithmetic errors such as division by zero.
- **Same-origin requests:** the Next.js proxy lets the browser reach the Go service without requiring cross-origin access. The backend URL is configured on the Next.js server.
- **Session-only history:** history lives in React state and resets on refresh. New entries appear at the bottom. An operator continues from the displayed result; a new digit starts a fresh expression. Continuing uses the displayed, rounded result.
- **Scope:** the API supports basic and optional operations. Scientific notation and implicit multiplication are not implemented. Power wraps the current expression as `(expression)^(`; Square root inserts the complete `(expression)^(1/2)`. Undefined values such as `0^0`, `0/0`, and negative square roots return HTTP 422. Percentage is a postfix division by 100. No database is required.

## Verification and coverage

From `backend/`:

```sh
go vet ./...
go test ./...
go test "-coverprofile=coverage.out" ./...
go tool cover "-func=coverage.out"
go tool cover "-html=coverage.out" -o coverage.html
```

Open the generated `coverage.html` for a browsable report. Backend tests cover arithmetic, precedence, invalid expressions, whitespace, exact large-number calculations, final rounding, size boundaries, and HTTP responses. The decimal literals for `10^25` and `10^32` have regression tests asserting unchanged output.

From `frontend/`:

```sh
npm test
npm run test:coverage
npm run lint
npx tsc --noEmit --incremental false
```

The frontend includes automated unit and component tests using Vitest, Testing Library, and jsdom. `npm test` runs the suite once; `npm run test:coverage` generates a report at `frontend/coverage/index.html`. For continuous feedback during development, run `npm run test:watch` separately. Tests cover expression validation, calculator interactions, and API request/error handling without requiring a live backend.

For a manual browser check, evaluate `2+3`, continue with `*4`, start fresh with a digit, try invalid input such as `+1`, `1++2`, or `1.2.3`, and add enough history entries to verify scrolling to the latest result.

## Project layout

| Path | Purpose |
| --- | --- |
| `backend/main.go` | HTTP server entry point |
| `backend/calculator/` | Expression parser, arithmetic, and unit tests |
| `backend/httpapi/` | HTTP handler, result formatting, and API tests |
| `frontend/app/` | Next.js pages and shared styles |
| `frontend/components/` | Calculator, keypad, and history UI |
| `frontend/lib/calculator-api.ts` | API client and error handling |
| `frontend/lib/expression-validation.ts` | Shared frontend expression validation |
| `frontend/tests/` | Automated validation, calculator, and API-client tests |
| [endpoints.md](endpoints.md) | API contract |
| [challenge.md](challenge.md) | Assignment requirements |
| [prompts.md](prompts.md) | Recorded AI development prompts |
