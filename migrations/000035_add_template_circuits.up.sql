CREATE TABLE IF NOT EXISTS template_circuits (
    id                 BIGSERIAL PRIMARY KEY,
    template_id        BIGINT NOT NULL REFERENCES templates(id) ON DELETE CASCADE,
    name               VARCHAR(255) NOT NULL,
    rounds             INT NOT NULL DEFAULT 1 CHECK (rounds >= 1),
    transition_seconds INT NOT NULL DEFAULT 0 CHECK (transition_seconds >= 0),
    sort_order         INT NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_template_circuits_template_id ON template_circuits (template_id);

ALTER TABLE template_exercises
    ADD COLUMN IF NOT EXISTS circuit_id BIGINT REFERENCES template_circuits(id) ON DELETE SET NULL;

ALTER TABLE template_exercises
    ADD COLUMN IF NOT EXISTS work_seconds INT NOT NULL DEFAULT 0 CHECK (work_seconds >= 0);
