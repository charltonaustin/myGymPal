CREATE TABLE users (
    id            BIGSERIAL   PRIMARY KEY,
    username      TEXT        NOT NULL UNIQUE,
    password_hash TEXT        NOT NULL,
    weight_unit   TEXT        NOT NULL DEFAULT 'lb'
                      CHECK (weight_unit IN ('lb', 'kg')),
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
