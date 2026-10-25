# Calculator API

```http
GET /math?query=<base64-expression>
```

Local backend: `http://localhost:8000`. The frontend proxies `/api/math` to this endpoint.

## Request

Provide exactly one nonempty `query` parameter: a standard RFC 4648 Base64-encoded expression, URL-encoded when building the request (for example, `=` becomes `%3D`). The decoded expression must be valid UTF-8, contain 1–64 characters, and contain no whitespace.

Supported syntax:

- Decimal numbers such as `12`, `-3`, and `0.5`; a decimal point requires digits on both sides (`.5` and `5.` are invalid).
- Binary `+`, `-`, `*`, `/`, `^`, postfix `%`, unary `+` and `-`, and parentheses. Roots use fractional exponents: `9^(1/2)` returns `3`; `50%` means `50/100`.
- Precedence: parentheses, postfix percentage, exponentiation, unary signs, multiplication/division, addition/subtraction. Powers associate right to left (`2^3^2 = 512`); other binary operators associate left to right. Thus `-2^2 = -4` and `2^-2 = 0.25`. Repeated `%` divides by 100 each time.

All other syntax is invalid, including scientific notation, implicit multiplication, variables, functions, and `**`. Expressions are parsed directly without `eval`.

## Responses

Success: HTTP `200`, `Content-Type: application/json`:

```json
{"result":"1.2345"}
```

The result is a decimal **string**. Basic arithmetic, percentages, and integer powers are exact until final formatting. Rational square roots are exact when possible; other square roots use 256-bit precision. Other fractional powers use floating-point approximation, so large or rounding-boundary results may lose precision. Formatting:

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
| `422` | Division by zero, undefined value (including invalid real powers), computation limit, or a rounded result exceeding 33 integer digits |
| `500` | Unexpected server failure |

## Examples

| Expression | Request | Response |
| --- | --- | --- |
| `2+3` | `/math?query=Misz` | `200` `{"result":"5"}` |
| `(2+3)*4` | `/math?query=KDIrMykqNA%3D%3D` | `200` `{"result":"20"}` |
| `1/3` | `/math?query=MS8z` | `200` `{"result":"0.3333"}` |
| `1/0` | `/math?query=MS8w` | `422` `{"error":"Division by zero"}` |

Both decimal literals and exponentiation such as `10^32` are supported. Integer exponents have magnitude at most 10,000, with a 65,536-bit estimated allocation budget per power. Negative bases require integer exponents; `0^0` and `0/0` return `422` with `{"error":"Undefined value"}`, and zero to a negative power is division by zero. Non-finite or underflowed fractional powers return `422`.
