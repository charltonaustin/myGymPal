---
level: 128w
parent: ../controllers.32w.md
deeper: auth.256w.md
relates-to:
  - account.128w.md
  - dashboard.128w.md
source: controllers/auth.go
---

The auth controller handles user authentication. GET /register and GET /login display their respective forms. POST
/login trims the username, looks it up via `Users.GetByUsername`, bcrypt-compares the password via `user.CheckPassword`,
and on success writes `user_id` and `username` to the session before redirecting to /dashboard. On failure it re-renders
`auth/login.tpl` with an inline `Error` and the submitted `Username`. POST /register validates that username and
password are non-empty, that the password is at least 8 characters, that the two password fields match, then calls
`Users.Create`. Duplicate-username errors set a specific message; other errors show a generic message. On success it
redirects to /login. GET /logout destroys the session and redirects to /login. No flash messages are used; errors are
passed directly via `c.Data["Error"]`.
