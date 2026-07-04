---
level: 128w
parent: programs.32w.md
deeper: programs-index.256w.md
relates-to:
  - ../controllers/programs.128w.md
source: views/programs/index.tpl
---

Renders the Training Programs list. Includes `partials/navbar.tpl`. Lists each program (name, start date, phase count,
weeks-per-phase) as a link to `/programs/:id`. A trash icon per row opens a Bootstrap delete-confirmation modal; on
confirm the modal sets the hidden form's action to `/programs/:id/delete` and submits via POST.

Template variables: `.Programs` (slice of program structs with `.ID`, `.Name`, `.StartDate`, `.NumPhases`,
`.WeeksPerPhase`), `.Success` (flash string, auto-dismissed after 3 s).

No AJAX calls. Registers `/sw.js`.
