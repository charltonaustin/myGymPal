-- Backfill the exercises library from historical session exercise data.
-- For each (user, exercise name) pair, take the most recently logged entry
-- so the goal weight reflects the latest value used.
INSERT INTO exercises (user_id, name, is_bodyweight, goal_weight, weight_unit)
SELECT DISTINCT ON (s.user_id, se.name)
    s.user_id,
    se.name,
    se.is_bodyweight,
    COALESCE(se.goal_weight, 0),
    se.weight_unit
FROM session_exercises se
JOIN sessions s ON s.id = se.session_id
ORDER BY s.user_id, se.name, se.id DESC
ON CONFLICT (user_id, name) DO NOTHING;
