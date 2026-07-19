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
indexed form fields (`exercise_name_0`, `is_bodyweight_0`, `block_0`, etc.). GET /templates/:id shows a template with
exercises grouped by block (main/abs/cardio/stretch). GET /templates/:id/edit and POST /templates/:id (Update) handle
editing. POST /templates/:id/delete removes a template. All handlers are session-gated. The New and Edit forms embed
`ExerciseLibraryJSON` (via the `exerciseLibraryJSON` helper from exercise.go) for autocomplete. Block grouping and
ordering logic (`blockOrder`, `blockLabels`) is defined in this file and shared with the session controller; block-name
sanitizing is `models.ValidBlock`. The New/Create and Edit/Update render paths set the `c.Data` chrome keys
(`Heading`, `FormAction`, `SubmitLabel`, `BackURL`, `BackLabel`) consumed by `partials/template_form.tpl`, via
`newFormChrome()` / `editFormChrome()`.
