---
level: 128w
parent: ../controllers.32w.md
deeper: default.256w.md
relates-to:
  - auth.128w.md
  - dashboard.128w.md
source: controllers/default.go
---

The default controller (`MainController`) handles the public-facing landing pages. GET / sets `LoggedIn` based on
whether `user_id` is present in the session (no redirect — unauthenticated users may view the page), sets `ActivePage`
to "home", and renders `index.tpl`. GET /example similarly checks session presence for the `LoggedIn` flag, sets
`ActivePage` to an empty string, and renders `example.tpl`. Neither handler performs any repository calls or enforces
authentication. These are the only routes accessible to unauthenticated visitors without a redirect, aside from /login,
/register, /error, and the PWA asset routes.
