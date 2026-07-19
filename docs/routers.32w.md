---
level: 32w
parent: README.16w.md
deeper: routers/router.128w.md
relates-to:
  - controllers.32w.md
  - models.32w.md
source: routers/
---

| File      | Summary                                                                                                                                                                      | Detail                         |
|-----------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|--------------------------------|
| router.go | Registers all HTTP routes (including POST /sessions/:id/exercises/:eid/link for superset chaining), injects repository instances into controller globals, and registers custom template functions (`isDev`, `fmtDuration`, `add`, `restMinutes`, `restSecs`) | [128w](routers/router.128w.md) |
