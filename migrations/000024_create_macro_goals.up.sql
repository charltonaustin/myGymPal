CREATE TABLE macro_goals (
    id         BIGSERIAL PRIMARY KEY,
    user_id    BIGINT NOT NULL UNIQUE,
    protein    NUMERIC(7,1) NOT NULL DEFAULT 0,
    carbs      NUMERIC(7,1) NOT NULL DEFAULT 0,
    fat        NUMERIC(7,1) NOT NULL DEFAULT 0,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
