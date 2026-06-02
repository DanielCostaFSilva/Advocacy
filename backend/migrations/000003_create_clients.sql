-- +goose Up
-- +goose StatementBegin
CREATE TABLE clients (
    id         UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    name       VARCHAR(255) NOT NULL,
    cpf        VARCHAR(11)  NOT NULL,
    email      VARCHAR(255) NOT NULL DEFAULT '',
    phone      VARCHAR(30)  NOT NULL DEFAULT '',
    created_at TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP    NOT NULL DEFAULT NOW()
);

CREATE UNIQUE INDEX uk_clients_cpf ON clients (cpf);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS clients;
-- +goose StatementEnd
