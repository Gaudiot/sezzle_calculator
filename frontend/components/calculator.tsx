"use client";

import { useRef, useState } from "react";
import type { ChangeEvent, ClipboardEvent, FormEvent, KeyboardEvent } from "react";
import { evaluateExpression } from "@/lib/calculator-api";
import { appendToExpression } from "@/lib/expression-editing";
import { getExpressionState, MAX_EXPRESSION_LENGTH } from "@/lib/expression-validation";
import { CalculatorKeypad } from "./calculator-keypad";
import { HistoryPanel } from "./history-panel";
import type { HistoryEntry } from "./history-panel";
import styles from "./calculator.module.css";

const INVALID_CHARACTERS = /[^0-9+\-*/^%().]/;

export function Calculator() {
  const [expression, setExpression] = useState("");
  const [isShowingResult, setIsShowingResult] = useState(false);
  const [error, setError] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [history, setHistory] = useState<HistoryEntry[]>([]);
  const inputRef = useRef<HTMLInputElement>(null);

  const updateExpression = (value: string) => {
    if (getExpressionState(value) === "invalid") {
      setError("That input would create an invalid expression.");
      return false;
    }

    setExpression(value);
    setIsShowingResult(false);
    setError("");
    return true;
  };

  const append = (value: string) => {
    if ((value === "^" || value === "root") && getExpressionState(expression) !== "complete") {
      setError("Complete the expression before adding an exponent.");
      inputRef.current?.focus();
      return;
    }
    updateExpression(appendToExpression(expression, value, isShowingResult));
    inputRef.current?.focus();
  };

  const clear = () => {
    setExpression("");
    setIsShowingResult(false);
    setError("");
    inputRef.current?.focus();
  };

  const evaluate = async () => {
    if (isLoading) return;

    if (!expression) {
      setError("Enter an expression first.");
      return;
    }

    if (getExpressionState(expression) !== "complete") {
      setError("Complete the expression before calculating.");
      return;
    }

    setIsLoading(true);
    setError("");

    try {
      const nextResult = await evaluateExpression(expression);
      setExpression(nextResult);
      setIsShowingResult(true);
      setHistory((entries) => [
        ...entries,
        { id: Date.now(), expression, result: nextResult },
      ]);
    } catch (caught) {
      setError(caught instanceof Error ? caught.message : "Unable to calculate right now.");
    } finally {
      setIsLoading(false);
    }
  };

  const handleSubmit = (event: FormEvent<HTMLFormElement>) => {
    event.preventDefault();
    void evaluate();
  };

  const handleChange = (event: ChangeEvent<HTMLInputElement>) => {
    updateExpression(event.target.value);
  };

  const handlePaste = (event: ClipboardEvent<HTMLInputElement>) => {
    event.preventDefault();
    const pastedText = event.clipboardData.getData("text");
    const start = event.currentTarget.selectionStart ?? expression.length;
    const end = event.currentTarget.selectionEnd ?? expression.length;
    const nextExpression = isShowingResult
      ? pastedText
      : `${expression.slice(0, start)}${pastedText}${expression.slice(end)}`;
    if (updateExpression(nextExpression)) {
      const caret = isShowingResult ? pastedText.length : start + pastedText.length;
      queueMicrotask(() => inputRef.current?.setSelectionRange(caret, caret));
    }
  };

  const handleKeyDown = (event: KeyboardEvent<HTMLInputElement>) => {
    if (event.key === "^" && !event.ctrlKey && !event.metaKey && !event.altKey) {
      event.preventDefault();
      append("^");
      return;
    }
    if (event.key === "Enter") {
      event.preventDefault();
      void evaluate();
      return;
    }

    if (
      isShowingResult &&
      event.key.length === 1 &&
      !INVALID_CHARACTERS.test(event.key) &&
      !event.ctrlKey &&
      !event.metaKey &&
      !event.altKey
    ) {
      event.preventDefault();
      append(event.key);
      return;
    }

    if (event.key.length === 1 && !event.ctrlKey && !event.metaKey && !event.altKey) {
      const start = event.currentTarget.selectionStart ?? expression.length;
      const end = event.currentTarget.selectionEnd ?? expression.length;
      const nextExpression = `${expression.slice(0, start)}${event.key}${expression.slice(end)}`;
      if (getExpressionState(nextExpression) === "invalid") {
        event.preventDefault();
        setError("That input would create an invalid expression.");
      }
    }
  };

  return (
    <section className={styles["calculator-workspace"]} aria-label="Calculator">
      <form className={styles["calculator-panel"]} onSubmit={handleSubmit}>
        <div className={styles.display}>
          <label htmlFor="expression">{isShowingResult ? "Result" : "Expression"}</label>
          <input
            ref={inputRef}
            id="expression"
            name="expression"
            type="text"
            inputMode="decimal"
            autoComplete="off"
            spellCheck={false}
            disabled={isLoading}
            maxLength={MAX_EXPRESSION_LENGTH}
            value={expression}
            placeholder="0"
            onChange={handleChange}
            onPaste={handlePaste}
            onKeyDown={handleKeyDown}
            aria-describedby={error ? "calculator-error" : "expression-hint"}
            aria-invalid={Boolean(error)}
          />
          <span className="sr-only" role="status">
            {isShowingResult ? `Result: ${expression}` : ""}
          </span>
          <div className={styles["result-row"]}>
            <span id="expression-hint">
              {isShowingResult ? "Use an operator to continue calculating." : `${expression.length}/${MAX_EXPRESSION_LENGTH}`}
            </span>
            {isLoading && <span className={styles["loading-dots"]} role="status" aria-label="Calculating">•••</span>}
          </div>
          <div className={styles["error-slot"]}>
            {error && (
              <p id="calculator-error" role="alert">
                {error}
              </p>
            )}
          </div>
        </div>

        <CalculatorKeypad
          disabled={isLoading}
          onAppend={append}
          onBackspace={() => updateExpression(expression.slice(0, -1))}
          onClear={clear}
          onEvaluate={() => void evaluate()}
        />
        <p className={styles["operation-hint"]}>Power and square root apply to the whole expression. Square root adds ^(1/2). % divides by 100.</p>
      </form>

      <HistoryPanel
        entries={history}
        onClear={() => setHistory([])}
        onReuse={(value) => {
          updateExpression(value);
          inputRef.current?.focus();
        }}
      />
    </section>
  );
}
