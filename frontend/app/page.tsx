import { Calculator } from "@/components/calculator";
import styles from "./page.module.css";

export default function Home() {
  return (
    <main className={styles["page-shell"]}>
      <header className={styles["page-header"]}>
        <div className={styles["brand-mark"]} aria-hidden="true">
          =
        </div>
        <div>
          <p className={styles.eyebrow}>A small arithmetic workspace</p>
          <h1 className={styles.heading}>Figure it out.</h1>
        </div>
      </header>

      <Calculator />

      <footer className={styles["page-footer"]}>
        <span>Four operations. No distractions.</span>
        <span>Keyboard ready</span>
      </footer>
    </main>
  );
}
