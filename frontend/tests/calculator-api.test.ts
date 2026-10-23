import { beforeEach, describe, expect, it, vi } from "vitest";
import { evaluateExpression } from "@/lib/calculator-api";

const fetchMock = vi.fn();

beforeEach(() => {
  fetchMock.mockReset();
  vi.stubGlobal("fetch", fetchMock);
});

describe("calculator API", () => {
  it("sends a URL-encoded Base64 query through the same-origin proxy", async () => {
    fetchMock.mockResolvedValue(new Response(JSON.stringify({ result: "20" })));
    await expect(evaluateExpression("(2+3)*4")).resolves.toBe("20");
    expect(fetchMock).toHaveBeenCalledWith("/api/math?query=KDIrMykqNA%3D%3D");
  });

  it("preserves large string results", async () => {
    const result = "100000000000000000000000000000000";
    fetchMock.mockResolvedValue(new Response(JSON.stringify({ result })));
    await expect(evaluateExpression(result)).resolves.toBe(result);
  });

  it.each(["", "+1", "1++2", "1.2.3", "1+", "1.", "(1", "2(3)", "1".repeat(65)])(
    "blocks invalid or incomplete %s before making a request", async (expression) => {
      await expect(evaluateExpression(expression)).rejects.toThrow("Enter a complete, valid expression.");
      expect(fetchMock).not.toHaveBeenCalled();
    },
  );

  it.each([" 1+2", "1+2 ", "1 +2", "1\t+2", "1\n+2", "1\u00a0+2", "1\u2003+2"])(
    "rejects whitespace in %s before making a request", async (expression) => {
      await expect(evaluateExpression(expression)).rejects.toThrow("Expressions cannot contain whitespace.");
      expect(fetchMock).not.toHaveBeenCalled();
    },
  );

  it.each([
    ['{"error":"Division by zero"}', "Division by zero"],
    ["Invalid expression", "Invalid expression"],
    ["", "Unable to calculate right now."],
  ])("handles an error response: %s", async (body, message) => {
    fetchMock.mockResolvedValue(new Response(body, { status: 422 }));
    await expect(evaluateExpression("1/0")).rejects.toThrow(message);
  });
});
