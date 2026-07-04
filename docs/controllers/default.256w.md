---
level: 256w
parent: ../controllers.32w.md
relates-to:
  - auth.256w.md
  - dashboard.256w.md
  - pwa.256w.md
source: controllers/default.go
---

## Routes

| Method | Path     | Handler                    |
|--------|----------|----------------------------|
| GET    | /        | Get (default Beego action) |
| GET    | /example | Example                    |

## Auth requirement

None. Both handlers read `c.GetSession("user_id")` only to determine whether the user is logged in for nav rendering — a
nil session does not cause a redirect.

## Session keys

- Read: `user_id` (any; nil check only, value not used)

## Template variables

**GET /:**

- `c.Data["LoggedIn"]` = `c.GetSession("user_id") != nil`
- `c.Data["ActivePage"]` = "home"

**GET /example:**

- `c.Data["LoggedIn"]` = `c.GetSession("user_id") != nil`
- `c.Data["ActivePage"]` = ""

## Templates

- `index.tpl`
- `example.tpl`

## Repository calls

None.

## Flash messages

None.

## Redirect paths

None. Both handlers always render their template.

## Relationship to other controllers

MainController is the unauthenticated entry point. The navigation rendered by `index.tpl` and `example.tpl` links to
/login and /register (AuthController) when `LoggedIn` is false, or to /dashboard and other authenticated routes when
`LoggedIn` is true. The Beego router maps GET / to `MainController.Get` via the default handler convention.
