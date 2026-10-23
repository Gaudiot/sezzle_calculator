import { getExpressionState } from "./expression-validation";

const API_URL = "/api/math";

type MathResponse = {
  result: string;
};

type ErrorResponse = {
  error: string;
};

export async function evaluateExpression(expression: string): Promise<string> {
  if (/\s/.test(expression)) {
    throw new Error("Expressions cannot contain whitespace.");
  }

  if (getExpressionState(expression) !== "complete") {
    throw new Error("Enter a complete, valid expression.");
  }

  const query = encodeURIComponent(btoa(expression));

  const response = await fetch(`${API_URL}?query=${query}`);

  if (!response.ok) {
    const body = await response.text();
    let message = body.trim() || "Unable to calculate right now.";

    try {
      const data: ErrorResponse = JSON.parse(body);
      if (typeof data.error === "string") message = data.error;
    } catch {
      // The backend also returns plain-text errors.
    }

    throw new Error(message);
  }

  const data: MathResponse = await response.json();
  return data.result;
}
