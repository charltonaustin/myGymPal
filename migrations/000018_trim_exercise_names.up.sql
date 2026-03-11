-- Trim whitespace from all exercise names.
UPDATE exercises        SET name = TRIM(name) WHERE name != TRIM(name);
UPDATE session_exercises SET name = TRIM(name) WHERE name != TRIM(name);
UPDATE template_exercises SET name = TRIM(name) WHERE name != TRIM(name);

-- Backfill any exercises that appeared in sessions after the first backfill ran.
INSERT INTO exercises (user_id, name, is_bodyweight, goal_weight, weight_unit)
SELECT DISTINCT ON (s.user_id, TRIM(se.name))
    s.user_id,
    TRIM(se.name),
    se.is_bodyweight,
    COALESCE(se.goal_weight, 0),
    se.weight_unit
FROM session_exercises se
JOIN sessions s ON s.id = se.session_id
WHERE TRIM(se.name) != ''
ORDER BY s.user_id, TRIM(se.name), se.id DESC
ON CONFLICT (user_id, name) DO NOTHING;
