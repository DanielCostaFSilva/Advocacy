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
