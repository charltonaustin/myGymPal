---
level: 256w
parent: ../controllers.32w.md
relates-to:
  - session.256w.md
  - template.256w.md
source: controllers/exercise.go
---

## Routes

| Method | Path                    | Handler               |
|--------|-------------------------|-----------------------|
| GET    | /exercises              | Index                 |
| GET    | /exercises/new          | New                   |
| POST   | /exercises/new          | Create                |
| GET    | /exercises/:id/edit     | Edit                  |
| POST   | /exercises/:id/edit     | Update                |
| POST   | /exercises/:id/delete   | Delete                |
| POST   | /exercises/goal-weight  | UpdateGoalWeightJSON  |
| POST   | /exercises/goal-reps    | UpdateGoalRepsJSON    |
| POST   | /exercises/goal-seconds | UpdateGoalSecondsJSON |

## Auth requirement

All handlers check `c.GetSession("user_id")`; nil returns 302 to /login for HTML handlers or JSON
`{"error": "unauthenticated"}` for AJAX handlers.

## Session keys

- Read: `user_id` (int64)

## Template variables — Index

- `c.Data["LoggedIn"]` = true
- `c.Data["ActivePage"]` = "exercises"
- `c.Data["Exercises"]` = `[]*models.Exercise` (goal weights converted to preferred unit)
- `c.Data["Success"]` = flash success string (if present)

## Template variables — New / Create failure

- `c.Data["WeightUnit"]` = user's preferred unit
- `c.Data["ExWeightUnit"]` = exercise-specific unit (lb or kg)
- `c.Data["DefaultBlock"]` = ""
- `c.Data["ShowDefaultBlock"]` = true
- `c.Data["Name"]`, `c.Data["IsBodyweight"]`, `c.Data["IsTimeBased"]`
- `c.Data["GoalWeight"]`, `c.Data["GoalHours"]`, `c.Data["GoalMinutes"]`, `c.Data["GoalSecsRemainder"]`
- `c.Data["GoalRepMin"]`, `c.Data["GoalRepMax"]`
- `c.Data["Error"]` (on failure)

## Template variables — Edit / Update failure

Same as New/Create plus `c.Data["Exercise"]` = `*models.Exercise`

## Templates

- `exercises/index.tpl`
- `exercises/new.tpl`
- `exercises/edit.tpl`

## Repository calls

- `Exercises.GetAllByUser(userID)` — Index
- `Exercises.Create(...)` — Create
- `Exercises.GetByID(id, userID)` — Edit, Update
- `Exercises.Update(...)` — Update, UpdateGoalWeightJSON, UpdateGoalRepsJSON, UpdateGoalSecondsJSON
- `Exercises.Delete(id, userID)` — Delete
- `Exercises.GetByName(userID, name)` — all three AJAX endpoints
- `Users.GetByID(userID)` — all handlers that need preferred unit

## Flash messages

- `flash.Success("%s added to your exercise library.", name)` — Create success
- `flash.Success("Exercise updated.")` — Update success
- `flash.Success("Exercise deleted.")` — Delete success

## AJAX JSON responses

- Success: `{"ok": true, "goal_weight": float, "weight_unit": string}` etc.
- Failure: `{"error": "message"}`

## Key invariant

Exercise names are matched via `Exercises.GetByName` which applies `LOWER(TRIM(...))` on both sides; goal
weight/reps/seconds updates from AJAX calls use `Exercises.Update` (not `UpdateGoalWeight`) to preserve all other
fields.

## Relationship to other controllers

`exerciseLibraryJSON` is called by `TemplateController.New/Edit` and `SessionController.Show` to populate JavaScript
autocomplete data.
