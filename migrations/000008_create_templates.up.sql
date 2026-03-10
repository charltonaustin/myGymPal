CREATE TABLE templates (
    id         BIGSERIAL PRIMARY KEY,
    name       VARCHAR(255) NOT NULL,
    focus      VARCHAR(255) NOT NULL DEFAULT '',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE template_exercises (
    id            BIGSERIAL PRIMARY KEY,
    template_id   BIGINT NOT NULL REFERENCES templates(id) ON DELETE CASCADE,
    name          VARCHAR(255) NOT NULL,
    is_bodyweight BOOLEAN NOT NULL DEFAULT FALSE,
    goal_weight   NUMERIC(8,2) NOT NULL DEFAULT 0,
    rep_min       INT NOT NULL,
    rep_max       INT NOT NULL,
    sort_order    INT NOT NULL DEFAULT 0
);
