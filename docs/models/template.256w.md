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

## Struct fields — TemplateCircuit

A circuit is a named, ordered group of exercises inside a template. It runs `Rounds` times, and
`TransitionSeconds` separates one exercise from the next.

| Field             | Go type | Notes                              |
|-------------------|---------|------------------------------------|
| ID                | int64   | `auto;pk`                          |
| TemplateID        | int64   | FK to templates.id, `ON DELETE CASCADE` |
| Name              | string  | trimmed; non-empty enforced        |
| Rounds            | int     | `>= 1` (CHECK + Go)                |
| TransitionSeconds | int     | `>= 0` (CHECK + Go)                |
| SortOrder         | int     | display order                      |

## Struct fields — TemplateExercise

| Field        | Go type | Notes                                              |
|--------------|---------|----------------------------------------------------|
| ID           | int64   | `auto;pk`                                          |
| TemplateID   | int64   | FK to templates.id                                 |
| Name         | string  | lowercased+trimmed; non-empty enforced             |
| IsBodyweight | bool    |                                                    |
| IsTimeBased  | bool    |                                                    |
| Block        | string  | `"main"`, `"abs"`, `"cardio"`, `"stretch"`         |
| SortOrder    | int     | display order                                      |
| CircuitID    | *int64  | nil = not in a circuit; FK, `ON DELETE SET NULL`   |
| WorkSeconds  | int     | `>= 0`; only meaningful when `CircuitID` is set    |

`CircuitID IS NULL` is the normal, non-circuit exercise, and is what every row predating migration `000035` is.
Deleting a circuit returns its exercises to the template as loose exercises rather than deleting them.

## Input helper structs

```go
type TemplateCircuitInput struct {
    Name              string
    Rounds            int
    TransitionSeconds int
    SortOrder         int
}

type TemplateExerciseInput struct {
    Name         string
    IsBodyweight bool
    IsTimeBased  bool
    Block        string
    SortOrder    int
    CircuitIndex int // index into the circuits slice; NoCircuit (-1) = loose
    WorkSeconds  int
}

const NoCircuit = -1
```

**`CircuitIndex`'s zero value, 0, means "the first circuit" — not "no circuit".** An input built without setting
it explicitly is filed under circuit 0. Non-circuit inputs must say `NoCircuit`. The repository rejects an index
that points at no circuit, so the mistake fails loudly rather than silently mis-filing an exercise, but the sharp
edge is real.

## Repository interface (TemplateRepository)

```go
Create(name, focus string, circuits []TemplateCircuitInput, exercises []TemplateExerciseInput) (*Template, error)
Update(id int64, name, focus string, circuits []TemplateCircuitInput, exercises []TemplateExerciseInput) (*Template, error)
GetAll() ([]*Template, error)
GetByID(id int64) (*Template, []*TemplateExercise, error)
GetCircuits(templateID int64) ([]*TemplateCircuit, error)
Delete(id int64) error
```

`GetByID` deliberately does **not** return circuits. The session and exercise controllers call it for the exercise
list alone; circuits are a separate read so those callers are unaffected by the feature.

## Notable behavior

- **Index, not id.** A circuit and the exercises inside it arrive in the same form submit, so an exercise cannot
  reference a `circuit_id` that does not exist yet. `TemplateExerciseInput.CircuitIndex` refers to a position in the
  `circuits` slice, and `insertTemplateBody` resolves it to a real id inside the existing transaction — circuits are
  inserted first, their ids kept, and each exercise's index mapped before it is inserted. There is no two-phase save
  and no insert-then-patch, so a half-written template is not reachable.
- **An out-of-range `CircuitIndex` is an error, not a NULL.** Writing NULL would quietly drop the exercise out of its
  circuit and read as data loss with nothing pointing back at the cause. The guard is also load-bearing rather than
  merely defensive: without it, the index into the id slice panics.
- `Create`: validates name non-empty, at least one exercise with a non-empty name, every circuit named, `Rounds >= 1`,
  seconds non-negative; normalizes exercise names with `strings.ToLower(strings.TrimSpace(...))`. Transactional.
- `Update`: same validation, then in one transaction updates the `Template` row, deletes all `template_exercises` and
  then all `template_circuits` for the template (exercises first — they hold the FK), and re-inserts both from the
  input. **`Update` is a delete-and-reinsert, so a field it does not copy is not preserved, it is destroyed.** That is
  why `Create` and `Update` share a single `insertTemplateBody`: the two paths cannot drift apart. `IsTimeBased` was
  once dropped exactly this way and reverted every time-based exercise to weighted on save.
  `TestUpdateTemplate_PreservesExerciseType` and `TestUpdateTemplate_PreservesCircuitsAndWorkSeconds` guard it.
- `GetAll`: `QueryTable.OrderBy("Name")` — alphabetical ordering.
- `GetByID`: reads template, then exercises ordered by `SortOrder`.
- `GetCircuits`: circuits for a template, ordered by `SortOrder`.
- `Delete`: deletes by primary key; `template_circuits` cascades away with the template.

## Key invariant

Template exercises store only identity — name, `is_bodyweight`, `is_time_based`, `block`, `sort_order` — plus, for a
circuit member, how long it is worked (`work_seconds`). Goal weight/reps/seconds still come from the Exercise library
and Phase config at render time (per CLAUDE.md); `work_seconds` is not a goal, it is the length of a circuit interval.

## Relationships

- Has many `TemplateCircuit` rows (replaced atomically on update, cascade-deleted with the template).
- Has many `TemplateExercise` rows (replaced atomically on update), each optionally belonging to one circuit.
- Template exercises reference Exercise names, not IDs — loose coupling to the Exercise library.
