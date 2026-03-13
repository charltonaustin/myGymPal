ALTER TABLE session_sets DROP COLUMN actual_seconds;
ALTER TABLE session_exercises DROP COLUMN goal_seconds, DROP COLUMN is_time_based;
ALTER TABLE exercises DROP COLUMN goal_seconds, DROP COLUMN is_time_based;
