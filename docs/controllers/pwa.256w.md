---
level: 256w
parent: ../controllers.32w.md
relates-to:
  - default.256w.md
source: controllers/pwa.go
---

## Routes

| Method | Path           | Handler       |
|--------|----------------|---------------|
| GET    | /sw.js         | ServiceWorker |
| GET    | /manifest.json | Manifest      |
| GET    | /offline       | Offline       |

## Auth requirement

None. No session check on any handler.

## Session keys

None read or written.

## Template variables

None set for ServiceWorker or Manifest (raw byte responses). Offline uses:

- No explicit `c.Data` keys — `offline.tpl` renders with whatever layout defaults Beego provides.

## Templates

- `offline.tpl` — rendered by Offline handler

## Repository calls

None.

## Flash messages

None.

## File serving details

**ServiceWorker (`/sw.js`):**

- Reads `static/sw.js` via `os.ReadFile`
- Sets `Content-Type: application/javascript; charset=utf-8`
- Sets `Cache-Control: no-cache` (ensures browser re-fetches on every load to pick up service worker updates)
- Sets `Service-Worker-Allowed: /` (grants the service worker scope over the entire origin)
- Returns HTTP 404 via `c.Ctx.Output.SetStatus(404)` if file is missing

**Manifest (`/manifest.json`):**

- Reads `static/manifest.json` via `os.ReadFile`
- Sets `Content-Type: application/manifest+json`
- Returns HTTP 404 if file is missing

## Route registration rationale

Routes are registered at `/sw.js` and `/manifest.json` (not under a `/static/` prefix) because browsers enforce that a
service worker's scope is limited to its URL path and below. Serving from the root gives the service worker scope over
all application routes.

## Relationship to other controllers

The offline page is a static fallback with no interaction with the rest of the application. The service worker and
manifest are consumed by the browser, not by any server-side controller. No other controller depends on PWAController.
