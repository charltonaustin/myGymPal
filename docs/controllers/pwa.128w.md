---
level: 128w
parent: ../controllers.32w.md
deeper: pwa.256w.md
relates-to:
  - default.128w.md
source: controllers/pwa.go
---

The PWA controller serves the three assets required to make myGymPal installable as a Progressive Web App. GET /sw.js
reads `static/sw.js` from disk and serves it with `Content-Type: application/javascript`, `Cache-Control: no-cache`, and
`Service-Worker-Allowed: /` headers; returns 404 if the file is missing. GET /manifest.json reads `static/manifest.json`
and serves it with `Content-Type: application/manifest+json`; returns 404 if missing. GET /offline renders
`offline.tpl` — the page served by the service worker when the user has no network connection. None of these handlers
perform any session check or repository call. The routes are registered at the root path so the service worker has the
correct scope for the entire application.
