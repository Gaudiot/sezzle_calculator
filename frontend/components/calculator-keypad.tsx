type CalculatorKeypadProps = {
  disabled: boolean;
  onAppend: (value: string) => void;
  onBackspace: () => void;
  onClear: () => void;
  onEvaluate: () => void;
};

type CalculatorKey = {
  label: string;
  value: string;
  kind?: "utility" | "operator" | "wide" | "equals";
  ariaLabel?: string;
};

const keys: CalculatorKey[] = [
  { label: "Clear", value: "clear", kind: "utility", ariaLabel: "Clear expression" },
  { label: "(", value: "(" },
  { label: ")", value: ")" },
  { label: "÷", value: "/", kind: "operator", ariaLabel: "Divide" },
  { label: "7", value: "7" },
  { label: "8", value: "8" },
  { label: "9", value: "9" },
  { label: "×", value: "*", kind: "operator", ariaLabel: "Multiply" },
  { label: "4", value: "4" },
  { label: "5", value: "5" },
  { label: "6", value: "6" },
  { label: "−", value: "-", kind: "operator", ariaLabel: "Subtract" },
  { label: "1", value: "1" },
  { label: "2", value: "2" },
  { label: "3", value: "3" },
  { label: "+", value: "+", kind: "operator", ariaLabel: "Add" },
  { label: "0", value: "0", kind: "wide" },
  { label: ".", value: "." },
  { label: "⌫", value: "backspace", ariaLabel: "Delete last character" },
  { label: "=", value: "evaluate", kind: "equals", ariaLabel: "Calculate result" },
];

export function CalculatorKeypad({
  disabled,
  onAppend,
  onBackspace,
  onClear,
  onEvaluate,
}: CalculatorKeypadProps) {
  const handleKey = (value: string) => {
    if (value === "clear") onClear();
    else if (value === "backspace") onBackspace();
    else if (value === "evaluate") onEvaluate();
    else onAppend(value);
  };

  return (
    <div className="keypad" aria-label="Calculator keypad">
      {keys.map((key) => (
        <button
          className={`key key--${key.kind ?? "number"}`}
          disabled={disabled}
          key={key.value}
          type="button"
          aria-label={key.ariaLabel}
          onClick={() => handleKey(key.value)}
        >
          {key.label}
        </button>
      ))}
    </div>
  );
}
