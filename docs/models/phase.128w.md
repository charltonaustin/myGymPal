---
level: 128w
parent: ../models.32w.md
deeper: phase.256w.md
relates-to:
  - program.128w.md
  - session.128w.md
  - ../controllers/program.128w.md
source: models/phase.go, models/phase_repository.go
---

# Phase

`Phase` represents a numbered training block within a `Program`. Each phase stores `PhaseNumber`, `RepMin`, `RepMax`,
`DefaultSets`, and `RestSeconds`. Phases have no timestamps; they are created once (during program creation) and updated
in bulk.

The `PhaseUpdate` helper struct carries all mutable fields for a batch update: `PhaseNumber`, `RepMin`, `RepMax`,
`DefaultSets`, and `RestSeconds`. `UpdatePhaseRepRanges` validates every update (rep_min > 0, rep_max >= rep_min,
default_sets > 0) before issuing any SQL, issuing one `QueryTable.Filter.Update` per phase number.

The `PhaseRepository` interface exposes only two methods: `GetByProgram` (ordered by `PhaseNumber`) and
`UpdateRepRanges`. Read is always scoped to a program ID; there is no single-phase lookup by ID. The program controller
uses these to display and edit phase configuration.
