---
level: 128w
parent: exercises.32w.md
deeper: exercises-edit.256w.md
relates-to:
  - ../controllers/exercises.128w.md
source: views/exercises/edit.tpl
---

Edit-exercise form page. Structurally identical to `exercises/new.tpl` except the form POSTs to `/exercises/:id/edit`
and the submit button reads "Save Changes". Includes `partials/navbar.tpl` and `partials/exercise_fields.tpl`. Error
alert auto-removes after 4 s. Bootstrap `was-validated` applied on submit.

Template variables: `.Exercise.ID` (used in form action), `.Error`, and all exercise-fields partial variables
pre-populated with current values.
