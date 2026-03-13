CREATE TABLE body_weights (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    date        DATE NOT NULL,
    weight      NUMERIC(6,2) NOT NULL,
    weight_unit VARCHAR(8) NOT NULL DEFAULT 'lb',
    UNIQUE (user_id, date)
);
