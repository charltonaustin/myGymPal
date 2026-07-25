ALTER TABLE session_exercises DROP COLUMN IF EXISTS work_seconds;
ALTER TABLE session_exercises DROP COLUMN IF EXISTS circuit_id;

DROP INDEX IF EXISTS idx_session_circuits_session_id;

DROP TABLE IF EXISTS session_circuits;
