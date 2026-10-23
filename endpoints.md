# Calculator API

```http
GET /math?query=<base64-expression>
```

Local backend: `http://localhost:8000`. The frontend proxies `/api/math` to this endpoint.

## Request

Provide exactly one nonempty `query` parameter: a standard RFC 4648 Base64-encoded expression, URL-encoded when building the request (for example, `=` becomes `%3D`). The decoded expression must be valid UTF-8, contain 1–64 characters, and contain no whitespace.

Supported syntax:

- Decimal numbers such as `12`, `-3`, and `0.5`; a decimal point requires digits on both sides (`.5` and `5.` are invalid).
- Binary `+`, `-`, `*`, `/`, unary `+` and `-`, and parentheses.
- Precedence: parentheses, unary signs, multiplication/division, addition/subtraction. Binary operators of equal precedence associate left to right.

All other syntax is invalid, including scientific notation, implicit multiplication, variables, functions, `%`, and exponentiation (`^` or `**`). Expressions are parsed directly without `eval`.

## Responses

Success: HTTP `200`, `Content-Type: application/json`:

```json
{"result":"1.2345"}
```

The result is a decimal **string**. Arithmetic is exact until final formatting:

- Round to at most **4 decimal places**, with halves rounded away from zero.
- Allow at most **33 integer digits after rounding**, excluding the minus sign.
- Remove unnecessary leading zeros from the integer part, trailing zeros from the fraction, and the decimal point when no fraction remains: `1.200` → `1.2`, `1.0` → `1`.
- Keep a zero before fractions (`0.5`, `-0.5`); normalize negative zero to `0`. Never return scientific notation, `NaN`, or infinity.

Calculation and validation errors use a JSON object:

```json
{"error":"Division by zero"}
```

| Status | Meaning |
| --- | --- |
| `400` | Missing/repeated/empty `query`, invalid Base64 or UTF-8, input over 64 characters, or invalid expression |
| `405` | Unsupported method; `Allow: GET` |
| `422` | Division by zero or a rounded result exceeding 33 integer digits |
| `500` | Unexpected server failure |

## Examples

| Expression | Request | Response |
| --- | --- | --- |
| `2+3` | `/math?query=Misz` | `200` `{"result":"5"}` |
| `(2+3)*4` | `/math?query=KDIrMykqNA%3D%3D` | `200` `{"result":"20"}` |
| `1/3` | `/math?query=MS8z` | `200` `{"result":"0.3333"}` |
| `1/0` | `/math?query=MS8w` | `422` `{"error":"Division by zero"}` |

The decimal literals for `10^25` and `10^32` return unchanged; the notation `10^32` itself is unsupported.
