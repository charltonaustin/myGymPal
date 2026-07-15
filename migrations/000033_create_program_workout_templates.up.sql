CREATE TABLE program_workout_templates (
    id             BIGSERIAL PRIMARY KEY,
    program_id     BIGINT NOT NULL REFERENCES programs(id) ON DELETE CASCADE,
    workout_number INT NOT NULL CHECK (workout_number > 0),
    template_id    BIGINT NOT NULL REFERENCES templates(id) ON DELETE CASCADE,
    UNIQUE (program_id, workout_number)
);
