---
level: 256w
parent: root.32w.md
relates-to: []
source: views/error.tpl
---

## Purpose

Generic catch-all error page rendered when the application encounters an unhandled error. Provides a user-friendly
fallback with a link to navigate home.

## Partials Included

None. Renders its own minimal dark navbar (brand link only) without the shared `partials/navbar.tpl`.

## Template Variables

None. All content is static.

## Content

- `<h1>`: "Something went wrong"
- `<p>`: "An unexpected error occurred. Please try again."
- Button: "Go home" → `/`

## JavaScript Behavior

Registers `/sw.js` PWA service worker. No other JS.

## AJAX / Fetch

None.

## Flash Messages

None.

## Notes

Does not use the shared navbar partial, so it does not require `.LoggedIn` or `.ActivePage`. Kept intentionally minimal
to avoid further errors during rendering.
