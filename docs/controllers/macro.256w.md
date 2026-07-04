---
level: 256w
parent: ../controllers.32w.md
relates-to:
  - dashboard.256w.md
  - weight.256w.md
source: controllers/macro.go
---

## Routes

| Method | Path               | Handler  |
|--------|--------------------|----------|
| GET    | /macros            | Index    |
| POST   | /macros            | Create   |
| POST   | /macros/goals      | SaveGoal |
| POST   | /macros/:id        | Update   |
| POST   | /macros/:id/delete | Delete   |

## Auth requirement

All handlers check `c.GetSession("user_id")`; nil redirects to /login.

## Session keys

- Read: `user_id` (int64)

## Template variables — Index

- `c.Data["LoggedIn"]` = true
- `c.Data["ActivePage"]` = "macros"
- `c.Data["Days"]` = `[]macroDay` — entries grouped by calendar date
- `c.Data["DefaultDate"]` = today in "2006-01-02" format
- `c.Data["Goal"]` = `*models.MacroGoal` (nil if none set)
- `c.Data["Summary"]` = `*macroSummary` (nil if no entries; covers up to 3 most recent days)
- `c.Data["FoodHistoryJSON"]` = `template.JS` — JSON array of distinct food entries for autocomplete

## Template

- `macros/index.tpl`

## Repository calls

- `Macros.GetAllByUser(userID)` — Index, (DashboardController also calls this)
- `Macros.GetDistinctFoods(userID)` — Index (for food history autocomplete)
- `MacroGoals.Get(userID)` — Index
- `Macros.Create(userID, date, foodName, servingWeight, servingUnit, protein, carbs, fat)` — Create
- `Macros.Update(id, userID, foodName, servingWeight, servingUnit, protein, carbs, fat)` — Update
- `Macros.Delete(id, userID)` — Delete
- `MacroGoals.Upsert(userID, protein, carbs, fat)` — SaveGoal

## Serving unit validation

Accepted values: "oz", "ml", "fl oz"; anything else defaults to "g".

## Helper functions (also used by DashboardController)

- `groupMacrosByDay(entries)` — returns `[]macroDay`, each with summed protein/carbs/fat/calories
- `buildMacroSummary(days, goal)` — returns `*macroSummary` averaging up to 3 most recent days; computes
  percentage-of-goal rows; calories = 4P + 4C + 9F

## Flash messages

None — all redirects go directly to /macros.

## Redirect paths

- All POST success → /macros
- All POST failure → /macros (errors are logged, not flashed)

## Relationship to other controllers

`groupMacrosByDay` and `buildMacroSummary` are defined in `macro.go` but called directly by `DashboardController.Get` to
power the macro summary card.
