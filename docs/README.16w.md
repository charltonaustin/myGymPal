---
level: 16w
parent: ../CLAUDE.md
deeper: controllers.32w.md
relates-to: [ ]
source: ./
---

Root index for myGymPal documentation. Each row is one section; follow the link to get 32-word summaries of everything
inside.

| Section      | 16-word summary                                                                        | Entry file                               |
|--------------|----------------------------------------------------------------------------------------|------------------------------------------|
| controllers  | HTTP request handlers for auth, workouts, exercises, programs, macros, weight, and PWA | [controllers.32w.md](controllers.32w.md) |
| models       | Domain structs and repository interfaces for every database entity in the application  | [models.32w.md](models.32w.md)           |
| views        | Beego HTML templates rendered by controllers, organized by feature area and partials   | [views.32w.md](views.32w.md)             |
| routers      | Route registration, repository injection, and custom template function setup           | [routers.32w.md](routers.32w.md)         |
| migrations   | PostgreSQL schema evolution via golang-migrate; 31 migrations from users to exercises  | [migrations.32w.md](migrations.32w.md)   |
| conf         | Application configuration for development, test, and production environments           | [conf.32w.md](conf.32w.md)               |
| static       | PWA manifest, service worker, offline sync handler, and app icons                      | [static.32w.md](static.32w.md)           |
| data-model   | Full PostgreSQL schema with all tables, indexes, and design decisions documented       | [data-model.md](data-model.md)           |
| user-stories | User story specifications indexed by US number, covering all implemented features      | [user-stories.md](user-stories.md)       |
