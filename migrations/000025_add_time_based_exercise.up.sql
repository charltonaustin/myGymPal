ALTER TABLE exercises
    ADD COLUMN is_time_based BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN goal_seconds  INT     NOT NULL DEFAULT 0;

ALTER TABLE session_exercises
    ADD COLUMN is_time_based BOOLEAN NOT NULL DEFAULT FALSE,
    ADD COLUMN goal_seconds  INT     NOT NULL DEFAULT 0;

ALTER TABLE session_sets
    ADD COLUMN actual_seconds INT NOT NULL DEFAULT 0;
