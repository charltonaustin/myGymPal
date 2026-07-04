---
level: 128w
parent: partials.32w.md
deeper: partials-navbar.256w.md
relates-to:
  - ../controllers/auth.128w.md
source: views/partials/navbar.tpl
---

Shared responsive dark navbar included in every authenticated page. Collapses into a hamburger on small screens. When
`.LoggedIn` is true, renders links to Home, Programs, Templates, Exercises, Weight, Macros, Dashboard, Settings, and a "
Log Out" button. When false, shows only Home, Log In, and Create Account buttons. Active page is indicated by
`fw-semibold active` class, controlled by `.ActivePage` string. When `isDev` template function returns true, a red "
DEVELOPMENT" banner is rendered below the navbar.
