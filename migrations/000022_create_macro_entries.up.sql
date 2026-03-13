CREATE TABLE macro_entries (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date        DATE NOT NULL,
    food_name   TEXT NOT NULL,
    protein     NUMERIC(6,1) NOT NULL DEFAULT 0,
    carbs       NUMERIC(6,1) NOT NULL DEFAULT 0,
    fat         NUMERIC(6,1) NOT NULL DEFAULT 0,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);
