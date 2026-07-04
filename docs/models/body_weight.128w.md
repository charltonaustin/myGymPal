---
level: 128w
parent: ../models.32w.md
deeper: body_weight.256w.md
relates-to:
  - ../controllers/weight.128w.md
source: models/body_weight.go, models/body_weight_repository.go
---

# BodyWeight

`BodyWeight` records a dated body-weight measurement for a user. Fields are `UserID`, `Date` (date-only), `Weight` (
float64), and `WeightUnit` (`"lb"` or `"kg"`). No timestamps beyond `Date`.

All reads are user-scoped. `GetAllByUser` returns entries newest-first (`OrderBy("-Date")`). `GetByID` reads by primary
key then checks ownership, returning `orm.ErrNoRows` on mismatch. `Update` patches only `weight` and `weight_unit`
fields after fetching the existing record. `Delete` fetches first to enforce ownership.

The `BodyWeightRepository` interface provides `Create`, `GetAllByUser`, `GetByID`, `Update`, and `Delete`. The weight
controller is the sole consumer, using this data for display and charting of body-weight trends over time.
