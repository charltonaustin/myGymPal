---
level: 128w
parent: programs.32w.md
deeper: programs-new.256w.md
relates-to:
  - ../controllers/programs.128w.md
source: views/programs/new.tpl
---

Renders the New Training Program form. Includes `partials/navbar.tpl`. Contains fields for program name, start date,
number of phases, weeks per phase, workouts per week, default rep range (min–max), and default sets per exercise. Form
has `data-offline-sync` attribute. Bootstrap client-side validation (`was-validated`) is applied on submit.

Template variables: `.Error` (string, auto-dismissed after 3 s), `.Name`, `.StartDate`, `.NumPhases`, `.WeeksPerPhase`,
`.WorkoutsPerWeek`, `.DefaultRepMin`, `.DefaultRepMax`, `.DefaultSets` (all repopulate fields on re-render). POSTs to
`/programs`.
