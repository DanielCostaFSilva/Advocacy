-- name: GetTimelineByCaseID :many
SELECT * FROM timeline_events
WHERE case_id = $1
ORDER BY created_at ASC;

-- name: CreateTimelineEvent :one
INSERT INTO timeline_events (case_id, type, description, metadata)
VALUES ($1, $2, $3, $4)
RETURNING *;
