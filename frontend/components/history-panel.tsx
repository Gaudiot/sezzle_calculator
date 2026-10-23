"use client";

import { useEffect, useRef } from "react";

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
    <aside className="history-panel" aria-label="Calculation history">
      <div className="history-header">
        <div>
          <p className="panel-label">Session</p>
          <h2>History</h2>
        </div>
        {entries.length > 0 && (
          <button className="text-button" type="button" onClick={onClear}>
            Clear
          </button>
        )}
      </div>

      <div className="history-list" ref={listRef} aria-live="polite">
        {entries.length === 0 ? (
          <div className="empty-history">
            <span aria-hidden="true">↳</span>
            <p>Your calculations will collect here.</p>
          </div>
        ) : (
          <div className="history-items">
            {entries.map((entry, index) => (
              <button
                className="history-entry"
                key={entry.id}
                type="button"
                onClick={() => onReuse(entry.expression)}
                aria-label={`Reuse ${entry.expression}`}
              >
                <span className="history-index">{String(index + 1).padStart(2, "0")}</span>
                <span className="history-calculation">
                  <span className="history-expression">{entry.expression}</span>
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
