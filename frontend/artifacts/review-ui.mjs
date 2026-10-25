import { writeFile } from "node:fs/promises";

const pages = await fetch(`http://localhost:${process.env.BROWSER_DEBUG_PORT ?? "9222"}/json`).then((r) => r.json());
const page = pages.find((p) => p.type === "page" && p.url.includes("localhost:3000"));
const socket = new WebSocket(page.webSocketDebuggerUrl);
await new Promise((resolve) => socket.addEventListener("open", resolve, { once: true }));
let id = 0;
const pending = new Map();
socket.addEventListener("message", ({ data }) => {
  const message = JSON.parse(data);
  if (pending.has(message.id)) {
    const { resolve, reject } = pending.get(message.id);
    pending.delete(message.id);
    if (message.error) reject(new Error(JSON.stringify(message.error)));
    else resolve(message.result);
  }
});
function send(method, params = {}) {
  return new Promise((resolve, reject) => {
    const requestId = ++id;
    const timeout = setTimeout(() => reject(new Error(`Timed out: ${method}`)), 30000);
    pending.set(requestId, {
      resolve: result => { clearTimeout(timeout); resolve(result); },
      reject: error => { clearTimeout(timeout); reject(error); },
    });
    socket.send(JSON.stringify({ id: requestId, method, params }));
  });
}

try {
  for (const [name, width, height] of [["mobile", 390, 1100], ["desktop", 1440, 1000]]) {
    await send("Emulation.setDeviceMetricsOverride", { width, height, deviceScaleFactor: 1, mobile: width < 720 });
    await send("Runtime.evaluate", { expression: "document.fonts.ready.then(() => new Promise(r => requestAnimationFrame(() => requestAnimationFrame(r))))", awaitPromise: true });
    const metrics = await send("Runtime.evaluate", { expression: "JSON.stringify({viewport:innerWidth, content:document.documentElement.scrollWidth, keys:[...document.querySelectorAll('[aria-label=\"Calculator keypad\"] button')].map(k=>({label:k.textContent,x:k.getBoundingClientRect().x,y:k.getBoundingClientRect().y}))})", returnByValue: true });
    const layout = JSON.parse(metrics.result.value);
    console.log(name, JSON.stringify(layout));
    if (layout.content > width || layout.viewport !== width) throw new Error("Viewport overflow");
    let previousRowY = -Infinity;
    for (const row of [["7", "8", "9"], ["4", "5", "6"], ["1", "2", "3"]]) {
      const keys = row.map(label => layout.keys.find(key => key.label === label));
      if (!keys.every(key => key.y === keys[0].y) || !(keys[0].x < keys[1].x && keys[1].x < keys[2].x) || keys[0].y <= previousRowY) throw new Error("Misaligned keypad");
      previousRowY = keys[0].y;
    }
    const screenshot = await send("Page.captureScreenshot", { format: "png", captureBeyondViewport: true });
    await writeFile(new URL(`./${process.env.SCREENSHOT_PREFIX ?? ""}${name}.png`, import.meta.url), Buffer.from(screenshot.data, "base64"));
  }
} finally {
  socket.close();
}
