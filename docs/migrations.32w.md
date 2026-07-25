---
level: 32w
parent: README.16w.md
deeper: data-model.md
relates-to:
  - models.32w.md
source: migrations/
---

36 numbered `golang-migrate` SQL files (000001–000036) build the schema incrementally: users, programs, phases,
templates, exercises, sessions, session exercises, cardio logs, body weights, macros, and macro goals. Migrations run
automatically on startup. See [data-model.md](data-model.md) for the current schema.

Recent migrations: `000032` split exercises into a global library, `000033` created program workout templates,
`000034` added `session_exercises.linked_to_next` (BOOLEAN NOT NULL DEFAULT FALSE) for superset chaining, and
`000035` created `template_circuits` and added `template_exercises.circuit_id` (nullable FK, `ON DELETE SET NULL`)
and `template_exercises.work_seconds` (INT NOT NULL DEFAULT 0), and `000036` did the same on the session side:
`session_circuits` plus `session_exercises.circuit_id` (nullable FK, `ON DELETE SET NULL`) and
`session_exercises.work_seconds` (INT NOT NULL DEFAULT 0).

`circuit_id` is nullable because NULL means "a normal, non-circuit exercise", which is every row that predates the
migration — so both apply to live data without a backfill. `ON DELETE SET NULL` rather than CASCADE: deleting a
circuit returns its exercises to the template or session instead of deleting logged work with it.

The session tables mirror the template ones rather than referencing them. A session is a record of what was
performed, so editing a template later must not rewrite a workout already done.
