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

A template may also own `TemplateCircuit` rows: named, ordered groups that run for `Rounds` with a
`TransitionSeconds` gap between exercises. An exercise inside a circuit has `CircuitID` set and a `WorkSeconds`
duration of its own; `CircuitID IS NULL` is the normal, non-circuit exercise and is what every row predating
migration `000035` is. Deleting a circuit returns its exercises to the template rather than deleting them.

`TemplateExercise` stores only the exercise identity (name + flags, plus a circuit member's work duration); goal
weight, reps, and seconds come from the Exercise library and Phase config at render time.

Both `Create` and `Update` run inside transactions and replace circuits and exercises atomically. Because a circuit
and its exercises are created in the same submit, `TemplateExerciseInput.CircuitIndex` references a circuit by
position in the input slice and the repository resolves it to a real id inside the transaction. An index pointing at
no circuit is an error, not a silent NULL. **`CircuitIndex`'s zero value means "circuit 0", not "no circuit" — loose
inputs must say `NoCircuit` explicitly.**

Because `Update` re-inserts rather than patches, **every** field must be copied from the input on both paths — a
field `Create` sets and `Update` omits is not merely ignored on edit, it is cleared. `IsTimeBased` was dropped this
way. `Create` and `Update` now share one `insertTemplateBody` so the two paths cannot drift;
`TestUpdateTemplate_PreservesExerciseType` and `TestUpdateTemplate_PreservesCircuitsAndWorkSeconds` guard it.

The `TemplateRepository` interface provides `Create`, `Update`, `GetAll` (ordered by name), `GetByID`, `GetCircuits`,
and `Delete`. `GetByID` returns exercises but not circuits, because the session and exercise controllers call it for
the exercise list alone. The template controller is the primary consumer.
