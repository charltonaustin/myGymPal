---
level: 256w
parent: ../models.32w.md
relates-to:
  - program.128w.md
  - session.128w.md
  - ../controllers/program.128w.md
source: models/phase.go, models/phase_repository.go
---

# Phase (full reference)

## Struct fields

| Field       | Go type | ORM tag / notes                     |
|-------------|---------|-------------------------------------|
| ID          | int64   | `auto;pk`                           |
| ProgramID   | int64   | FK to programs.id                   |
| PhaseNumber | int     | 1-indexed within program            |
| RepMin      | int     | minimum reps for sets in this phase |
| RepMax      | int     | maximum reps for sets in this phase |
| DefaultSets | int     | default number of sets per exercise |
| RestSeconds | int     | default rest period between sets    |

No `CreatedAt`/`UpdatedAt` — phases are structural config, not event records.

## Helper struct

```go
type PhaseUpdate struct {
    PhaseNumber int
    RepMin      int
    RepMax      int
    DefaultSets int
    RestSeconds int
}
```

## Repository interface (PhaseRepository)

```go
GetByProgram(programID int64) ([]*Phase, error)
UpdateRepRanges(programID int64, updates []PhaseUpdate) error
```

## Notable behavior

- `GetByProgram` orders phases by `PhaseNumber` ascending.
- `UpdateRepRanges` validates all updates before writing: `RepMin > 0`, `RepMax >= RepMin`, `DefaultSets > 0`. Fails
  fast on the first invalid entry without partial writes (loop exits before any DB calls if validation fails on first
  item, but validation loop runs sequentially so a later invalid entry aborts before its DB call).
- There is no single-phase `GetByID` — phases are always read as a set.
- No `Create` or `Delete` on `PhaseRepository`; phases are created inside `CreateProgram`'s transaction and have no
  separate delete path.

## ORM / SQL patterns

- `orm.RegisterModel(&Phase{})` in `init()`.
- `GetByProgram`: `QueryTable.Filter("ProgramID").OrderBy("PhaseNumber")`.
- `UpdateRepRanges`: one `QueryTable.Filter("ProgramID").Filter("PhaseNumber").Update(orm.Params{...})` per entry in the
  updates slice.

## Relationships

- Belongs to `Program` via `ProgramID`.
- Phase rep ranges are used at session render time to set per-exercise goals for non-bodyweight exercises; bodyweight
  exercises use exercise-level `GoalRepMin`/`GoalRepMax` when set.
