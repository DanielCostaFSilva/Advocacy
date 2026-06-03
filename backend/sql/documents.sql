-- name: CreateDocument :one
INSERT INTO documents (case_id, name, description, type, file_name, mime_type, file_size, storage_key)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
RETURNING *;

-- name: GetDocumentByID :one
SELECT * FROM documents
WHERE id = $1
LIMIT 1;

-- name: ListDocumentsByCaseID :many
SELECT * FROM documents
WHERE case_id = $1
ORDER BY created_at ASC;

-- name: CountDocumentsByCaseID :one
SELECT COUNT(*) FROM documents
WHERE case_id = $1;
