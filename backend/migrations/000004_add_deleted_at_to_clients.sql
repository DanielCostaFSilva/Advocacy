-- +goose Up
-- +goose StatementBegin
ALTER TABLE clients
ADD COLUMN deleted_at TIMESTAMP NULL;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
ALTER TABLE clients
DROP COLUMN IF EXISTS deleted_at;
-- +goose StatementEnd
