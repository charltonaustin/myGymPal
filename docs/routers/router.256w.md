---
level: 256w
parent: router.128w.md
relates-to:
  - ../../controllers.32w.md
  - ../../models.32w.md
source: routers/router.go
---

# router.go — Full Detail

## Repository Injection

`Register()` instantiates all repositories and assigns them to controller globals before routes are registered:
`Users`, `Programs`, `Phases`, `Templates`, `Sessions`, `SessionExercises`, `Exercises`, `BodyWeights`, `Macros`,
`MacroGoals`.

## Template Functions (`init`)

| Function      | Signature           | Behavior                                          |
|---------------|---------------------|---------------------------------------------------|
| `isDev`       | `() bool`           | Returns true when `env` config key equals `"dev"` |
| `fmtDuration` | `(secs int) string` | Formats as `m:ss`; includes hours when `h > 0`    |
| `restMinutes` | `(secs int) int`    | Integer division `secs / 60`                      |
| `restSecs`    | `(secs int) int`    | Remainder `secs % 60`                             |

## Route Table

| Method   | Path                                              | Controller:Action                        |
|----------|---------------------------------------------------|------------------------------------------|
| GET      | `/sw.js`                                          | PWAController:ServiceWorker              |
| GET      | `/manifest.json`                                  | PWAController:Manifest                   |
| GET      | `/offline`                                        | PWAController:Offline                    |
| GET/POST | `/`                                               | MainController (default)                 |
| GET      | `/example`                                        | MainController:Example                   |
| GET      | `/register`                                       | AuthController:Register                  |
| POST     | `/register`                                       | AuthController:RegisterPost              |
| GET      | `/login`                                          | AuthController:Login                     |
| POST     | `/login`                                          | AuthController:LoginPost                 |
| GET      | `/logout`                                         | AuthController:Logout                    |
| GET/POST | `/dashboard`                                      | DashboardController (default)            |
| GET      | `/settings`                                       | AccountController:Settings               |
| POST     | `/settings`                                       | AccountController:SettingsPost           |
| POST     | `/account/delete`                                 | AccountController:DeleteAccount          |
| GET      | `/error`                                          | ErrorController (default)                |
| GET      | `/programs`                                       | ProgramController:Index                  |
| POST     | `/programs`                                       | ProgramController:Create                 |
| GET      | `/programs/new`                                   | ProgramController:New                    |
| GET      | `/programs/:id`                                   | ProgramController:Show                   |
| POST     | `/programs/:id`                                   | ProgramController:UpdatePhases           |
| POST     | `/programs/:id/delete`                            | ProgramController:Delete                 |
| GET      | `/programs/:id/sessions/new`                      | SessionController:New                    |
| POST     | `/programs/:id/sessions`                          | SessionController:Create                 |
| GET      | `/sessions/:id`                                   | SessionController:Show                   |
| POST     | `/sessions/:id/delete`                            | SessionController:Delete                 |
| POST     | `/sessions/:id/exercises`                         | SessionController:AddExercise            |
| POST     | `/sessions/:id/exercises/reorder`                 | SessionController:ReorderExercises       |
| POST     | `/sessions/:id/exercises/:eid/delete`             | SessionController:DeleteExercise         |
| POST     | `/sessions/:id/exercises/:eid/sets`               | SessionController:LogSet                 |
| POST     | `/sessions/:id/exercises/:eid/sets/:sid/delete`   | SessionController:DeleteSet              |
| POST     | `/sessions/:id/cardio`                            | SessionController:AddCardioActivity      |
| POST     | `/sessions/:id/exercises/:eid/cardio`             | SessionController:LogCardio              |
| POST     | `/sessions/:id/exercises/:eid/cardio/:lid/delete` | SessionController:DeleteCardioLog        |
| GET      | `/macros`                                         | MacroController:Index                    |
| POST     | `/macros`                                         | MacroController:Create                   |
| POST     | `/macros/goals`                                   | MacroController:SaveGoal                 |
| POST     | `/macros/:id`                                     | MacroController:Update                   |
| POST     | `/macros/:id/delete`                              | MacroController:Delete                   |
| GET      | `/weight`                                         | WeightController:Index                   |
| POST     | `/weight`                                         | WeightController:Create                  |
| POST     | `/weight/:id`                                     | WeightController:Update                  |
| POST     | `/weight/:id/delete`                              | WeightController:Delete                  |
| GET      | `/exercises`                                      | ExerciseController:Index                 |
| GET      | `/exercises/new`                                  | ExerciseController:New                   |
| POST     | `/exercises/new`                                  | ExerciseController:Create                |
| POST     | `/exercises/goal-weight`                          | ExerciseController:UpdateGoalWeightJSON  |
| POST     | `/exercises/goal-reps`                            | ExerciseController:UpdateGoalRepsJSON    |
| POST     | `/exercises/goal-seconds`                         | ExerciseController:UpdateGoalSecondsJSON |
| GET      | `/exercises/:id/edit`                             | ExerciseController:Edit                  |
| POST     | `/exercises/:id/edit`                             | ExerciseController:Update                |
| POST     | `/exercises/:id/delete`                           | ExerciseController:Delete                |
| GET      | `/templates`                                      | TemplateController:Index                 |
| GET      | `/templates/new`                                  | TemplateController:New                   |
| POST     | `/templates/new`                                  | TemplateController:Create                |
| GET      | `/templates/:id`                                  | TemplateController:Show                  |
| POST     | `/templates/:id`                                  | TemplateController:Update                |
| GET      | `/templates/:id/edit`                             | TemplateController:Edit                  |
| POST     | `/templates/:id/delete`                           | TemplateController:Delete                |
