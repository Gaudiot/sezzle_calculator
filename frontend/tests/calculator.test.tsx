import { fireEvent, render, screen, waitFor, within } from "@testing-library/react";
import userEvent from "@testing-library/user-event";
import { beforeEach, describe, expect, it, vi } from "vitest";
import { Calculator } from "@/components/calculator";
import { evaluateExpression } from "@/lib/calculator-api";

vi.mock("@/lib/calculator-api", () => ({ evaluateExpression: vi.fn() }));
const evaluateMock = vi.mocked(evaluateExpression);

beforeEach(() => {
  evaluateMock.mockReset();
  evaluateMock.mockResolvedValue("5");
});

function expressionInput() {
  return screen.getByRole("textbox") as HTMLInputElement;
}

describe("calculator input", () => {
  it.each(["Exponentiation"])("wraps the whole expression using %s and submits the entered exponent", async (name) => {
    const user = userEvent.setup();
    render(<Calculator />);
    await user.type(expressionInput(), "2+7");
    await user.click(screen.getByRole("button", { name }));
    expect(expressionInput()).toHaveValue("(2+7)^(");
    await user.click(screen.getByRole("button", { name: "Calculate result" }));
    expect(evaluateMock).not.toHaveBeenCalled();
    await user.type(expressionInput(), "1/2){Enter}");
    await waitFor(() => expect(evaluateMock).toHaveBeenCalledWith("(2+7)^(1/2)"));
  });

  it("inserts a complete square root and submits without further input", async () => {
    const user = userEvent.setup();
    render(<Calculator />);
    await user.type(expressionInput(), "2+7");
    await user.click(screen.getByRole("button", { name: "Square root" }));
    expect(expressionInput()).toHaveValue("(2+7)^(1/2)");
    await user.click(screen.getByRole("button", { name: "Calculate result" }));
    await waitFor(() => expect(evaluateMock).toHaveBeenCalledWith("(2+7)^(1/2)"));
    await waitFor(() => expect(expressionInput()).toHaveAccessibleName("Result"));
    await user.click(screen.getByRole("button", { name: "Square root" }));
    expect(expressionInput()).toHaveValue("(5)^(1/2)");
  });

  it("keeps the expression when a square root would exceed the length limit", async () => {
    const user = userEvent.setup();
    render(<Calculator />);
    fireEvent.change(expressionInput(), { target: { value: "1".repeat(57) } });
    await user.click(screen.getByRole("button", { name: "Square root" }));
    expect(expressionInput()).toHaveValue("1".repeat(57));
  });

  it("wraps keyboard exponentiation and accepts postfix percentage", async () => {
    const user = userEvent.setup();
    render(<Calculator />);
    await user.type(expressionInput(), "2+3^2)%{Enter}");
    await waitFor(() => expect(evaluateMock).toHaveBeenCalledWith("(2+3)^(2)%"));
  });

  it("rejects exponentiation for incomplete expressions and when wrapping exceeds the limit", async () => {
    const user = userEvent.setup();
    render(<Calculator />);
    await user.click(screen.getByRole("button", { name: "Square root" }));
    expect(expressionInput()).toHaveValue("");
    expect(screen.getByRole("alert")).toHaveTextContent("Complete the expression");
    fireEvent.change(expressionInput(), { target: { value: "1".repeat(61) } });
    await user.click(screen.getByRole("button", { name: "Exponentiation" }));
    expect(expressionInput()).toHaveValue("1".repeat(61));
  });

  it("applies percentage and power to a displayed result", async () => {
    const user = userEvent.setup();
    render(<Calculator />);
    await user.type(expressionInput(), "2+3{Enter}");
    await waitFor(() => expect(expressionInput()).toHaveAccessibleName("Result"));
    await user.click(screen.getByRole("button", { name: "Percentage" }));
    expect(expressionInput()).toHaveValue("5%");
    await user.click(screen.getByRole("button", { name: "Exponentiation" }));
    expect(expressionInput()).toHaveValue("(5%)^(");
  });

  it("blocks leading plus and invalid keypad operator sequences", async () => {
    const user = userEvent.setup();
    render(<Calculator />);
    await user.click(screen.getByRole("button", { name: "Add" }));
    expect(expressionInput()).toHaveValue("");
    await user.click(screen.getByRole("button", { name: "1" }));
    await user.click(screen.getByRole("button", { name: "Add" }));
    await user.click(screen.getByRole("button", { name: "Multiply" }));
    expect(expressionInput()).toHaveValue("1+");
    await user.click(screen.getByRole("button", { name: "2" }));
    expect(expressionInput()).toHaveValue("1+2");
  });

  it("blocks invalid keyboard characters, operators, and repeated decimal points", async () => {
    const user = userEvent.setup();
    render(<Calculator />);
    await user.type(expressionInput(), "+1.2.*3 x");
    expect(expressionInput()).toHaveValue("1.2*3");
    await user.type(expressionInput(), "+*2");
    expect(expressionInput()).toHaveValue("1.2*3+2");
  });

  it("validates insertion and selection replacement at the caret", async () => {
    const user = userEvent.setup();
    render(<Calculator />);
    await user.type(expressionInput(), "12");
    expressionInput().setSelectionRange(1, 1);
    await user.keyboard("+");
    expect(expressionInput()).toHaveValue("1+2");
    expressionInput().setSelectionRange(1, 2);
    await user.keyboard("*");
    expect(expressionInput()).toHaveValue("1*2");
    expressionInput().setSelectionRange(2, 2);
    await user.keyboard("+");
    expect(expressionInput()).toHaveValue("1*2");
  });

  it("rejects invalid direct changes without overwriting valid input", () => {
    render(<Calculator />);
    fireEvent.change(expressionInput(), { target: { value: "1.2" } });
    fireEvent.change(expressionInput(), { target: { value: "1.2.3" } });
    expect(expressionInput()).toHaveValue("1.2");
    fireEvent.change(expressionInput(), { target: { value: "1".repeat(65) } });
    expect(expressionInput()).toHaveValue("1.2");
    fireEvent.change(expressionInput(), { target: { value: "1".repeat(64) } });
    expect(expressionInput()).toHaveValue("1".repeat(64));
  });

  it.each(["+1", "1++2", "1.2.3", "(1+2))", "2(3)", "1 +2", "1\n+2", "1\u00a0+2"])(
    "rejects pasted %s atomically", (text) => {
      render(<Calculator />);
      fireEvent.change(expressionInput(), { target: { value: "7" } });
      expressionInput().setSelectionRange(0, 1);
      fireEvent.paste(expressionInput(), { clipboardData: { getData: () => text } });
      expect(expressionInput()).toHaveValue("7");
    },
  );

  it("validates pasted text against the full expression and selected range", async () => {
    render(<Calculator />);
    fireEvent.change(expressionInput(), { target: { value: "12" } });
    expressionInput().setSelectionRange(1, 1);
    fireEvent.paste(expressionInput(), { clipboardData: { getData: () => "+" } });
    expect(expressionInput()).toHaveValue("1+2");
    await waitFor(() => expect(expressionInput().selectionStart).toBe(2));
    expressionInput().setSelectionRange(1, 2);
    fireEvent.paste(expressionInput(), { clipboardData: { getData: () => "*" } });
    expect(expressionInput()).toHaveValue("1*2");
    await waitFor(() => expect(expressionInput().selectionStart).toBe(2));
    expressionInput().setSelectionRange(2, 2);
    fireEvent.paste(expressionInput(), { clipboardData: { getData: () => "+" } });
    expect(expressionInput()).toHaveValue("1*2");
  });

  it.each(["-", "1+", "1.", "(", "(1+2"])("allows prefix %s but blocks evaluation", async (prefix) => {
    const user = userEvent.setup();
    render(<Calculator />);
    await user.type(expressionInput(), prefix);
    expect(expressionInput()).toHaveValue(prefix);
    await user.click(screen.getByRole("button", { name: "Calculate result" }));
    expect(evaluateMock).not.toHaveBeenCalled();
    expect(screen.getByRole("alert")).toHaveTextContent("Complete the expression");
  });

  it("accepts negative operands grouped after an operator and balanced parentheses", async () => {
    const user = userEvent.setup();
    render(<Calculator />);
    await user.type(expressionInput(), "1*(-2)");
    await user.keyboard("{Enter}");
    await waitFor(() => expect(evaluateMock).toHaveBeenCalledWith("1*(-2)"));
  });

  it("leaves domain validation to the backend and displays its error", async () => {
    evaluateMock.mockRejectedValue(new Error("Division by zero"));
    const user = userEvent.setup();
    render(<Calculator />);
    await user.type(expressionInput(), "1/0");
    await user.keyboard("{Enter}");
    expect(await screen.findByRole("alert")).toHaveTextContent("Division by zero");
    expect(evaluateMock).toHaveBeenCalledWith("1/0");
    expect(expressionInput()).toHaveValue("1/0");
  });
});

describe("results and history", () => {
  it("preserves result mode when an invalid paste or key is rejected", async () => {
    const user = userEvent.setup();
    render(<Calculator />);
    await user.type(expressionInput(), "2+3{Enter}");
    await waitFor(() => expect(expressionInput()).toHaveAccessibleName("Result"));
    fireEvent.paste(expressionInput(), { clipboardData: { getData: () => "1.2.3" } });
    expect(expressionInput()).toHaveValue("5");
    expect(expressionInput()).toHaveAccessibleName("Result");
    await user.keyboard(")");
    expect(expressionInput()).toHaveValue("5");
    expect(expressionInput()).toHaveAccessibleName("Result");
    await user.keyboard("+");
    expect(expressionInput()).toHaveValue("5+");
  });

  it("shows a result alone, continues through the keypad, and starts fresh with a digit", async () => {
    evaluateMock.mockResolvedValueOnce("5").mockResolvedValueOnce("20");
    const user = userEvent.setup();
    render(<Calculator />);
    await user.type(expressionInput(), "2+3{Enter}");
    await waitFor(() => expect(expressionInput()).toHaveAccessibleName("Result"));
    expect(expressionInput()).toHaveValue("5");
    await user.click(screen.getByRole("button", { name: "Multiply" }));
    await user.click(screen.getByRole("button", { name: "4" }));
    expect(expressionInput()).toHaveValue("5*4");
    await user.click(screen.getByRole("button", { name: "Calculate result" }));
    await waitFor(() => expect(expressionInput()).toHaveValue("20"));
    await user.click(screen.getByRole("button", { name: "7" }));
    expect(expressionInput()).toHaveValue("7");
    expect(expressionInput()).toHaveAccessibleName("Expression");
  });

  it("continues a result with the keyboard, replaces it on paste, and reuses history", async () => {
    const user = userEvent.setup();
    render(<Calculator />);
    await user.type(expressionInput(), "2+3{Enter}");
    await waitFor(() => expect(expressionInput()).toHaveValue("5"));
    await user.click(expressionInput());
    await user.keyboard("*4");
    expect(expressionInput()).toHaveValue("5*4");
    await user.keyboard("{Enter}");
    await waitFor(() => expect(expressionInput()).toHaveAccessibleName("Result"));
    fireEvent.paste(expressionInput(), { clipboardData: { getData: () => "12.5" } });
    expect(expressionInput()).toHaveValue("12.5");
    await user.click(screen.getByRole("button", { name: "Reuse 2+3" }));
    expect(expressionInput()).toHaveValue("2+3");
    const history = screen.getByRole("complementary", { name: "Calculation history" });
    expect(within(history).getAllByRole("button", { name: /^Reuse/ }).map((entry) => entry.getAttribute("aria-label")))
      .toEqual(["Reuse 2+3", "Reuse 5*4"]);
    await user.click(within(history).getByRole("button", { name: "Clear" }));
    expect(within(history).getByText("Your calculations will collect here.")).toBeInTheDocument();
  });

  it("can clear and backspace after validation errors", async () => {
    const user = userEvent.setup();
    render(<Calculator />);
    await user.type(expressionInput(), "12+");
    await user.click(screen.getByRole("button", { name: "Delete last character" }));
    expect(expressionInput()).toHaveValue("12");
    await user.click(screen.getByRole("button", { name: "Clear expression" }));
    expect(expressionInput()).toHaveValue("");
    await user.click(screen.getByRole("button", { name: "Calculate result" }));
    expect(screen.getByRole("alert")).toHaveTextContent("Enter an expression first.");
  });
});
