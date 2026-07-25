CREATE TABLE IF NOT EXISTS session_circuits (
    id                 BIGSERIAL PRIMARY KEY,
    session_id         BIGINT NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    name               VARCHAR(255) NOT NULL,
    rounds             INT NOT NULL DEFAULT 1 CHECK (rounds >= 1),
    transition_seconds INT NOT NULL DEFAULT 0 CHECK (transition_seconds >= 0),
    sort_order         INT NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_session_circuits_session_id ON session_circuits (session_id);

ALTER TABLE session_exercises
    ADD COLUMN IF NOT EXISTS circuit_id BIGINT REFERENCES session_circuits(id) ON DELETE SET NULL;

ALTER TABLE session_exercises
    ADD COLUMN IF NOT EXISTS work_seconds INT NOT NULL DEFAULT 0 CHECK (work_seconds >= 0);
