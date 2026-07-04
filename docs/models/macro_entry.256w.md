---
level: 256w
parent: ../models.32w.md
relates-to:
  - macro_goal.128w.md
  - ../controllers/macro.128w.md
source: models/macro_entry.go, models/macro_entry_repository.go
---

# MacroEntry (full reference)

## Struct fields

| Field         | Go type   | ORM tag / notes                              |
|---------------|-----------|----------------------------------------------|
| ID            | int64     | `auto;pk`                                    |
| UserID        | int64     | FK to users.id                               |
| Date          | time.Time | `type(date)` — date only                     |
| FoodName      | string    |                                              |
| ServingWeight | float64   | quantity of the serving                      |
| ServingUnit   | string    | unit for serving weight (e.g. `"g"`, `"oz"`) |
| Protein       | float64   | grams                                        |
| Carbs         | float64   | grams                                        |
| Fat           | float64   | grams                                        |
| CreatedAt     | time.Time | `auto_now_add`                               |

No `UpdatedAt` on the struct; the `Update` repository method patches fields in place using ORM's named-field update.

## Repository interface (MacroRepository — defined in macro_entry_repository.go)

```go
Create(userID int64, date time.Time, foodName string, servingWeight float64, servingUnit string, protein, carbs, fat float64) (*MacroEntry, error)
GetAllByUser(userID int64) ([]*MacroEntry, error)
GetDistinctFoods(userID int64) ([]*MacroEntry, error)
GetByID(id, userID int64) (*MacroEntry, error)
Update(id, userID int64, foodName string, servingWeight float64, servingUnit string, protein, carbs, fat float64) (*MacroEntry, error)
Delete(id, userID int64) error
```

Note: `MacroRepository` is declared locally in `macro_entry_repository.go`, not in `interfaces.go`.

## Notable behavior

- `GetAllByUser`: `QueryTable.Filter("UserID").OrderBy("-Date", "-CreatedAt")` — newest first.
- `GetDistinctFoods`: raw PostgreSQL SQL using `DISTINCT ON (food_name)`, ordered by `food_name, created_at DESC`.
  Returns one entry per unique food name (the most recent). Used for autocomplete.
- `GetByID`: reads by primary key, then checks `e.UserID != userID`; returns `orm.ErrNoRows` on mismatch (unlike most
  other models that return `"not found"` string error).
- `Update`: fetches via `GetByID` first (enforces ownership), then patches `food_name`, `serving_weight`,
  `serving_unit`, `protein`, `carbs`, `fat` using ORM named-field update.
- `Delete`: fetches via `GetByID` first, then `o.Delete(e)`.

## ORM / SQL patterns

- `orm.RegisterModel(&MacroEntry{})` in `init()`.
- `GetDistinctFoods` uses raw SQL — the only place in this model file that bypasses the ORM query builder.

## Relationships

- Belongs to `User` via `UserID`.
- Logically paired with `MacroGoal` for daily tracking comparisons.
