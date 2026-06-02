-- name: CreateCase :one
INSERT INTO cases (client_id, number, title, description, court, status)
VALUES ($1, $2, $3, $4, $5, $6)
RETURNING *;

-- name: GetCaseByID :one
SELECT * FROM cases
WHERE id = $1
LIMIT 1;

-- name: GetCaseByNumber :one
SELECT * FROM cases
WHERE number = $1
LIMIT 1;
