CREATE TABLE phases (
    id           BIGSERIAL PRIMARY KEY,
    program_id   BIGINT NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
    phase_number INT    NOT NULL CHECK (phase_number > 0),
    rep_min      INT    NOT NULL DEFAULT 0,
    rep_max      INT    NOT NULL DEFAULT 0,
    UNIQUE (program_id, phase_number)
);
