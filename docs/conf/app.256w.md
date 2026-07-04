---
level: 256w
parent: app.128w.md
relates-to:
  - ../../routers.32w.md
source: conf/app.conf, conf/app.test.conf, conf/app.prod.conf
---

# Configuration — Full Detail

## Key Reference

| Key                     | dev                | test               | prod               | Purpose                                                           |
|-------------------------|--------------------|--------------------|--------------------|-------------------------------------------------------------------|
| `appname`               | `myGymPal`         | `myGymPal`         | `myGymPal`         | Beego application name                                            |
| `httpaddr`              | `localhost`        | `localhost`        | `localhost`        | Bind address (Caddy fronts in prod)                               |
| `httpport`              | `8080`             | `8080`             | `8080`             | HTTP listen port                                                  |
| `runmode`               | `dev`              | `dev`              | `prod`             | Beego run mode; `prod` disables template reload                   |
| `env`                   | `dev`              | _(absent)_         | `prod`             | Custom key read by `isDev` template function                      |
| `sessionon`             | `true`             | `true`             | `true`             | Enable Beego session middleware                                   |
| `sessionname`           | `mygympal_session` | `mygympal_session` | `mygympal_session` | Session cookie name                                               |
| `sessiongcmaxlifetime`  | `2592000`          | `86400`            | `2592000`          | Session GC lifetime in seconds (30 days / 1 day)                  |
| `sessioncookielifetime` | `2592000`          | _(absent)_         | `2592000`          | Cookie `Max-Age` in seconds; absent means session cookie in tests |
| `sessionprovider`       | `postgresql`       | `memory`           | `postgresql`       | Session storage backend                                           |
| `sessionproviderconfig` | DSN string         | _(empty)_          | DSN string         | Connection string for PostgreSQL session provider                 |
| `db_driver`             | `postgres`         | `postgres`         | `postgres`         | ORM driver (passed to golang-migrate and Beego ORM)               |
| `db_host`               | `localhost`        | `localhost`        | `localhost`        | PostgreSQL host                                                   |
| `db_port`               | `5432`             | `5432`             | `5432`             | PostgreSQL port                                                   |
| `db_name`               | `mygympal`         | `mygympal_test`    | `mygympal`         | Database name                                                     |
| `db_user`               | `postgres`         | `postgres`         | `postgres`         | Database user                                                     |
| `db_password`           | `postgres`         | `postgres`         | `postgres`         | Database password                                                 |
| `db_sslmode`            | `disable`          | `disable`          | `disable`          | PostgreSQL SSL mode                                               |

## Notes

- `app.prod.conf` is listed in `.gitignore` and copied to the DigitalOcean server manually via the deploy script.
- The `sessionproviderconfig` DSN in dev/prod is a libpq keyword/value string, not a URL.
- `runmode=prod` causes Beego to cache compiled templates; `runmode=dev` reloads them on every request.
- The `env` key (not a Beego built-in) is read exclusively by the `isDev` template function registered in
  `routers/router.go`.
