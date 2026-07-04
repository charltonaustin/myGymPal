---
level: 256w
parent: ../controllers.32w.md
relates-to:
  - auth.256w.md
source: controllers/account.go
---

## Routes

| Method | Path            | Handler       |
|--------|-----------------|---------------|
| GET    | /settings       | Settings      |
| POST   | /settings       | SettingsPost  |
| POST   | /account/delete | DeleteAccount |

## Auth requirement

All three handlers check `c.GetSession("user_id")`; nil redirects to /login.

## Session keys

- Read: `user_id` (int64)
- Destroyed on account deletion: `c.DestroySession()`

## Template variables

**Settings (GET):**

- `c.Data["LoggedIn"]` = true
- `c.Data["ActivePage"]` = "settings"
- `c.Data["WeightUnit"]` = user.WeightUnit
- `c.Data["Success"]` = flash success string (if present)

**Settings (POST failure):**

- `c.Data["LoggedIn"]` = true
- `c.Data["ActivePage"]` = "settings"
- `c.Data["Error"]` = "Invalid weight unit."
- `c.Data["WeightUnit"]` = submitted unit string

**DeleteAccount (failure):**

- `c.Data["LoggedIn"]` = true
- `c.Data["ActivePage"]` = "settings"
- `c.Data["Error"]` = "Failed to delete account. Please try again."
- `c.Data["WeightUnit"]` = reloaded from DB (if available)

## Template

- `settings.tpl`

## Repository calls

- `Users.GetByID(userID)` — Settings GET, DeleteAccount failure fallback
- `Users.UpdateWeightUnit(userID, unit)` — SettingsPost
- `Users.DeleteByID(userID)` — DeleteAccount

## Flash messages

- `flash.Success("Settings saved.")` — SettingsPost success

## Redirect paths

- SettingsPost success → /settings
- SettingsPost failure → re-renders settings.tpl
- DeleteAccount success → /login
- DeleteAccount failure → re-renders settings.tpl

## Relationship to other controllers

AuthController writes the session that AccountController reads. Weight unit set here propagates to WeightController,
ExerciseController, SessionController, and DashboardController, which all call `Users.GetByID` to read the preferred
unit.
