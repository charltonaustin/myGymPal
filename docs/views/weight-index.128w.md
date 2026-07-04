---
level: 128w
parent: weight.32w.md
deeper: weight-index.256w.md
relates-to:
  - ../controllers/weight.128w.md
source: views/weight/index.tpl
---

Body-weight tracking page. Includes `partials/navbar.tpl`. Shows an N-day rolling average (`.WeightAvg.Weight`,
`.WeightAvg.Unit`, `.WeightAvg.Days`) in a card. A log form POSTs to `/weight` with date, weight (number), and unit (
lb/kg select). Below, all entries are listed; pencil toggles an inline edit form (`POST /weight/:id`); trash submits
`POST /weight/:id/delete`. `.Success` flash auto-dismissed after 3 s. No AJAX.

Template variables: `.WeightAvg`, `.Entries` (slice with `.ID`, `.Weight`, `.WeightUnit`, `.Date`), `.DefaultDate`,
`.WeightUnit`, `.Success`.
