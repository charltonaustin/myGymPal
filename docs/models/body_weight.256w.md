---
level: 256w
parent: ../models.32w.md
relates-to:
  - ../controllers/weight.128w.md
source: models/body_weight.go, models/body_weight_repository.go
---

# BodyWeight (full reference)

## Struct fields

| Field      | Go type   | ORM tag / notes          |
|------------|-----------|--------------------------|
| ID         | int64     | `auto;pk`                |
| UserID     | int64     | FK to users.id           |
| Date       | time.Time | `type(date)` — date only |
| Weight     | float64   | numeric weight value     |
| WeightUnit | string    | `"lb"` or `"kg"`         |

No `CreatedAt` or `UpdatedAt` — the entry date is the `Date` field.

## Repository interface (BodyWeightRepository)

```go
Create(userID int64, date time.Time, weight float64, weightUnit string) (*BodyWeight, error)
GetAllByUser(userID int64) ([]*BodyWeight, error)
GetByID(id, userID int64) (*BodyWeight, error)
Update(id, userID int64, weight float64, weightUnit string) (*BodyWeight, error)
Delete(id, userID int64) error
```

## Notable behavior

- `GetAllByUser`: `QueryTable.Filter("UserID").OrderBy("-Date")` — newest first; no limit.
- `GetByID`: reads by primary key, checks `bw.UserID != userID`; returns `orm.ErrNoRows` on mismatch (consistent with
  MacroEntry; differs from other models that return `"not found"` string error).
- `Update`: calls `GetBodyWeightByID` first (enforces ownership), then `o.Update(bw, "weight", "weight_unit")` — partial
  update, preserves `Date` and `UserID`.
- `Delete`: calls `GetBodyWeightByID` first (enforces ownership), then `o.Delete(bw)`.
- No validation of `WeightUnit` in the model layer (unlike `User.UpdateWeightUnit`).

## ORM / SQL patterns

- `orm.RegisterModel(&BodyWeight{})` in `init()`.
- All operations use the default ORM instance; no raw SQL.
- `Update` uses named-field ORM update: `o.Update(bw, "weight", "weight_unit")`.
- No transactions — all operations are single-row.

## Relationships

- Belongs to `User` via `UserID`.
- No FK to other models; standalone time-series data.

## Usage

The weight controller (`controllers/weight.go`) handles all CRUD operations and chart data rendering. `ConvertWeight` in
`models/convert.go` can be used to normalize values across `"lb"`/`"kg"` for display.
