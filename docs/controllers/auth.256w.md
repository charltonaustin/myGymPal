---
level: 256w
parent: ../controllers.32w.md
relates-to:
  - account.256w.md
  - dashboard.256w.md
source: controllers/auth.go
---

## Routes

| Method | Path      | Handler      |
|--------|-----------|--------------|
| GET    | /register | Register     |
| POST   | /register | RegisterPost |
| GET    | /login    | Login        |
| POST   | /login    | LoginPost    |
| GET    | /logout   | Logout       |

## Session keys

- Read: none (GET handlers do not check for existing session)
- Written on login success: `user_id` (int64), `username` (string)
- Destroyed on logout: `DestroySession()`

## Template variables

**Login (GET + POST failure):** `c.Data["Error"]` (string), `c.Data["Username"]` (string)

**Register (GET + POST failure):** `c.Data["Error"]` (string), `c.Data["Username"]` (string)

## Templates

- `auth/login.tpl`
- `auth/register.tpl`

## Repository calls

- `Users.GetByUsername(username)` — POST /login
- `user.CheckPassword(password)` — method on returned User model
- `Users.Create(username, password, "lb")` — POST /register (default weight unit is `lb`)

## Redirect paths

- POST /login success → /dashboard
- POST /login failure → re-renders auth/login.tpl (no redirect)
- POST /register success → /login
- POST /register failure → re-renders auth/register.tpl (no redirect)
- GET /logout → /login
- Any session-destroy error in Logout → /error

## Validation (RegisterPost)

1. Username and password must be non-empty.
2. Password must be at least 8 characters.
3. Password must match confirm_password.
4. Unique/duplicate error from DB maps to "That username is already taken."

## Flash messages

None — errors are passed inline via `c.Data["Error"]`.

## Relationship to other controllers

Once session is established by LoginPost, every other authenticated controller reads `c.GetSession("user_id")` and
redirects to /login if nil. The DashboardController is the immediate destination after a successful login.
