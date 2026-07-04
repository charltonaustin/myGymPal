---
level: 128w
parent: macros.32w.md
deeper: macros-index.256w.md
relates-to:
  - ../controllers/macros.128w.md
source: views/macros/index.tpl
---

Single-page macro tracker. Includes `partials/navbar.tpl`. Shows an N-day average summary (protein, carbs, fat, calories
vs goals, colour-coded green/red). Two forms on page: daily goals (`POST /macros/goals`) and log-food (`POST /macros`).
Food-name input has autocomplete from `.FoodHistoryJSON` that pre-fills serving and macro fields. Per-day sections list
logged entries; pencil toggles an inline edit form per entry (`POST /macros/:id`); trash submits
`POST /macros/:id/delete`. Template variables: `.Summary`, `.Goal`, `.Days`, `.DefaultDate`, `.FoodHistoryJSON`.
