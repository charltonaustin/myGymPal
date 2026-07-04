---
level: 256w
parent: root.32w.md
relates-to: []
source: views/offline.tpl
---

## Purpose

Offline fallback page that the PWA service worker (`/sw.js`) serves when a user navigates to a page that is not in the
cache and the network is unavailable.

## Partials Included

None. Has its own minimal dark navbar (brand link only).

## Template Variables

None. All content is static.

## Content

- Large emoji: signal-bars (📶) rendered as `class="display-1"`.
- `<h1>`: "You're Offline"
- `<p>`: "This page isn't available offline. Check your connection and try again."
- Button: "Retry" — calls `window.location.reload()` via inline `onclick`.

## JavaScript Behavior

Inline `onclick="window.location.reload()"` on Retry button. No service-worker registration (this page cannot register
one while offline in the same context it is served from).

## AJAX / Fetch

None.

## Flash Messages

None.

## Notes

This page must be pre-cached by the service worker at install time so it is available when all other caches miss. It
intentionally avoids loading any CDN assets that would require a network request.
