---
level: 256w
parent: programs.32w.md
relates-to:
  - ../controllers/programs.128w.md
source: views/programs/new.tpl
---

## Purpose

Form to create a new training program. All phase-level settings (rep range, sets) set here become the default applied to
every phase; users can override per phase on the program detail page.

## Partials Included

- `partials/navbar.tpl`

## Template Variables

| Variable           | Type   | Description                                                  |
|--------------------|--------|--------------------------------------------------------------|
| `.Error`           | string | Server-side validation error; auto-dismissed after 3 000 ms. |
| `.Name`            | string | Pre-fills program name input.                                |
| `.StartDate`       | string | Pre-fills date input (YYYY-MM-DD).                           |
| `.NumPhases`       | int    | Pre-fills number of phases.                                  |
| `.WeeksPerPhase`   | int    | Pre-fills weeks per phase.                                   |
| `.WorkoutsPerWeek` | int    | Pre-fills workouts per week.                                 |
| `.DefaultRepMin`   | int    | Pre-fills minimum rep range.                                 |
| `.DefaultRepMax`   | int    | Pre-fills maximum rep range.                                 |
| `.DefaultSets`     | int    | Pre-fills default sets per exercise.                         |

## Form Fields

Form: `id="program-form" method="POST" action="/programs" novalidate data-offline-sync`

| Field             | Name                | Type   | Constraints     |
|-------------------|---------------------|--------|-----------------|
| Program Name      | `name`              | text   | required        |
| Start Date        | `start_date`        | date   | required        |
| Number of Phases  | `num_phases`        | number | min=1, required |
| Weeks per Phase   | `weeks_per_phase`   | number | min=1, required |
| Workouts per Week | `workouts_per_week` | number | min=1, required |
| Default Rep Min   | `default_rep_min`   | number | min=1, required |
| Default Rep Max   | `default_rep_max`   | number | min=1, required |
| Default Sets      | `default_sets`      | number | min=1, required |

## JavaScript Behavior

- Error alert auto-dismiss after 3 000 ms.
- Form submit: `checkValidity()` guard; adds `was-validated` class for Bootstrap error display.

## AJAX / Fetch

None.
