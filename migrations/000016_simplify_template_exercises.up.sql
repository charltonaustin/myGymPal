ALTER TABLE template_exercises
    DROP COLUMN IF EXISTS goal_weight,
    DROP COLUMN IF EXISTS weight_unit,
    DROP COLUMN IF EXISTS rep_min,
    DROP COLUMN IF EXISTS rep_max;
