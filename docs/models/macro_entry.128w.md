---
level: 128w
parent: ../models.32w.md
deeper: macro_entry.256w.md
relates-to:
  - macro_goal.128w.md
  - ../controllers/macro.128w.md
source: models/macro_entry.go, models/macro_entry_repository.go
---

# MacroEntry

`MacroEntry` records a single food log entry for a user on a specific date. Fields include `FoodName`, `ServingWeight`,
`ServingUnit`, `Protein`, `Carbs`, and `Fat` (all floats). `CreatedAt` is auto-set; there is no `UpdatedAt` (the struct
exposes an `Update` path in the repository that patches named fields directly).

`GetDistinctFoodsByUser` uses a PostgreSQL `DISTINCT ON (food_name)` query to return the most recent entry per unique
food name — used for food autocomplete/quick-add UX.

Note: `MacroRepository` is defined inside `macro_entry_repository.go`, not in `interfaces.go`. Its interface includes
`Create`, `GetAllByUser`, `GetDistinctFoods`, `GetByID`, `Update`, and `Delete`. All reads are scoped by `userID`. The
macro controller is the primary consumer.
