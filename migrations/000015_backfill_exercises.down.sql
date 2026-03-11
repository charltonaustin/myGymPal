-- Remove backfilled exercises (only those that came from session history,
-- identified by having a matching session_exercise record for the same user).
DELETE FROM exercises e
WHERE EXISTS (
    SELECT 1
    FROM session_exercises se
    JOIN sessions s ON s.id = se.session_id
    WHERE s.user_id = e.user_id AND se.name = e.name
);
