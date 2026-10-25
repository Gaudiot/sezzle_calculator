# Calculator frontend

React and TypeScript with Next.js. Install dependencies with `npm install`, then
run `npm run dev` and open http://localhost:3000. Run the Go backend at
http://localhost:8000. The frontend proxies `/api/math` to the backend; set
`BACKEND_URL` to another backend origin before starting or building when needed.

Use `npm run build` and `npm start` for a production build.

## Operations

Enter numbers, parentheses, `+`, `-`, `*`, `/`, `^`, and postfix `%`.
The Power button and keyboard `^` wrap the current expression and add `^(`,
leaving the exponent for you to enter and close with `)`. Square root inserts
the complete `(expression)^(1/2)`; for example, `2+7` becomes `(2+7)^(1/2)`.
Percentage divides its preceding operand by 100: `50%` is `0.5`, and `(20+30)%` is also `0.5`.

Inputs are limited to 64 characters. Whitespace, implicit multiplication, and
invalid operator sequences are blocked. Negative operands after an operator
must be grouped, such as `2^(-3)`. Calculation errors appear below the input.

## Design and verification

`components/calculator.tsx` owns UI state and coordinates input and requests.
The keypad and history are presentation components. `lib/expression-editing.ts`
contains pure editing rules, `lib/expression-validation.ts` checks syntax and
editable prefixes, and `lib/calculator-api.ts` handles HTTP. The backend owns
arithmetic, precedence, precision, and domain validation; there is no client-side
evaluation or duplicate calculation engine.

Each component keeps its scoped CSS Module alongside its TypeScript file. Page
layout styles live in `app/page.module.css`; `app/globals.css` contains only the
shared theme, resets, and reduced-motion preferences. The existing coordinator,
keypad, and history boundaries remain small enough to maintain without extra
component layers.

Run `npm test`, `npm run test:coverage`, `npm run lint`, and `npx tsc --noEmit`.
Coverage reports are written to `coverage/`. Tests cover expression editing and
validation, API requests and errors, and user interaction through results and
history.
