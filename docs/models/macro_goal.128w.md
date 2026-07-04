---
level: 128w
parent: ../models.32w.md
deeper: macro_goal.256w.md
relates-to:
  - macro_entry.128w.md
  - ../controllers/macro.128w.md
source: models/macro_goal.go, models/macro_goal_repository.go
---

# MacroGoal

`MacroGoal` stores a single per-user daily macro target: `Protein`, `Carbs`, and `Fat` in grams (float64). `UserID` has
a unique constraint, enforcing one goal record per user. An `UpdatedAt` timestamp is auto-managed.

The repository uses upsert semantics: `UpsertMacroGoal` tries to read an existing record by `UserID`, inserts if not
found, or updates the three macro values and `updated_at` if one exists. There is no separate create/update split
exposed to callers.

`Get` returns `nil, nil` (not an error) when no goal has been set yet — callers must handle the nil case.

The `MacroGoalRepository` interface has two methods: `Get(userID int64)` and
`Upsert(userID int64, protein, carbs, fat float64)`. The macro controller uses both to display and save daily targets.
