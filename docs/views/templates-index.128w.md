---
level: 128w
parent: templates.32w.md
deeper: templates-index.256w.md
relates-to:
  - ../controllers/templates.128w.md
source: views/templates/index.tpl
---

Workout Templates list page. Includes `partials/navbar.tpl`. Renders each template as a link row showing name and
optional focus text. A trash icon per row opens `#deleteModal`, which sets the hidden form's action to
`/templates/:id/delete` and submits via POST. `.Success` flash auto-dismisses after 3 s.

Template variables: `.Templates` (slice of structs with `.ID`, `.Name`, `.Focus`), `.Success` (string). No AJAX calls.
Registers `/sw.js` PWA service worker.
