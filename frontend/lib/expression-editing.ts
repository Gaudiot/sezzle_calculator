// Editing stays separate from evaluation: arithmetic belongs to the backend.
export function appendToExpression(expression: string, value: string, isShowingResult: boolean): string {
  if (value === "root") {
    return `(${expression})^(1/2)`;
  }
  if (value === "^") {
    return `(${expression})^(`;
  }

  const continuesResult = /^[+\-*/%]$/.test(value);
  const prefix = isShowingResult && !continuesResult ? "" : expression;
  return `${prefix}${value}`;
}
