-- Create global exercises table (name only, shared across all users)
CREATE TABLE IF NOT EXISTS exercises_new (
    id         BIGSERIAL   PRIMARY KEY,
    name       TEXT        NOT NULL UNIQUE,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Create per-user goals/configuration table
CREATE TABLE IF NOT EXISTS user_exercise_goals (
    id            BIGSERIAL    PRIMARY KEY,
    user_id       BIGINT       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    exercise_id   BIGINT       NOT NULL REFERENCES exercises_new(id) ON DELETE CASCADE,
    is_bodyweight BOOLEAN      NOT NULL DEFAULT FALSE,
    is_time_based BOOLEAN      NOT NULL DEFAULT FALSE,
    goal_weight   NUMERIC(6,2) NOT NULL DEFAULT 0,
    weight_unit   VARCHAR(8)   NOT NULL DEFAULT 'lb',
    goal_seconds  INT          NOT NULL DEFAULT 0,
    goal_rep_min  INT          NOT NULL DEFAULT 0,
    goal_rep_max  INT          NOT NULL DEFAULT 0,
    default_block VARCHAR(20)  NOT NULL DEFAULT 'main',
    created_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at    TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    UNIQUE (user_id, exercise_id)
);

-- Migrate data and swap tables.
-- Uses EXECUTE with dollar-quoting so column references in the old schema are
-- resolved at runtime, making this block safe to re-run if exercises is already
-- in the new format (i.e. no user_id column).
DO $outer$
BEGIN
    IF EXISTS (
        SELECT 1 FROM information_schema.columns
        WHERE table_schema = 'public' AND table_name = 'exercises' AND column_name = 'user_id'
    ) THEN
        EXECUTE $q$
            INSERT INTO exercises_new (name)
            SELECT DISTINCT LOWER(TRIM(name)) FROM exercises
            ON CONFLICT (name) DO NOTHING
        $q$;
        EXECUTE $q$
            INSERT INTO user_exercise_goals
                (user_id, exercise_id, is_bodyweight, is_time_based, goal_weight,
                 weight_unit, goal_seconds, goal_rep_min, goal_rep_max, default_block)
            SELECT e.user_id, n.id, e.is_bodyweight, e.is_time_based,
                   e.goal_weight, e.weight_unit, e.goal_seconds,
                   e.goal_rep_min, e.goal_rep_max, e.default_block
            FROM exercises e
            JOIN exercises_new n ON n.name = LOWER(TRIM(e.name))
            ON CONFLICT (user_id, exercise_id) DO NOTHING
        $q$;
        EXECUTE $q$ DROP TABLE exercises $q$;
        EXECUTE $q$ ALTER TABLE exercises_new RENAME TO exercises $q$;
    ELSE
        -- exercises already in new format; drop exercises_new if left over
        DROP TABLE IF EXISTS exercises_new;
    END IF;
END $outer$;
