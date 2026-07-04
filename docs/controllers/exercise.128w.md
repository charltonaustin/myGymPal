---
level: 128w
parent: ../controllers.32w.md
deeper: exercise.256w.md
relates-to:
  - session.128w.md
  - template.128w.md
source: controllers/exercise.go
---

The exercise controller manages the user's personal exercise library. GET /exercises lists all exercises with goal
weights converted to the user's preferred unit and reads a flash success. GET /exercises/new and POST /exercises/new (
Create) display and submit the new-exercise form. GET /exercises/:id/edit and POST /exercises/:id/edit (Update) display
and submit the edit form. POST /exercises/:id/delete removes an exercise. All five CRUD handlers are session-gated and
redirect unauthenticated requests to /login.

Three AJAX POST endpoints allow in-session goal updates without a page reload: POST /exercises/goal-weight updates goal
weight and unit, POST /exercises/goal-reps updates the bodyweight rep range, and POST /exercises/goal-seconds updates
the time-based goal duration. These serve JSON responses. The package-level helper `exerciseLibraryJSON` serializes the
full exercise library as `template.JS` for safe embedding in `<script>` tags, used by SessionController and
TemplateController.
