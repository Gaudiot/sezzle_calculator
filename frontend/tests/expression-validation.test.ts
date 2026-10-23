import { describe, expect, it } from "vitest";
import { getExpressionState, MAX_EXPRESSION_LENGTH } from "@/lib/expression-validation";

describe("complete expressions", () => {
  it.each([
    "0", "12", "-3", "0.5", "0012.50", "2+3*4", "(2+3)*4", "-(2+3)",
    "1*(-2)", "1+(-2)", "-(-2)", "((1))", "1/0", "1.2+3.4", "(1)+(2)",
    "100000000000000000000000000000000",
  ])("accepts %s", (expression) => {
    expect(getExpressionState(expression)).toBe("complete");
  });
});

describe("editable prefixes", () => {
  it.each(["", "-", "(", "(-", "1+", "1*", "1/", "1-", "1.", "(1.", "(1+2", "1*(-"])(
    "allows unfinished %s without permitting evaluation", (expression) => {
      expect(getExpressionState(expression)).toBe("incomplete");
    },
  );
});

describe("invalid expressions", () => {
  it.each([
    "+", "+1", "*1", "/1", "1++2", "1+-2", "1*-2", "1/-2", "--1", "1--2",
    "1**2", "1//2", "1.2.3", ".5", "1.)", "1.+2", "()", "(+)" , "(-)",
    "1)", "(1))", "2(3)", "(2)3", "(2)(3)", "(2).3", "1 2", " 1+2", "1+2 ",
    "1\t+2", "1\n+2", "1\u00a0+2", "1\u2003+2", "1e2", "2^3", "2%", "x+1",
  ])("rejects %s", (expression) => {
    expect(getExpressionState(expression)).toBe("invalid");
  });

  it("accepts exactly 64 characters and rejects 65", () => {
    expect(getExpressionState("1".repeat(MAX_EXPRESSION_LENGTH))).toBe("complete");
    expect(getExpressionState("1".repeat(MAX_EXPRESSION_LENGTH + 1))).toBe("invalid");
  });
});
