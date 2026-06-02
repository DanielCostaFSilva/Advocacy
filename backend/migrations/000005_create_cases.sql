-- +goose Up
-- +goose StatementBegin
CREATE TABLE cases (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id   UUID         NOT NULL,
    number      VARCHAR(100) NOT NULL,
    title       VARCHAR(255) NOT NULL,
    description TEXT         NOT NULL DEFAULT '',
    court       VARCHAR(255) NOT NULL,
    status      VARCHAR(30)  NOT NULL DEFAULT 'draft',
    created_at  TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP    NOT NULL DEFAULT NOW()
);

ALTER TABLE cases
    ADD CONSTRAINT fk_cases_client
    FOREIGN KEY (client_id)
    REFERENCES clients(id);

CREATE UNIQUE INDEX uk_cases_number ON cases (number);
CREATE INDEX idx_cases_client_id ON cases (client_id);
CREATE INDEX idx_cases_status ON cases (status);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS cases;
-- +goose StatementEnd
