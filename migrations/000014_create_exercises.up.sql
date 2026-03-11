CREATE TABLE exercises (
    id            BIGSERIAL PRIMARY KEY,
    user_id       BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name          TEXT NOT NULL,
    is_bodyweight BOOLEAN NOT NULL DEFAULT FALSE,
    goal_weight   NUMERIC(6,2) NOT NULL DEFAULT 0,
    weight_unit   VARCHAR(8) NOT NULL DEFAULT 'lb',
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, name)
);
