CREATE TABLE sessions (
    id             BIGSERIAL PRIMARY KEY,
    program_id     BIGINT NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
    user_id        BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    phase_number   INT NOT NULL,
    week_number    INT NOT NULL,
    workout_number INT NOT NULL,
    is_deload      BOOLEAN NOT NULL DEFAULT FALSE,
    date           DATE NOT NULL DEFAULT CURRENT_DATE,
    created_at     TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at     TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
