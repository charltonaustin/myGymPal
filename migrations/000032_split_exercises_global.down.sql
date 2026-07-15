CREATE TABLE exercises_old (
    id            BIGSERIAL    PRIMARY KEY,
    user_id       BIGINT       NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name          TEXT         NOT NULL,
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
    UNIQUE (user_id, name)
);

INSERT INTO exercises_old (user_id, name, is_bodyweight, is_time_based, goal_weight, weight_unit, goal_seconds, goal_rep_min, goal_rep_max, default_block)
SELECT g.user_id, e.name, g.is_bodyweight, g.is_time_based, g.goal_weight, g.weight_unit, g.goal_seconds, g.goal_rep_min, g.goal_rep_max, g.default_block
FROM user_exercise_goals g
JOIN exercises e ON e.id = g.exercise_id;

DROP TABLE user_exercise_goals;
DROP TABLE exercises;
ALTER TABLE exercises_old RENAME TO exercises;
