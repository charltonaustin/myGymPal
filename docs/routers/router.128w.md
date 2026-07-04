---
level: 128w
parent: ../routers.32w.md
deeper: router.256w.md
relates-to:
  - ../controllers.32w.md
  - ../models.32w.md
source: routers/router.go
---

# router.go

`Register()` is the application entry point called before `beego.Run()`. It instantiates all ten ORM-backed repository
implementations and assigns them to the package-level globals in `controllers/repos.go`, injecting dependencies without
a DI framework.

Routes are organized into seven groups: PWA assets (`/sw.js`, `/manifest.json`, `/offline`), root/example, auth (
`/register`, `/login`, `/logout`), app pages (`/dashboard`, `/settings`, `/account/delete`, `/error`), programs,
sessions, macros, weight, exercises, and templates. URL parameters use Beego's `:id` / `:eid` / `:sid` / `:lid` syntax.

The `init()` function registers four custom template functions: `isDev` (reads the `env` config key), `fmtDuration` (
formats seconds as `m:ss` or `h:mm:ss`), `restMinutes`, and `restSecs`.
