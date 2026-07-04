---
level: 128w
parent: ../controllers.32w.md
deeper: dashboard.256w.md
relates-to:
  - session.128w.md
  - weight.128w.md
  - macro.128w.md
source: controllers/dashboard.go
---

The dashboard controller renders the single authenticated home page at GET /dashboard. It requires a session and
redirects to /login if none exists. It fetches the 10 most recent sessions via `Sessions.GetRecentByUser`, all
body-weight entries via `BodyWeights.GetAllByUser`, all macro entries via `Macros.GetAllByUser`, and the macro goal via
`MacroGoals.Get`. Body-weight entries are converted to the user's preferred unit before a 3-day rolling average is
computed. Macro entries are grouped by calendar day and a 3-day summary is built. The template receives the username
from session, recent sessions, the weight average struct, and the macro summary struct. Errors from any repository call
are logged but do not stop rendering — the page degrades gracefully. Renders `dashboard.tpl`.
