---
level: 32w
parent: README.16w.md
relates-to:
  - routers.32w.md
source: static/
---

PWA assets served from `/static/`. `manifest.json` declares standalone display, dark theme color (`#212529`), and two
icon sizes (192x192, 512x512). `sw.js` pre-caches the `/offline` fallback and Bootstrap CDN assets on install, then
applies cache-first for CDN requests and network-first (with offline fallback) for same-origin app pages. Only GET
requests are intercepted; POST and other mutations pass through unchanged.
