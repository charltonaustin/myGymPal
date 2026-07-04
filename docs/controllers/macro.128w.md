---
level: 128w
parent: ../controllers.32w.md
deeper: macro.256w.md
relates-to:
  - dashboard.128w.md
  - weight.128w.md
source: controllers/macro.go
---

The macro controller manages daily food and macro-nutrient logging. GET /macros lists all entries grouped by calendar
day with a 3-day average summary, the current macro goal, and a food-history JSON blob for autocomplete. POST /macros
creates a new macro entry for a given date. POST /macros/:id updates an existing entry. POST /macros/:id/delete removes
an entry. POST /macros/goals upserts the user's daily macro goal (protein, carbs, fat). All handlers are session-gated.
The controller also contains the `groupMacrosByDay` and `buildMacroSummary` helper functions, which are re-used by
`DashboardController.Get` to render the macro summary widget. Calories are computed as 4×protein + 4×carbs + 9×fat. The
summary covers the most recent 3 logged days.
