---
level: 128w
parent: auth.32w.md
deeper: auth-login.256w.md
relates-to:
  - ../controllers/auth.128w.md
source: views/auth/login.tpl
---

Renders the login page. Standalone HTML document with a minimal dark navbar (no auth links). A centered card holds a "
Log In" heading, an optional `.Error` danger alert, and a form that POSTs `username` and `password` to `/POST /login`. A
link below the form goes to `/register`.

Template variables: `.Error` (string, shown when login fails), `.Username` (string, repopulates the username field after
a failed attempt).

No partials included. No AJAX calls. Registers the PWA service worker (`/sw.js`) and loads `/static/offline-sync.js`.
Bootstrap 5.3.3 is loaded from CDN.
