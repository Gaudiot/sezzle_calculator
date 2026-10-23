import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "Tally — Calculator",
  description: "A focused calculator with a running history.",
};

export default function RootLayout({ children }: LayoutProps<"/">) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}
