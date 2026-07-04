---
level: 32w
parent: README.16w.md
deeper: data-model.md
relates-to:
  - models.32w.md
source: migrations/
---

31 numbered `golang-migrate` SQL files (000001–000031) build the schema incrementally: users, programs, phases,
templates, exercises, sessions, session exercises, cardio logs, body weights, macros, and macro goals. Migrations run
automatically on startup. See [data-model.md](data-model.md) for the current schema.
