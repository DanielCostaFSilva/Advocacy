-- name: CreateClient :one
INSERT INTO clients (name, cpf, email, phone)
VALUES ($1, $2, $3, $4)
RETURNING *;

-- name: GetClientByID :one
SELECT * FROM clients
WHERE id = $1
LIMIT 1;

-- name: GetClientByCPF :one
SELECT * FROM clients
WHERE cpf = $1
LIMIT 1;
