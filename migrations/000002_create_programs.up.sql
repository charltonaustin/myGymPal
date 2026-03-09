CREATE TABLE programs (
    id          BIGSERIAL   PRIMARY KEY,
    user_id     BIGINT      NOT NULL REFERENCES users(id),
    name        TEXT        NOT NULL,
    start_date  DATE        NOT NULL,
    num_phases  INT         NOT NULL CHECK (num_phases > 0),
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
