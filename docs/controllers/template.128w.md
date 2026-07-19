---
level: 128w
parent: ../controllers.32w.md
deeper: template.256w.md
relates-to:
  - session.128w.md
  - exercise.128w.md
source: controllers/template.go
---

The template controller manages reusable workout templates (not HTML templates). GET /templates lists all templates. GET
/templates/new and POST /templates/new (Create) display and submit the creation form; exercises are collected from
indexed form fields (`exercise_name_0`, `is_bodyweight_0`, `block_0`, `work_seconds_0`, `circuit_index_0`, etc.), and
circuits from `circuit_count` plus `circuit_name_0` / `circuit_rounds_0` / `circuit_transition_0`. GET /templates/:id
shows a template with loose exercises grouped by block (main/abs/cardio/stretch) and each circuit rendered as its own
card; circuit members are excluded from the block grouping so they are not listed twice. GET /templates/:id/edit and
POST /templates/:id (Update) handle editing. POST /templates/:id/delete removes a template. All handlers are
session-gated. The New and Edit forms embed `ExerciseLibraryJSON` (via the `exerciseLibraryJSON` helper from
exercise.go) for autocomplete. Block grouping and ordering logic (`blockOrder`, `blockLabels`) is defined in this file
and shared with the session controller; block-name sanitizing is `models.ValidBlock`, and rounds/seconds are clamped
with `models.ValidRounds` / `models.ValidSeconds` — the form is never trusted.

`parseTemplateForm` renumbers circuits contiguously as it reads them, dropping any whose name the user cleared, and
maps each exercise's submitted `circuit_index` through that compaction. Without it, deleting a circuit would slide
every exercise below it onto the wrong one. An exercise whose circuit is gone becomes loose rather than a dangling
reference.

## The `c.Data` contract for `partials/template_form.tpl`

There are four render paths into the form — `New`, `Create`'s error re-render, `Edit`, and `Update`'s error re-render
— and a `c.Data` key missed on one of them renders as an empty string, not an error. The keys are therefore set in two
places and nowhere else:

- **Page chrome** (differs per *page*): `newFormChrome()` / `editFormChrome(tmpl)` set `Heading`, `FormAction`,
  `SubmitLabel`, `BackURL`, `BackLabel`.
- **Form body** (differs per *render path*): each path builds a `templateFormData` and passes it to `setFormBody`,
  which sets `Name`, `Focus`, `Exercises` (loose only), `Circuits`, `ExerciseCount` and `CircuitCount`. Because the
  body is a struct, a forgotten field is a compile error rather than a blank field on the page.

New keys the form needs go in one of those two places. Adding one to a render path directly re-creates the hazard both
exist to remove.
