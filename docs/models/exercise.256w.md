---
level: 256w
parent: ../models.32w.md
relates-to:
  - session_exercise.128w.md
  - template.128w.md
  - ../controllers/exercise.128w.md
source: models/exercise.go, models/exercise_repository.go
---

# Exercise (full reference)

## Struct fields

| Field        | Go type   | ORM tag / notes                                      |
|--------------|-----------|------------------------------------------------------|
| ID           | int64     | `auto;pk`                                            |
| UserID       | int64     | FK to users.id                                       |
| Name         | string    | always lowercase+trimmed                             |
| IsBodyweight | bool      |                                                      |
| GoalWeight   | float64   |                                                      |
| WeightUnit   | string    | `"lb"` or `"kg"`                                     |
| IsTimeBased  | bool      |                                                      |
| GoalSeconds  | int       |                                                      |
| GoalRepMin   | int       | overrides phase rep range for bodyweight when > 0    |
| GoalRepMax   | int       | overrides phase rep range for bodyweight when > 0    |
| DefaultBlock | string    | `"main"` (default), `"abs"`, `"cardio"`, `"stretch"` |
| CreatedAt    | time.Time | `auto_now_add`                                       |
| UpdatedAt    | time.Time | `auto_now`                                           |

## Repository interface (ExerciseRepository)

```go
Create(userID int64, name string, isBodyweight bool, goalWeight float64, weightUnit string, isTimeBased bool, goalSeconds int, goalRepMin int, goalRepMax int, defaultBlock string) (*Exercise, error)
GetAllByUser(userID int64) ([]*Exercise, error)
GetByID(id, userID int64) (*Exercise, error)
GetByName(userID int64, name string) (*Exercise, error)
Update(id, userID int64, name string, isBodyweight bool, goalWeight float64, weightUnit string, isTimeBased bool, goalSeconds int, goalRepMin int, goalRepMax int, defaultBlock string) (*Exercise, error)
UpdateGoalWeight(id int64, goalWeight float64, weightUnit string) error
Delete(id, userID int64) error
```

## Notable behavior

- `Create` and `Update` both normalize name: `strings.ToLower(strings.TrimSpace(name))`; error if result is empty.
- `DefaultBlock` is validated by `validDefaultBlock()`: accepts `"abs"`, `"cardio"`, `"stretch"`; all others become
  `"main"`.
- `GetByName` uses raw SQL: `LOWER(TRIM(name)) = LOWER(TRIM(?))` — case-insensitive, whitespace-insensitive lookup.
- `UpdateGoalWeight` does a targeted ORM update of only `GoalWeight`, `WeightUnit`, `UpdatedAt` fields — used by the
  auto-goal-promotion logic after 3rd set at or above goal reps.
- `GetAllByUser`: `QueryTable.Filter("UserID").OrderBy("Name")` — alphabetical.
- `GetByID` and `Delete` enforce user ownership by comparing `ex.UserID != userID`.

## ORM / SQL patterns

- `orm.RegisterModel(&Exercise{})` in `init()`.
- `GetByName`: raw SQL query returning into `[]*Exercise` slice, returns first element.
- `Update`: calls `GetExerciseByID` first (enforces ownership), mutates struct, then `o.Update(ex)` (full update).
- `UpdateGoalWeight`: `o.Update(ex, "GoalWeight", "WeightUnit", "UpdatedAt")` (partial update).

## Key invariants

- Exercise names are always lowercase+trimmed (enforced in model layer).
- `GoalRepMin`/`GoalRepMax` on the exercise override the Phase rep range for bodyweight exercises at session render
  time.

## Relationships

- Belongs to `User` via `UserID`.
- Referenced by `SessionExercise` by name (not by FK).
- Referenced by `TemplateExercise` by name.
