-- name: CreateHearing :one
INSERT INTO hearings (case_id, title, description, type, location, scheduled_at)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetHearingByID :one
SELECT * FROM hearings
WHERE id = $1
LIMIT 1;

-- name: CountHearingsByCaseID :one
SELECT COUNT(*) FROM hearings
WHERE case_id = $1;

-- name: ListHearingsByCaseID :many
SELECT * FROM hearings
WHERE case_id = $1
ORDER BY scheduled_at ASC;
