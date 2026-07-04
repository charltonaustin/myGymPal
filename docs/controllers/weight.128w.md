---
level: 128w
parent: ../controllers.32w.md
deeper: weight.256w.md
relates-to:
  - dashboard.128w.md
  - account.128w.md
source: controllers/weight.go
---

The weight controller manages the body-weight log. GET /weight lists all entries, converts each to the user's preferred
unit for display, computes a 3-day rolling average via `computeWeightAverage`, reads a flash success, and renders
`weight/index.tpl`. POST /weight logs a new entry with a date, weight value, and unit (lb or kg). POST /weight/:id
updates an existing entry's weight and unit. POST /weight/:id/delete removes an entry. All handlers are session-gated.
The `computeWeightAverage` helper is defined here and re-used by `DashboardController.Get`. Weight conversion is
performed via `models.ConvertWeight` before display, but stored values retain their original unit in the database.
