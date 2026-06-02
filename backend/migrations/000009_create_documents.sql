-- +goose Up
-- +goose StatementBegin
CREATE TABLE documents (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    case_id     UUID         NOT NULL,
    name        VARCHAR(255) NOT NULL,
    description TEXT,
    type        VARCHAR(100) NOT NULL,
    file_name   VARCHAR(255) NOT NULL,
    mime_type   VARCHAR(255) NOT NULL,
    file_size   BIGINT       NOT NULL,
    storage_key VARCHAR(500) NOT NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP    NOT NULL DEFAULT NOW()
);

ALTER TABLE documents
    ADD CONSTRAINT fk_documents_case
    FOREIGN KEY (case_id)
    REFERENCES cases(id);

CREATE INDEX idx_documents_case_id ON documents (case_id);
CREATE INDEX idx_documents_type ON documents (type);
CREATE INDEX idx_documents_created_at ON documents (created_at);
CREATE INDEX idx_documents_storage_key ON documents (storage_key);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS documents;
-- +goose StatementEnd
