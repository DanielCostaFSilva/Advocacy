-- +goose Up
-- +goose StatementBegin
CREATE TABLE contracts (
    id          UUID             PRIMARY KEY DEFAULT gen_random_uuid(),
    client_id   UUID             NOT NULL,
    case_id     UUID             NOT NULL,
    title       VARCHAR(255)     NOT NULL,
    description TEXT,
    type        VARCHAR(50)      NOT NULL,
    amount      NUMERIC(18,2)    NOT NULL,
    start_date  DATE             NOT NULL,
    end_date    DATE,
    active      BOOLEAN          NOT NULL DEFAULT TRUE,
    created_at  TIMESTAMP        NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP        NOT NULL DEFAULT NOW()
);

ALTER TABLE contracts
    ADD CONSTRAINT fk_contracts_client
    FOREIGN KEY (client_id)
    REFERENCES clients(id);

ALTER TABLE contracts
    ADD CONSTRAINT fk_contracts_case
    FOREIGN KEY (case_id)
    REFERENCES cases(id);

CREATE INDEX idx_contracts_client_id ON contracts (client_id);
CREATE INDEX idx_contracts_case_id ON contracts (case_id);
CREATE INDEX idx_contracts_type ON contracts (type);
CREATE INDEX idx_contracts_active ON contracts (active);
CREATE INDEX idx_contracts_created_at ON contracts (created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS contracts;
-- +goose StatementEnd
