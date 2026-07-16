---
level: 32w
parent: README.16w.md
deeper: controllers/auth.128w.md
relates-to:
  - models.32w.md
  - routers.32w.md
source: controllers/
---

| Controller | Summary                                                                                                                                                                        | Detail                                                                        |
|------------|--------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|-------------------------------------------------------------------------------|
| auth       | Handles GET/POST /login and /register for user authentication; writes user_id and username to session on login; renders error inline; redirects to /dashboard on success.      | [128w](controllers/auth.128w.md) · [256w](controllers/auth.256w.md)           |
| account    | Manages GET/POST /settings for weight-unit preferences and POST /account/delete for account removal; requires session; flashes success on save.                                | [128w](controllers/account.128w.md) · [256w](controllers/account.256w.md)     |
| dashboard  | Renders GET /dashboard with recent sessions, 3-day body-weight average, and 3-day macro summary; requires session; aggregates data from five repositories.                     | [128w](controllers/dashboard.128w.md) · [256w](controllers/dashboard.256w.md) |
| default    | Serves GET / (home) and GET /example; sets LoggedIn from session presence; no auth requirement; renders index.tpl and example.tpl.                                             | [128w](controllers/default.128w.md) · [256w](controllers/default.256w.md)     |
| error      | Minimal GET /error handler; renders error.tpl with no data and no auth check.                                                                                                  | [128w](controllers/error.128w.md) · [256w](controllers/error.256w.md)         |
| exercise   | Full CRUD for the user's exercise library at /exercises; three AJAX POST endpoints for live goal-weight, goal-reps, and goal-seconds updates from within a session. GET /exercises/history renders a GitHub-style activity heatmap of the last year (days exercised, shaded by that day's set count) plus a Chart.js line chart where users can compare multiple exercises over time, pre-populated by default with exercises performed in the last two weeks (14 days), with a range selector (2 weeks / 1 / 3 / 6 months / 1 year / all) that clamps the plotted window; GET /exercises/history/data returns JSON series data (max weight or max reps per session per exercise), accepting a `days` query param to limit the window. | [128w](controllers/exercise.128w.md) · [256w](controllers/exercise.256w.md)   |
| macro      | Full CRUD for daily food/macro entries at /macros; POST /macros/goals to upsert daily macro goals; groups entries by day and computes 3-day summary.                           | [128w](controllers/macro.128w.md) · [256w](controllers/macro.256w.md)         |
| program    | Manages training programs at /programs; shows phases; POST /programs/:id updates phase rep ranges and rest periods via PhaseRepository.                                        | [128w](controllers/program.128w.md) · [256w](controllers/program.256w.md)     |
| session    | Most complex controller; full workout-session lifecycle at /programs/:id/sessions and /sessions/:id; logs sets and cardio; auto-promotes goal weight after 3rd successful set. | [128w](controllers/session.128w.md) · [256w](controllers/session.256w.md)     |
| template   | CRUD for reusable workout templates at /templates; stores exercise name, bodyweight flag, time-based flag, and block; injects exercise library JSON for autocomplete.          | [128w](controllers/template.128w.md) · [256w](controllers/template.256w.md)   |
| weight     | CRUD for body-weight log entries at /weight; converts stored units to user preference at display time; computes 3-day rolling average shared with dashboard.                   | [128w](controllers/weight.128w.md) · [256w](controllers/weight.256w.md)       |
| pwa        | Serves PWA assets: GET /sw.js (service worker), GET /manifest.json (web app manifest), GET /offline (offline fallback page); no auth or session use.                           | [128w](controllers/pwa.128w.md) · [256w](controllers/pwa.256w.md)             |
