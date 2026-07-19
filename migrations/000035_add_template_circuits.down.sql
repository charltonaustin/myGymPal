ALTER TABLE template_exercises DROP COLUMN IF EXISTS work_seconds;
ALTER TABLE template_exercises DROP COLUMN IF EXISTS circuit_id;

DROP INDEX IF EXISTS idx_template_circuits_template_id;

DROP TABLE IF EXISTS template_circuits;
