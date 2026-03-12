CREATE TABLE cardio_logs (
    id BIGSERIAL PRIMARY KEY,
    session_exercise_id BIGINT NOT NULL REFERENCES session_exercises(id) ON DELETE CASCADE,
    cardio_type VARCHAR(100) NOT NULL DEFAULT '',
    goal_duration INT NOT NULL DEFAULT 0,
    actual_duration INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);
