import { Calculator } from "@/components/calculator";

export default function Home() {
  return (
    <main className="page-shell">
      <header className="page-header">
        <div className="brand-mark" aria-hidden="true">
          =
        </div>
        <div>
          <p className="eyebrow">A small arithmetic workspace</p>
          <h1>Figure it out.</h1>
        </div>
      </header>

      <Calculator />

      <footer className="page-footer">
        <span>Four operations. No distractions.</span>
        <span>Keyboard ready</span>
      </footer>
    </main>
  );
}
