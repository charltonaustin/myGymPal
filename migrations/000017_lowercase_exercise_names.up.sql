-- For exercises that would collide after lowercasing, keep the one with the
-- highest goal_weight (break ties by keeping the lowest id), then delete the rest.
WITH ranked AS (
    SELECT id,
           ROW_NUMBER() OVER (
               PARTITION BY user_id, LOWER(name)
               ORDER BY goal_weight DESC, id ASC
           ) AS rn
    FROM exercises
)
DELETE FROM exercises WHERE id IN (SELECT id FROM ranked WHERE rn > 1);

-- Lowercase all exercise names.
UPDATE exercises SET name = LOWER(name);

-- Lowercase names in session_exercises and template_exercises.
UPDATE session_exercises SET name = LOWER(name);
UPDATE template_exercises SET name = LOWER(name);
