CREATE TABLE session_exercises (
    id              BIGSERIAL PRIMARY KEY,
    session_id      BIGINT NOT NULL REFERENCES sessions(id) ON DELETE CASCADE,
    name            TEXT NOT NULL,
    is_bodyweight   BOOLEAN NOT NULL DEFAULT FALSE,
    goal_weight     NUMERIC(6,2),
    weight_unit     VARCHAR(8) NOT NULL DEFAULT 'lb',
    sort_order      INT NOT NULL DEFAULT 0
);

CREATE TABLE session_sets (
    id                      BIGSERIAL PRIMARY KEY,
    session_exercise_id     BIGINT NOT NULL REFERENCES session_exercises(id) ON DELETE CASCADE,
    set_number              INT NOT NULL,
    actual_weight           NUMERIC(6,2),
    weight_unit             VARCHAR(8) NOT NULL DEFAULT 'lb',
    actual_reps             INT NOT NULL,
    UNIQUE (session_exercise_id, set_number)
);
