# myGymPal — PostgreSQL Data Model

## Overview

This document describes the complete relational schema for myGymPal. The stack is Go (Beego v2) with `lib/pq` connecting to PostgreSQL. All tables use `BIGSERIAL` surrogate primary keys. Timestamps are stored as `TIMESTAMPTZ`. Weight values are stored in the user's preferred unit (lb or kg) as entered — no normalisation to a canonical unit is applied at the database layer.

---

## Tables

### `users`

```sql
CREATE TABLE users (
    id              BIGSERIAL       PRIMARY KEY,
    username        TEXT            NOT NULL UNIQUE,
    password_hash   TEXT            NOT NULL,
    weight_unit     TEXT            NOT NULL DEFAULT 'lb'
                        CHECK (weight_unit IN ('lb', 'kg')),
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);
```

---

### `exercises`

Global catalogue of exercises. Shared across all users and referenced by templates.

```sql
CREATE TABLE exercises (
    id              BIGSERIAL       PRIMARY KEY,
    name            TEXT            NOT NULL UNIQUE,
    exercise_type   TEXT            NOT NULL
                        CHECK (exercise_type IN ('weighted', 'bodyweight')),
    block_type      TEXT            NOT NULL
                        CHECK (block_type IN ('main', 'abs', 'cardio', 'stretch')),
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);
```

`block_type` is stored on the exercise to prevent a main-block exercise from being placed in the abs block. Template editors can only select exercises whose `block_type` matches the block being configured.

---

### `templates`

Reusable workout blueprints, shared across all users.

```sql
CREATE TABLE templates (
    id              BIGSERIAL       PRIMARY KEY,
    created_by      BIGINT          REFERENCES users(id) ON DELETE SET NULL,
    name            TEXT            NOT NULL,
    focus           TEXT,
    is_public       BOOLEAN         NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);
```

`created_by` is SET NULL if the creator deletes their account so templates persist for other users.

---

### `template_main_exercises`

Ordered main-block exercises within a template.

```sql
CREATE TABLE template_main_exercises (
    id                  BIGSERIAL       PRIMARY KEY,
    template_id         BIGINT          NOT NULL REFERENCES templates(id) ON DELETE CASCADE,
    exercise_id         BIGINT          NOT NULL REFERENCES exercises(id),
    position            INT             NOT NULL,
    sets                INT             NOT NULL DEFAULT 3,
    weight_increment    NUMERIC(6,2)    NOT NULL DEFAULT 5.00,
    rep_increment       INT             NOT NULL DEFAULT 1,
    goal_weight         NUMERIC(6,2),   -- NULL for bodyweight exercises
    goal_reps           INT,            -- NULL for weighted (phase rep range is used instead)
    UNIQUE (template_id, position)
);
```

For **weighted** exercises: `goal_weight` is set, `goal_reps` is NULL (the phase's `rep_min`/`rep_max` drives reps).
For **bodyweight** exercises: `goal_reps` is set, `goal_weight` is NULL.

---

### `template_abs_exercises`

Abs-block exercises within a template. Goal reps are fixed by the template, not driven by the phase rep range.

```sql
CREATE TABLE template_abs_exercises (
    id                  BIGSERIAL       PRIMARY KEY,
    template_id         BIGINT          NOT NULL REFERENCES templates(id) ON DELETE CASCADE,
    exercise_id         BIGINT          NOT NULL REFERENCES exercises(id),
    position            INT             NOT NULL,
    sets                INT             NOT NULL DEFAULT 3,
    weight_increment    NUMERIC(6,2)    NOT NULL DEFAULT 5.00,
    rep_increment       INT             NOT NULL DEFAULT 1,
    goal_weight         NUMERIC(6,2),   -- NULL for bodyweight abs
    goal_reps           INT             NOT NULL, -- always present; never derived from phase range
    UNIQUE (template_id, position)
);
```

---

### `template_cardio_blocks`

```sql
CREATE TABLE template_cardio_blocks (
    id                      BIGSERIAL   PRIMARY KEY,
    template_id             BIGINT      NOT NULL REFERENCES templates(id) ON DELETE CASCADE,
    position                INT         NOT NULL,
    cardio_type             TEXT        NOT NULL,
    goal_duration_minutes   INT         NOT NULL,
    UNIQUE (template_id, position)
);
```

---

### `template_stretch_blocks`

```sql
CREATE TABLE template_stretch_blocks (
    id                  BIGSERIAL   PRIMARY KEY,
    template_id         BIGINT      NOT NULL REFERENCES templates(id) ON DELETE CASCADE,
    exercise_id         BIGINT      NOT NULL REFERENCES exercises(id),
    position            INT         NOT NULL,
    goal_hold_seconds   INT         NOT NULL,
    UNIQUE (template_id, position)
);
```

---

### `programs`

A user's training program. One user can have multiple programs.

```sql
CREATE TABLE programs (
    id              BIGSERIAL       PRIMARY KEY,
    user_id         BIGINT          NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name            TEXT            NOT NULL,
    start_date      DATE            NOT NULL,
    num_phases      INT             NOT NULL DEFAULT 6,
    weeks_per_phase INT             NOT NULL DEFAULT 8,
    is_active       BOOLEAN         NOT NULL DEFAULT TRUE,
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);
```

The last week of each phase (`week_number = weeks_per_phase`) is always the deload week. This is a computation, not stored on the phase.

---

### `phases`

One row per phase within a program. Each phase has its own rep range.

```sql
CREATE TABLE phases (
    id              BIGSERIAL   PRIMARY KEY,
    program_id      BIGINT      NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
    phase_number    INT         NOT NULL,
    rep_min         INT         NOT NULL,
    rep_max         INT         NOT NULL,
    UNIQUE (program_id, phase_number),
    CONSTRAINT chk_rep_range CHECK (rep_min <= rep_max AND rep_min > 0)
);
```

---

### `workouts`

A single workout instance, linked to a program, phase, week, and the template it was built from.

```sql
CREATE TABLE workouts (
    id              BIGSERIAL       PRIMARY KEY,
    user_id         BIGINT          NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    program_id      BIGINT          NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
    phase_id        BIGINT          NOT NULL REFERENCES phases(id),
    template_id     BIGINT          REFERENCES templates(id) ON DELETE SET NULL,
    week_number     INT             NOT NULL,
    workout_date    DATE            NOT NULL DEFAULT CURRENT_DATE,
    workout_type    TEXT            NOT NULL DEFAULT 'normal'
                        CHECK (workout_type IN ('normal', 'deload')),
    status          TEXT            NOT NULL DEFAULT 'in_progress'
                        CHECK (status IN ('in_progress', 'completed')),
    created_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ     NOT NULL DEFAULT NOW(),
    completed_at    TIMESTAMPTZ
);
```

`workout_type` is stored explicitly (not computed from `week_number`) so that:
1. Progression queries use a simple `WHERE workout_type = 'normal'` without joining back to `programs`.
2. Historical workout types are preserved even if `weeks_per_phase` is later edited.
3. New types (e.g. `'test'`, `'competition'`) can be added without a schema change.

---

### `program_exercise_state`

Live goal weight or goal reps for each exercise within a program. Evolves via progression events and carries across phases unchanged.

```sql
CREATE TABLE program_exercise_state (
    id                          BIGSERIAL       PRIMARY KEY,
    program_id                  BIGINT          NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
    exercise_id                 BIGINT          NOT NULL REFERENCES exercises(id),
    block_type                  TEXT            NOT NULL CHECK (block_type IN ('main', 'abs')),
    goal_weight                 NUMERIC(6,2),   -- NULL for bodyweight
    goal_reps                   INT,            -- NULL for weighted main exercises
    consecutive_max_workouts    INT             NOT NULL DEFAULT 0,
    weight_increment            NUMERIC(6,2)    NOT NULL DEFAULT 5.00,
    rep_increment               INT             NOT NULL DEFAULT 1,
    UNIQUE (program_id, exercise_id, block_type)
);
```

`consecutive_max_workouts` is the progression state machine. On workout completion:
- `workout_type != 'normal'` → skip all updates.
- All sets hit `phase.rep_max` → increment; if reaches 3, apply increase and reset to 0.
- Any set missed `phase.rep_min` → apply decrease and reset to 0.
- Between min and max → reset to 0.

---

### `progression_events`

Immutable audit log of every auto-progression change. Powers the post-workout indicator and trend history.

```sql
CREATE TABLE progression_events (
    id                  BIGSERIAL       PRIMARY KEY,
    workout_id          BIGINT          NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
    exercise_id         BIGINT          NOT NULL REFERENCES exercises(id),
    block_type          TEXT            NOT NULL CHECK (block_type IN ('main', 'abs')),
    event_type          TEXT            NOT NULL CHECK (event_type IN ('increase', 'decrease')),
    old_goal_weight     NUMERIC(6,2),
    new_goal_weight     NUMERIC(6,2),
    old_goal_reps       INT,
    new_goal_reps       INT,
    created_at          TIMESTAMPTZ     NOT NULL DEFAULT NOW()
);
```

The event is recorded against the triggering workout. The new goal is applied when the **next workout is created** — the application reads pending events for each exercise when pre-filling the new workout.

---

### `workout_main_sets`

Logged sets for main-block exercises. Goals are snapshotted from `program_exercise_state` and `phases` at workout creation so historical targets are preserved even after progression fires.

```sql
CREATE TABLE workout_main_sets (
    id              BIGSERIAL       PRIMARY KEY,
    workout_id      BIGINT          NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
    exercise_id     BIGINT          NOT NULL REFERENCES exercises(id),
    set_number      INT             NOT NULL,
    goal_weight     NUMERIC(6,2),   -- NULL for bodyweight; snapshotted at workout creation
    goal_reps       INT             NOT NULL, -- phase rep_max for weighted; goal_reps for bodyweight; rep_min-2 for deload
    UNIQUE (workout_id, exercise_id, set_number)
);
```

---

### `workout_main_set_segments`

Sub-rows within a single set to capture mid-set weight drops. A normal set has one segment; a drop set has two or more.

```sql
CREATE TABLE workout_main_set_segments (
    id              BIGSERIAL       PRIMARY KEY,
    set_id          BIGINT          NOT NULL REFERENCES workout_main_sets(id) ON DELETE CASCADE,
    segment_number  INT             NOT NULL,
    actual_weight   NUMERIC(6,2),   -- NULL for bodyweight
    actual_reps     INT             NOT NULL,
    UNIQUE (set_id, segment_number)
);
```

Progression logic sums `actual_reps` across all segments of a set before comparing to `goal_reps`.

---

### `workout_abs_sets`

```sql
CREATE TABLE workout_abs_sets (
    id              BIGSERIAL       PRIMARY KEY,
    workout_id      BIGINT          NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
    exercise_id     BIGINT          NOT NULL REFERENCES exercises(id),
    set_number      INT             NOT NULL,
    goal_weight     NUMERIC(6,2),   -- NULL for bodyweight abs
    goal_reps       INT             NOT NULL,
    UNIQUE (workout_id, exercise_id, set_number)
);
```

### `workout_abs_set_segments`

```sql
CREATE TABLE workout_abs_set_segments (
    id              BIGSERIAL       PRIMARY KEY,
    set_id          BIGINT          NOT NULL REFERENCES workout_abs_sets(id) ON DELETE CASCADE,
    segment_number  INT             NOT NULL,
    actual_weight   NUMERIC(6,2),
    actual_reps     INT             NOT NULL,
    UNIQUE (set_id, segment_number)
);
```

---

### `workout_cardio_logs`

```sql
CREATE TABLE workout_cardio_logs (
    id                      BIGSERIAL   PRIMARY KEY,
    workout_id              BIGINT      NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
    position                INT         NOT NULL,
    cardio_type             TEXT        NOT NULL,
    goal_duration_minutes   INT         NOT NULL,
    actual_duration_minutes INT         NOT NULL,
    notes                   TEXT
);
```

`actual_duration_minutes` may be less than `goal_duration_minutes` — this is valid and never treated as an error.

---

### `workout_stretch_logs`

```sql
CREATE TABLE workout_stretch_logs (
    id                  BIGSERIAL   PRIMARY KEY,
    workout_id          BIGINT      NOT NULL REFERENCES workouts(id) ON DELETE CASCADE,
    exercise_id         BIGINT      NOT NULL REFERENCES exercises(id),
    position            INT         NOT NULL,
    goal_hold_seconds   INT         NOT NULL,
    actual_hold_seconds INT         NOT NULL
);
```

---

## Index Summary

| Table | Index | Rationale |
|---|---|---|
| `users` | UNIQUE on `username` | Login lookup |
| `templates` | `(created_by)` | "My templates" list |
| `templates` | `(is_public)` | Browse public templates |
| `programs` | `(user_id, is_active)` | Find active program for a user |
| `workouts` | `(user_id, status, workout_date DESC)` | History list |
| `workouts` | `(program_id, phase_id, workout_type)` | Progression look-back |
| `workout_main_sets` | `(workout_id, exercise_id)` | Sets for an exercise in a workout |
| `workout_abs_sets` | `(workout_id, exercise_id)` | Same, for abs block |
| `program_exercise_state` | UNIQUE on `(program_id, exercise_id, block_type)` | State lookup |
| `progression_events` | `(workout_id)` | Post-workout indicator display |
| `progression_events` | `(exercise_id)` | Exercise trend history |

---

## Entity Relationship Summary

```
users
  └── programs
        ├── phases
        ├── program_exercise_state  → exercises
        └── workouts
              ├── workout_main_sets        → exercises
              │     └── workout_main_set_segments
              ├── workout_abs_sets         → exercises
              │     └── workout_abs_set_segments
              ├── workout_cardio_logs
              └── workout_stretch_logs     → exercises

exercises  (global catalogue)

templates  (shared, owned by a user)
  ├── template_main_exercises    → exercises
  ├── template_abs_exercises     → exercises
  ├── template_cardio_blocks
  └── template_stretch_blocks    → exercises

progression_events → workouts, exercises
```

---

## Key Design Decisions

### `workout_type` stored explicitly on workouts
Using a text enum (`'normal'`, `'deload'`) rather than a boolean avoids re-deriving the type from `programs.weeks_per_phase` on every progression query, protects historical data if the program config changes, and makes it easy to add new types (e.g. `'test'`, `'competition'`) in the future with only a CHECK constraint update.

### Set segments for mid-set weight drops
A drop set "55 lb x4 → 45 lb x6" is one logical set with two segments. This cleanly supports any number of drops without fixed-width columns.

### `program_exercise_state` separate from templates
Template rows define starting goal values. Once a program is running, each user's goals drift independently via progression events. Storing live state on the template would corrupt it for all users.

### Goals snapshotted onto workout set rows
The goal at workout creation is what the user was targeting. Snapshotting preserves the historical record even after future progression changes the goal.

### Deload target reps computed, not stored
Deload reps = `phase.rep_min - 2`. The application computes this when creating the workout and writes it into `workout_main_sets.goal_reps`. The phase columns are the single source of truth.

### Template copying
A copy creates a new `templates` row with `created_by` set to the copying user, then duplicates all child rows. No `source_template_id` link is maintained — the copy is fully independent.
