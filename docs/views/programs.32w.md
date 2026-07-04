---
level: 32w
parent: ../views.32w.md
deeper: programs-index.128w.md
relates-to:
  - ../controllers/programs.128w.md
source: views/programs/
---

| Template | Summary                                                                                                                                                                  | Detail                         |
|----------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------------------|
| index    | Lists all programs; delete-confirmation modal POSTs to `/programs/:id/delete`; success flash auto-dismisses                                                              | [128w](programs-index.128w.md) |
| new      | Form to create a training program: name, start date, phases, weeks/phase, workouts/week, default rep range, default sets; POSTs to `/programs`                           | [128w](programs-new.128w.md)   |
| show     | Program detail: workout history list with delete buttons, phase-settings editor (reps/sets/rest per phase with "copy to all"), "Start Workout" and "Log Workout" buttons | [128w](programs-show.128w.md)  |
