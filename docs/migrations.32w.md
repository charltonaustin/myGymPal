---
level: 32w
parent: README.16w.md
deeper: data-model.md
relates-to:
  - models.32w.md
source: migrations/
---

34 numbered `golang-migrate` SQL files (000001–000034) build the schema incrementally: users, programs, phases,
templates, exercises, sessions, session exercises, cardio logs, body weights, macros, and macro goals. Migrations run
automatically on startup. See [data-model.md](data-model.md) for the current schema.

Recent migrations: `000032` split exercises into a global library, `000033` created program workout templates, and
`000034` added `session_exercises.linked_to_next` (BOOLEAN NOT NULL DEFAULT FALSE) for superset chaining.
