---
level: 256w
parent: ../controllers.32w.md
relates-to:
  - session.256w.md
  - template.256w.md
source: controllers/program.go
---

## Routes

| Method | Path                 | Handler      |
|--------|----------------------|--------------|
| GET    | /programs            | Index        |
| GET    | /programs/new        | New          |
| POST   | /programs            | Create       |
| GET    | /programs/:id        | Show         |
| POST   | /programs/:id        | UpdatePhases |
| POST   | /programs/:id/delete | Delete       |

## Auth requirement

All handlers check `c.GetSession("user_id")`; nil redirects to /login.

## Session keys

- Read: `user_id` (int64)

## Template variables — Index

- `c.Data["LoggedIn"]` = true
- `c.Data["ActivePage"]` = "programs"
- `c.Data["Programs"]` = `[]*models.Program`
- `c.Data["Success"]` = flash success (if present)

## Template variables — New / Create failure

- `c.Data["WeeksPerPhase"]` = "8" (default)
- `c.Data["WorkoutsPerWeek"]` = "4" (default)
- `c.Data["DefaultRepMin"]` = "12" (default)
- `c.Data["DefaultRepMax"]` = "14" (default)
- `c.Data["DefaultSets"]` = "3" (default)
- `c.Data["Error"]`, `c.Data["Name"]`, `c.Data["StartDate"]`, `c.Data["NumPhases"]`, etc. (on failure)

## Template variables — Show

- `c.Data["Program"]` = `*models.Program`
- `c.Data["Phases"]` = `[]*models.Phase`
- `c.Data["Templates"]` = `[]*models.Template`
- `c.Data["Sessions"]` = `[]*models.Session`
- `c.Data["Success"]` = flash success (if present)
- `c.Data["Error"]` = inline error (UpdatePhases failure)

## Templates

- `programs/index.tpl`
- `programs/new.tpl`
- `programs/show.tpl`

## Repository calls

- `Programs.GetAllByUser(userID)` — Index
- `Programs.Create(userID, name, startDate, numPhases, weeksPerPhase, workoutsPerWeek, repMin, repMax, sets)` — Create
- `Programs.GetByID(id, userID)` — Show, UpdatePhases, Delete
- `Programs.Delete(id, userID)` — Delete
- `Phases.GetByProgram(id)` — Show, UpdatePhases
- `Phases.UpdateRepRanges(id, []PhaseUpdate)` — UpdatePhases
- `Templates.GetAll()` — Show
- `Sessions.GetByProgram(id)` — Show

## Flash messages

- `flash.Success("Program created.")` — Create success
- `flash.Success("Program deleted.")` — Delete success
- `flash.Success("Phase settings saved.")` — UpdatePhases success

## Redirect paths

- Create success → /programs
- Create failure → re-renders programs/new.tpl
- Delete success → /programs
- UpdatePhases success → /programs/:id
- UpdatePhases failure → re-renders programs/show.tpl with error

## Relationship to other controllers

Show page lists sessions managed by SessionController and templates managed by TemplateController. SessionController.New
links back to /programs/:id on error.
