---
level: 256w
parent: ../controllers.32w.md
relates-to:
  - session.256w.md
  - exercise.256w.md
source: controllers/template.go
---

## Routes

| Method | Path                  | Handler |
|--------|-----------------------|---------|
| GET    | /templates            | Index   |
| GET    | /templates/new        | New     |
| POST   | /templates/new        | Create  |
| GET    | /templates/:id        | Show    |
| POST   | /templates/:id        | Update  |
| GET    | /templates/:id/edit   | Edit    |
| POST   | /templates/:id/delete | Delete  |

## Auth requirement

All handlers check `c.GetSession("user_id")`; nil redirects to /login.

## Session keys

- Read: `user_id` (int64)

## Template variables — Index

- `c.Data["LoggedIn"]` = true
- `c.Data["ActivePage"]` = "templates"
- `c.Data["Templates"]` = `[]*models.Template`
- `c.Data["Success"]` = flash success (if present)

## Template variables — New / Create failure

- `c.Data["Exercises"]` = `[]exerciseForm` (one default entry with Block="main" on GET)
- `c.Data["ExerciseLibraryJSON"]` = `template.JS` (safe embedded JSON for autocomplete)
- `c.Data["Error"]`, `c.Data["Name"]`, `c.Data["Focus"]` (on failure)

## Template variables — Show

- `c.Data["Template"]` = `*models.Template`
- `c.Data["ExerciseBlocks"]` = `[]templateExerciseBlock` (grouped by block)
- `c.Data["Success"]` = flash success (if present)

## Template variables — Edit / Update failure

- `c.Data["Template"]`, `c.Data["Name"]`, `c.Data["Focus"]`
- `c.Data["Exercises"]` = `[]exerciseForm`
- `c.Data["ExerciseLibraryJSON"]` = `template.JS`
- `c.Data["Error"]` (on failure)

## Templates (HTML)

- `templates/index.tpl`
- `templates/new.tpl`
- `templates/show.tpl`
- `templates/edit.tpl`

## Repository calls

- `Templates.GetAll()` — Index
- `Templates.Create(name, focus, []TemplateExerciseInput)` — Create
- `Templates.GetByID(id)` — Show, Edit, Update
- `Templates.Update(id, name, focus, []TemplateExerciseInput)` — Update
- `Templates.Delete(id)` — Delete

## Exercise form parsing (Create / Update)

Reads `exercise_count` then iterates `0..count-1`, reading `exercise_name_N`, `is_bodyweight_N`, `is_time_based_N`,
`block_N`. Blank names are skipped. Block is validated via `validBlock()` (defaults to "main"). `SortOrder` is set from
the filtered index position.

## Flash messages

- `flash.Success("Template created.")` — Create success (redirects to /templates/:id)
- `flash.Success("Template updated.")` — Update success (redirects to /templates/:id)
- `flash.Success("Template deleted.")` — Delete success (redirects to /templates)

## Shared utilities defined in template.go

- `blockOrder` = ["main", "abs", "cardio", "stretch"]
- `blockLabels` = map of block key → display label
- `validBlock(b)` — sanitizes block string to one of the four valid values
- `groupTemplateExercises(exercises)` — returns `[]templateExerciseBlock` (used by Show)
- `groupSessionExercises(exercises)` — returns `[]sessionExerciseBlock` (used by SessionController.Show)

## Relationship to other controllers

Templates are listed in ProgramController.Show and consumed by SessionController.Create to pre-populate a session's
exercises. `exerciseLibraryJSON` (defined in exercise.go) is called here for autocomplete on New/Edit forms.
