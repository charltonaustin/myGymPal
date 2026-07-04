---
level: 256w
parent: ../models.32w.md
relates-to:
  - macro_entry.128w.md
  - ../controllers/macro.128w.md
source: models/macro_goal.go, models/macro_goal_repository.go
---

# MacroGoal (full reference)

## Struct fields

| Field     | Go type   | ORM tag / notes              |
|-----------|-----------|------------------------------|
| ID        | int64     | `auto;pk`                    |
| UserID    | int64     | `unique` — one goal per user |
| Protein   | float64   | daily target grams           |
| Carbs     | float64   | daily target grams           |
| Fat       | float64   | daily target grams           |
| UpdatedAt | time.Time | `auto_now`                   |

No `CreatedAt` field.

## Repository interface (MacroGoalRepository)

```go
Get(userID int64) (*MacroGoal, error)
Upsert(userID int64, protein, carbs, fat float64) (*MacroGoal, error)
```

## Notable behavior

- `Get`: reads by `UserID` field (`o.Read(g, "UserID")`). Returns `nil, nil` when `orm.ErrNoRows` — callers must handle
  the nil case gracefully (no goal yet set).
- `Upsert`: reads by `UserID` first:
    - On `ErrNoRows`: sets fields and calls `o.Insert(g)`.
    - On success (existing row): sets `Protein`, `Carbs`, `Fat`, then calls
      `o.Update(g, "protein", "carbs", "fat", "updated_at")` — partial update, preserves `UserID` and `ID`.
    - Other read errors are returned immediately.
- No delete method — macro goals are permanent once set (only updatable).
- No create-only method — `Upsert` is the sole write path.

## ORM / SQL patterns

- `orm.RegisterModel(&MacroGoal{})` in `init()`.
- All queries use the default ORM instance; no raw SQL or transactions.
- `o.Read(g, "UserID")` uses the `UserID` unique index as a lookup key instead of primary key.
- `o.Update(g, "protein", "carbs", "fat", "updated_at")` — named-field partial update prevents overwriting `UserID`.

## Relationships

- Belongs to `User` via `UserID` (unique FK — one-to-one).
- Logically paired with `MacroEntry` records for daily macro tracking and comparison in the macro controller.
