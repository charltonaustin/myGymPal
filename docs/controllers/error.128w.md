---
level: 128w
parent: ../controllers.32w.md
deeper: error.256w.md
relates-to:
  - default.128w.md
source: controllers/error.go
---

The error controller is a minimal catch-all page for unrecoverable failures. GET /error renders `error.tpl` with no
template data set and no session check. It is the redirect target used throughout the application when an unexpected
repository error occurs — for example when a session cannot be written during login, when a required program or phase
cannot be fetched, or when other data-layer operations fail in a way that makes the requested page impossible to render.
Because it performs no auth check, it is safely reachable even when the session is broken or absent. No flash messages,
no repository calls, and no `c.Data` keys are set.
