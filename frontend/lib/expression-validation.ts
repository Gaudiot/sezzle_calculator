export const MAX_EXPRESSION_LENGTH = 64;

export type ExpressionState = "invalid" | "incomplete" | "complete";

// An unfinished but valid prefix must remain editable, such as "1+" or "(2".
export function getExpressionState(expression: string): ExpressionState {
  if (expression.length > MAX_EXPRESSION_LENGTH || /[^0-9+\-*/^%().]/.test(expression)) {
    return "invalid";
  }

  let position = 0;
  let openParentheses = 0;
  let expectsOperand = true;
  let allowsUnaryMinus = true;

  while (position < expression.length) {
    const character = expression[position];

    if (expectsOperand) {
      if (character === "-" && allowsUnaryMinus) {
        allowsUnaryMinus = false;
        position++;
        continue;
      }

      if (character === "(") {
        openParentheses++;
        allowsUnaryMinus = true;
        position++;
        continue;
      }

      if (!/[0-9]/.test(character)) return "invalid";

      while (position < expression.length && /[0-9]/.test(expression[position])) position++;

      if (expression[position] === ".") {
        position++;
        const fractionStart = position;
        while (position < expression.length && /[0-9]/.test(expression[position])) position++;
        if (position === fractionStart) {
          return position === expression.length ? "incomplete" : "invalid";
        }
      }

      expectsOperand = false;
      allowsUnaryMinus = false;
      continue;
    }

    if (character === ")" && openParentheses > 0) {
      openParentheses--;
      position++;
      continue;
    }

    if (character === "%") {
      position++;
      continue;
    }

    if (/[+\-*/^]/.test(character)) {
      expectsOperand = true;
      // Negative operands after an operator must be grouped: 1*(-2), not 1*-2.
      allowsUnaryMinus = false;
      position++;
      continue;
    }

    return "invalid";
  }

  return expectsOperand || openParentheses > 0 ? "incomplete" : "complete";
}
