---
level: 128w
parent: ../conf.32w.md
deeper: app.256w.md
relates-to:
  - ../routers.32w.md
source: conf/
---

# Configuration Files

Three Beego INI config files control environment-specific behavior. All share the same key set; only values differ.

**app.conf** (development): connects to the `mygympal` PostgreSQL database on `localhost:5432` with `postgres/postgres`
credentials. Sessions are PostgreSQL-backed with a 30-day lifetime (`sessiongcmaxlifetime=2592000`). `runmode=dev`,
`env=dev`.

**app.test.conf** (tests): uses the `mygympal_test` database to isolate test data. Sessions use the `memory` provider
with a 1-day lifetime (86400 s) — no PostgreSQL session table required. `runmode=dev`, no `env` key (so `isDev` returns
false in test templates).

**app.prod.conf** (production): mirrors `app.conf` with `runmode=prod` and `env=prod`. Git-excluded and copied to the
server manually. Caddy reverse-proxies to `localhost:8080`.
