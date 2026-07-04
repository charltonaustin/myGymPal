---
level: 128w
parent: root.32w.md
deeper: root-index.256w.md
relates-to:
  - ../controllers/auth.128w.md
source: views/index.tpl
---

Landing/home page for unauthenticated (and authenticated) users. Includes `partials/navbar.tpl`. Four informational
cards: "Quick Start" (link to example, two-step guide), "Full Setup" (four-step exercise-library workflow), "Tracking
Weight" (rolling average explanation, link to /weight and /settings), and "Tracking Macros" (food logging instructions
with serving-size and goal details).

No template variables used in the body (navbar receives `.LoggedIn`, `.ActivePage`). No forms. No AJAX. Loads
`offline-sync.js` and registers `/sw.js`.
