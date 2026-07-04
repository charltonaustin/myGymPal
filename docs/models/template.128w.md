---
level: 128w
parent: ../models.32w.md
deeper: template.256w.md
relates-to:
  - exercise.128w.md
  - session_exercise.128w.md
  - ../controllers/template.128w.md
source: models/template.go, models/template_repository.go
---

# Template

`Template` is a reusable workout blueprint with a `Name`, optional `Focus` string, and auto-managed `CreatedAt`/
`UpdatedAt` timestamps. Each template owns an ordered set of `TemplateExercise` rows storing exercise `Name` (
lowercased/trimmed), `IsBodyweight`, `IsTimeBased`, `Block`, and `SortOrder`.

`TemplateExercise` stores only the exercise identity (name + flags); goal weight, reps, and seconds come from the
Exercise library and Phase config at render time.

Both `Create` and `Update` run inside transactions and replace exercises atomically. `Update` deletes all existing
`template_exercises` for the template, then inserts the new list. At least one exercise is required; all exercise names
must be non-empty.

The `TemplateRepository` interface provides `Create`, `Update`, `GetAll` (ordered by name), `GetByID`, and `Delete`. The
template controller is the primary consumer.
