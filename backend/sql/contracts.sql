-- name: CreateContract :one
INSERT INTO contracts (client_id, case_id, title, description, type, amount, start_date, end_date)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetContractByID :one
SELECT * FROM contracts
WHERE id = $1;

-- name: ListContractsByCaseID :many
SELECT * FROM contracts
WHERE case_id = $1
ORDER BY created_at DESC;

-- name: UpdateContract :one
UPDATE contracts
SET title = $2, description = $3, type = $4, amount = $5, start_date = $6, end_date = $7, active = $8, updated_at = NOW()
WHERE id = $1
RETURNING *;
