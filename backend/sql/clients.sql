-- name: CreateClient :one
INSERT INTO clients (name, cpf, email, phone)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetClientByID :one
SELECT * FROM clients
WHERE id = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: GetClientByCPF :one
SELECT * FROM clients
WHERE cpf = $1 AND deleted_at IS NULL
LIMIT 1;

-- name: UpdateClient :one
UPDATE clients
SET name = $2, email = $3, phone = $4, updated_at = NOW()
WHERE id = $1
RETURNING *;

-- name: SoftDeleteClient :exec
UPDATE clients
SET deleted_at = NOW()
WHERE id = $1 AND deleted_at IS NULL;

-- name: CountClients :one
SELECT COUNT(*) FROM clients
WHERE (sqlc.narg('name') IS NULL OR name ILIKE '%' || sqlc.narg('name') || '%')
  AND (sqlc.narg('cpf') IS NULL OR cpf = sqlc.narg('cpf'))
  AND deleted_at IS NULL;
