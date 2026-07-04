---
level: 128w
parent: ../controllers.32w.md
deeper: program.256w.md
relates-to:
  - session.128w.md
  - template.128w.md
source: controllers/program.go
---

The program controller manages training programs. GET /programs lists the user's programs (flash success if present).
GET /programs/new renders the creation form with sensible defaults. POST /programs creates a program after validating
all numeric fields; on success it flashes "Program created." and redirects to /programs. GET /programs/:id shows a
program with its phases, all templates, and all sessions belonging to it. POST /programs/:id updates the rep ranges,
default sets, and rest periods for all phases at once via `Phases.UpdateRepRanges`. POST /programs/:id/delete deletes
the program and redirects to /programs. All handlers are session-gated. The Show page aggregates data from Programs,
Phases, Templates, and Sessions repositories to render a complete program overview.
