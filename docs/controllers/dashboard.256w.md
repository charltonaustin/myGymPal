---
level: 256w
parent: ../controllers.32w.md
relates-to:
  - session.256w.md
  - weight.256w.md
  - macro.256w.md
source: controllers/dashboard.go
---

## Routes

| Method | Path       | Handler |
|--------|------------|---------|
| GET    | /dashboard | Get     |

## Auth requirement

Checks `c.GetSession("user_id")`; nil redirects to /login.

## Session keys

- Read: `user_id` (int64), `username` (any, passed directly to template)

## Repository calls

- `Sessions.GetRecentByUser(userID, 10)` — fetches 10 most recent sessions
- `BodyWeights.GetAllByUser(userID)` — all body-weight entries
- `Macros.GetAllByUser(userID)` — all macro entries
- `MacroGoals.Get(userID)` — single macro goal record
- `Users.GetByID(userID)` — reads `WeightUnit` for unit conversion

## Data processing

1. Body-weight entries are each passed through `models.ConvertWeight(e.Weight, e.WeightUnit, preferredUnit)` in-place.
2. `computeWeightAverage(entries, preferredUnit)` averages the most recent 3 entries into a
   `weightAverage{Days, Weight, Unit}` struct.
3. `groupMacrosByDay(entries)` groups macro entries into `[]macroDay` keyed by date string.
4. `buildMacroSummary(days, goal)` averages the most recent 3 days of macros and computes percentage-of-goal rows for
   protein, carbs, fat, and calories.

## Template variables

- `c.Data["LoggedIn"]` = true
- `c.Data["ActivePage"]` = "dashboard"
- `c.Data["Username"]` = `c.GetSession("username")`
- `c.Data["RecentSessions"]` = `[]*models.RecentSession` (up to 10)
- `c.Data["WeightAvg"]` = `*weightAverage` (nil if no entries)
- `c.Data["MacroSummary"]` = `*macroSummary` (nil if no entries)

## Template

- `dashboard.tpl`

## Flash messages

None.

## Error handling

All repository errors are logged via `logs.Error` but do not abort rendering. The page displays partial data gracefully
when any fetch fails.

## Relationship to other controllers

The dashboard aggregates summaries produced by WeightController (`computeWeightAverage`) and MacroController (
`groupMacrosByDay`, `buildMacroSummary`) — these helper functions live in `weight.go` and `macro.go` respectively but
are called from `dashboard.go`. The `RecentSessions` list links to individual sessions managed by SessionController.
