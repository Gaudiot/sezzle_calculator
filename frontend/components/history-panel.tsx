"use client";

import { useEffect, useRef } from "react";
import styles from "./history-panel.module.css";

export type HistoryEntry = {
  id: number;
  expression: string;
  result: string;
};

type HistoryPanelProps = {
  entries: HistoryEntry[];
  onClear: () => void;
  onReuse: (expression: string) => void;
};

export function HistoryPanel({ entries, onClear, onReuse }: HistoryPanelProps) {
  const listRef = useRef<HTMLDivElement>(null);

  useEffect(() => {
    if (listRef.current) {
      listRef.current.scrollTop = listRef.current.scrollHeight;
    }
  }, [entries]);

  return (
    <aside className={styles["history-panel"]} aria-label="Calculation history">
      <div className={styles["history-header"]}>
        <div>
          <p className={styles["panel-label"]}>Session</p>
          <h2>History</h2>
        </div>
        {entries.length > 0 && (
          <button className={styles["text-button"]} type="button" onClick={onClear}>
            Clear
          </button>
        )}
      </div>

      <div className={styles["history-list"]} ref={listRef} aria-live="polite">
        {entries.length === 0 ? (
          <div className={styles["empty-history"]}>
            <span aria-hidden="true">↳</span>
            <p>Your calculations will collect here.</p>
          </div>
        ) : (
          <div className={styles["history-items"]}>
            {entries.map((entry, index) => (
              <button
                className={styles["history-entry"]}
                key={entry.id}
                type="button"
                onClick={() => onReuse(entry.expression)}
                aria-label={`Reuse ${entry.expression}`}
              >
                <span className={styles["history-index"]}>{String(index + 1).padStart(2, "0")}</span>
                <span className={styles["history-calculation"]}>
                  <span className={styles["history-expression"]}>{entry.expression}</span>
                  <strong>{entry.result}</strong>
                </span>
              </button>
            ))}
          </div>
        )}
      </div>
    </aside>
  );
}
