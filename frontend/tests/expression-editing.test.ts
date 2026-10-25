import { describe, expect, it } from "vitest";
import { appendToExpression } from "@/lib/expression-editing";

describe("expression editing", () => {
  it.each(["^"])("groups the entire base for %s", (operation) => {
    expect(appendToExpression("2+7", operation, false)).toBe("(2+7)^(");
    expect(appendToExpression("9", operation, true)).toBe("(9)^(");
  });

  it("adds a complete square root to the expression or displayed result", () => {
    expect(appendToExpression("2+7", "root", false)).toBe("(2+7)^(1/2)");
    expect(appendToExpression("9", "root", true)).toBe("(9)^(1/2)");
  });

  it("continues results with operations and replaces results with new input", () => {
    expect(appendToExpression("50", "%", true)).toBe("50%");
    expect(appendToExpression("5", "+", true)).toBe("5+");
    expect(appendToExpression("5", "2", true)).toBe("2");
    expect(appendToExpression("5", "2", false)).toBe("52");
  });
});
