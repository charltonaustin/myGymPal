---
level: 256w
parent: ../models.32w.md
relates-to:
  - exercise.128w.md
  - session_exercise.128w.md
  - ../controllers/template.128w.md
source: models/template.go, models/template_repository.go
---

# Template (full reference)

## Struct fields — Template

| Field     | Go type   | ORM tag / notes                 |
|-----------|-----------|---------------------------------|
| ID        | int64     | `auto;pk`                       |
| Name      | string    | non-empty enforced              |
| Focus     | string    | optional description/focus area |
| CreatedAt | time.Time | `auto_now_add`                  |
| UpdatedAt | time.Time | `auto_now`                      |

## Struct fields — TemplateExercise

| Field        | Go type | Notes                                      |
|--------------|---------|--------------------------------------------|
| ID           | int64   | `auto;pk`                                  |
| TemplateID   | int64   | FK to templates.id                         |
| Name         | string  | lowercased+trimmed; non-empty enforced     |
| IsBodyweight | bool    |                                            |
| IsTimeBased  | bool    |                                            |
| Block        | string  | `"main"`, `"abs"`, `"cardio"`, `"stretch"` |
| SortOrder    | int     | display order                              |

## Input helper struct

```go
type TemplateExerciseInput struct {
    Name         string
    IsBodyweight bool
    IsTimeBased  bool
    Block        string
    SortOrder    int
}
```

## Repository interface (TemplateRepository)

```go
Create(name, focus string, exercises []TemplateExerciseInput) (*Template, error)
Update(id int64, name, focus string, exercises []TemplateExerciseInput) (*Template, error)
GetAll() ([]*Template, error)
GetByID(id int64) (*Template, []*TemplateExercise, error)
Delete(id int64) error
```

## Notable behavior

- `Create`: validates name non-empty, at least one exercise with non-empty name; normalizes all exercise names with
  `strings.ToLower(strings.TrimSpace(...))`. Transactional: inserts template then all exercises; rolls back on failure.
- `Update`: validates same rules as `Create`, then in a transaction: updates the `Template` row, deletes all existing
  `template_exercises` via raw SQL (`DELETE FROM template_exercises WHERE template_id = ?`), and re-inserts the new
  exercise list. Note: `IsTimeBased` is not re-inserted in the update path in the current implementation.
- `GetAll`: `QueryTable.OrderBy("Name")` — alphabetical ordering.
- `GetByID`: reads template, then `QueryTable.Filter("TemplateID").OrderBy("SortOrder")` for exercises.
- `Delete`: deletes by primary key; cascade behavior depends on DB constraints.

## Key invariant

Template exercises store only name + is_bodyweight; goal weight/reps/seconds come from the Exercise library and Phase
config at render time (per CLAUDE.md).

## Relationships

- Has many `TemplateExercise` rows (replaced atomically on update).
- Template exercises reference Exercise names, not IDs — loose coupling to the Exercise library.
