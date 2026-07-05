# myGymPal — PostgreSQL Data Model

## Overview

This document describes the complete relational schema for myGymPal. The stack is Go (Beego v2) with `lib/pq` connecting
to PostgreSQL. All tables use `BIGSERIAL` surrogate primary keys. Timestamps are stored as `TIMESTAMPTZ`. Weight values
are stored in the user's preferred unit (lb or kg) as entered — no normalisation to a canonical unit is applied at the
database layer.

---

## Tables

### `users`

```sql
CREATE TABLE users (
    id            BIGSERIAL   PRIMARY KEY,
    username      TEXT        NOT NULL UNIQUE,
    password_hash TEXT        NOT NULL,
    weight_unit   TEXT        NOT NULL DEFAULT 'lb'
        CHECK (weight_unit IN ('lb', 'kg')),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

### `session`

Beego HTTP session store. Not application data.

```sql
CREATE TABLE session (
    session_key    CHAR(64)  NOT NULL,
    session_data   BYTEA,
    session_expiry TIMESTAMP NOT NULL,
    CONSTRAINT session_pkey PRIMARY KEY (session_key)
);
```

---

### `programs`

A user's training program. One user can have multiple programs.

```sql
CREATE TABLE programs (
    id                BIGSERIAL   PRIMARY KEY,
    user_id           BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name              TEXT        NOT NULL,
    start_date        DATE        NOT NULL,
    num_phases        INT         NOT NULL CHECK (num_phases > 0),
    weeks_per_phase   INT         NOT NULL DEFAULT 8 CHECK (weeks_per_phase > 0),
    workouts_per_week INT         NOT NULL DEFAULT 4 CHECK (workouts_per_week > 0),
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

### `phases`

One row per phase within a program. Each phase has its own rep range, default set count, and rest period.

```sql
CREATE TABLE phases (
    id           BIGSERIAL PRIMARY KEY,
    program_id   BIGINT    NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
    phase_number INT       NOT NULL CHECK (phase_number > 0),
    rep_min      INT       NOT NULL DEFAULT 0,
    rep_max      INT       NOT NULL DEFAULT 0,
    default_sets INT       NOT NULL DEFAULT 3,
    rest_seconds INT       NOT NULL DEFAULT 0,
    UNIQUE (program_id, phase_number)
);
```

The last week of each phase (`week_number = weeks_per_phase`) is always the deload week. This is a computation, not
stored on the phase.

---

### `templates`

Reusable workout blueprints. No user ownership — all templates are visible to all authenticated users.

```sql
CREATE TABLE templates (
    id         BIGSERIAL    PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    focus      VARCHAR(255) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
```

---

### `template_exercises`

Ordered exercises within a template. Stores identity only — goal weight and rep targets come from the user's exercise
library and the current phase at workout creation time.

```sql
CREATE TABLE template_exercises (
    id            BIGSERIAL    PRIMARY KEY,
    template_id   BIGINT       NOT NULL REFERENCES templates(id) ON DELETE CASCADE,
    name          VARCHAR(255) NOT NULL,
    is_bodyweight BOOLEAN      NOT NULL DEFAULT FALSE,
    is_time_based BOOLEAN      NOT NULL DEFAULT FALSE,
    sort_order    INT          NOT NULL DEFAULT 0,
    block         VARCHAR(20)  NOT NULL DEFAULT 'main'
);
```

`name` is always stored lowercase and trimmed. Exercise names are the coupling point between templates and a user's
exercise library — there is no FK to `exercises`.

---

### `sessions`

A single workout instance linked to a program and phase.

```sql
CREATE TABLE sessions (
    id             BIGSERIAL   PRIMARY KEY,
    program_id     BIGINT      NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
    user_id        BIGINT      NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    phase_number   INT         NOT NULL,
    week_number    INT         NOT NULL,
    workout_number INT         NOT NULL,
    is_deload      BOOLEAN     NOT NULL DEFAULT FALSE,
    date           DATE        NOT NULL DEFAULT CURRENT_DATE,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

---

### `session_exercises`

Exercises logged within a session. Goals are snapshotted from the user's exercise library and the current phase at
session creation so historical targets are preserved even after the library is updated.

```sql
CREATE TABLE session_exercises (
    id            BIGSERIAL    PRIMARY KEY,
    session_id    BIGINT       NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    name          TEXT         NOT NULL,
    is_bodyweight BOOLEAN      NOT NULL DEFAULT FALSE,
    is_time_based BOOLEAN      NOT NULL DEFAULT FALSE,
    goal_weight   NUMERIC(6,2),
    weight_unit   VARCHAR(8)   NOT NULL DEFAULT 'lb',
    goal_reps     INT          NOT NULL DEFAULT 0,
    goal_seconds  INT          NOT NULL DEFAULT 0,
    block         VARCHAR(20)  NOT NULL DEFAULT 'main',
    sort_order    INT          NOT NULL DEFAULT 0
);
```

---

### `session_sets`

Individual sets logged within a session exercise.

```sql
CREATE TABLE session_sets (
    id                  BIGSERIAL   PRIMARY KEY,
    session_exercise_id BIGINT      NOT NULL REFERENCES session_exercises(id) ON DELETE CASCADE,
    set_number          INT         NOT NULL,
    actual_weight       NUMERIC(6,2),
    weight_unit         VARCHAR(8)  NOT NULL DEFAULT 'lb',
    actual_reps         INT         NOT NULL,
    actual_seconds      INT         NOT NULL DEFAULT 0,
    activity_type       VARCHAR(50) NOT NULL DEFAULT '',
    UNIQUE (session_exercise_id, set_number)
);
```

---

### `cardio_logs`

Cardio entries linked to a session exercise.

```sql
CREATE TABLE cardio_logs (
    id                  BIGSERIAL    PRIMARY KEY,
    session_exercise_id BIGINT       NOT NULL REFERENCES session_exercises(id) ON DELETE CASCADE,
    cardio_type         VARCHAR(100) NOT NULL DEFAULT '',
    goal_duration       INT          NOT NULL DEFAULT 0,
    actual_duration     INT          NOT NULL DEFAULT 0,
    created_at          TIMESTAMP    NOT NULL DEFAULT NOW()
);
```

---

### `exercises`

Per-user exercise library. Each user maintains their own named list of exercises with personal goal targets.

```sql
CREATE TABLE exercises (
    id            BIGSERIAL    PRIMARY KEY,
    user_id       BIGINT       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name          TEXT         NOT NULL,
    is_bodyweight BOOLEAN      NOT NULL DEFAULT FALSE,
    is_time_based BOOLEAN      NOT NULL DEFAULT FALSE,
    goal_weight   NUMERIC(6,2) NOT NULL DEFAULT 0,
    weight_unit   VARCHAR(8)   NOT NULL DEFAULT 'lb',
    goal_seconds  INT          NOT NULL DEFAULT 0,
    goal_rep_min  INT          NOT NULL DEFAULT 0,
    goal_rep_max  INT          NOT NULL DEFAULT 0,
    default_block VARCHAR(20)  NOT NULL DEFAULT 'main',
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, name)
);
```

`name` is always stored lowercase and trimmed. The `UNIQUE (user_id, name)` constraint enforces one entry per exercise
name per user. There is no sharing between users — each user's library is fully independent.

---

### `body_weights`

Daily body weight log per user.

```sql
CREATE TABLE body_weights (
    id          BIGSERIAL    PRIMARY KEY,
    user_id     BIGINT       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date        DATE         NOT NULL,
    weight      NUMERIC(6,2) NOT NULL,
    weight_unit VARCHAR(8)   NOT NULL DEFAULT 'lb',
    UNIQUE (user_id, date)
);
```

---

### `macro_entries`

Individual food entries per user per day.

```sql
CREATE TABLE macro_entries (
    id             BIGSERIAL    PRIMARY KEY,
    user_id        BIGINT       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date           DATE         NOT NULL,
    food_name      TEXT         NOT NULL,
    protein        NUMERIC(6,1) NOT NULL DEFAULT 0,
    carbs          NUMERIC(6,1) NOT NULL DEFAULT 0,
    fat            NUMERIC(6,1) NOT NULL DEFAULT 0,
    serving_weight NUMERIC(7,1) NOT NULL DEFAULT 0,
    serving_unit   VARCHAR(8)   NOT NULL DEFAULT 'g',
    created_at     TIMESTAMP    NOT NULL DEFAULT NOW()
);
```

---

### `macro_goals`

Daily macro targets per user. One row per user (enforced by UNIQUE on `user_id`).

```sql
CREATE TABLE macro_goals (
    id         BIGSERIAL    PRIMARY KEY,
    user_id    BIGINT       NOT NULL UNIQUE,
    protein    NUMERIC(7,1) NOT NULL DEFAULT 0,
    carbs      NUMERIC(7,1) NOT NULL DEFAULT 0,
    fat        NUMERIC(7,1) NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ  NOT NULL DEFAULT NOW()
);
```

---

## Index Summary

| Table               | Index                                | Rationale                             |
|---------------------|--------------------------------------|---------------------------------------|
| `users`             | UNIQUE on `username`                 | Login lookup                          |
| `exercises`         | UNIQUE on `(user_id, name)`          | Per-user exercise lookup by name      |
| `phases`            | UNIQUE on `(program_id, phase_number)` | Phase lookup within a program       |
| `session_sets`      | UNIQUE on `(session_exercise_id, set_number)` | Sets within an exercise      |
| `body_weights`      | UNIQUE on `(user_id, date)`          | One entry per user per day            |
| `macro_goals`       | UNIQUE on `user_id`                  | Single goal row per user              |

---

## Entity Relationship Summary

```
users
  ├── programs
  │     └── phases
  │     └── sessions
  │           ├── session_exercises
  │           │     ├── session_sets
  │           │     └── cardio_logs
  ├── exercises       (per-user library)
  ├── body_weights
  ├── macro_entries
  └── macro_goals

templates             (global, no user ownership)
  └── template_exercises
```

---

## Key Design Decisions

### Exercise names are the coupling between templates and the exercise library

`template_exercises.name` is a plain text field, not a FK to `exercises`. When a session is created from a template,
the app looks up each exercise by name in the user's personal exercise library to copy goal weight, reps, and seconds
into the session snapshot. This means a user's library must contain matching names for template goals to populate —
new users start blank and build their library over time.

### Template exercises store identity only

`template_exercises` holds only `name`, `is_bodyweight`, `is_time_based`, `block`, and `sort_order`. Goal values
(weight, reps, seconds) are not stored on the template. They come from the user's exercise library at session creation
time, so each user gets their own personalised targets from the same shared template.

### Goals snapshotted onto session rows at creation

`session_exercises` copies goal weight, reps, and seconds from the exercise library at the moment the session is
created. This preserves the historical target even if the user later updates their exercise goals.

### `is_deload` stored explicitly on sessions

Deload status is stored as a boolean on `sessions` rather than computed from `week_number` and `programs.weeks_per_phase`.
This protects historical records if the program config is later edited.

### Time-based exercises are a flag, not a separate type

`is_time_based` is a boolean on `exercises`, `session_exercises`, and `template_exercises`. When true, `goal_seconds`
(and `actual_seconds` on sets) is the relevant target; `goal_weight`/`actual_reps` are unused. This avoids a separate
table or type hierarchy for what is a minor behavioural variant.

### Block stored on exercise entries, not inferred

`block` (`'main'`, `'abs'`, `'cardio'`, `'stretch'`) is stored on `session_exercises` and `template_exercises` rather
than looked up from the exercise library. The user's `exercises.default_block` is used as a suggestion at entry time
but is not authoritative — the block is confirmed per-use.