-- +goose Up
-- +goose StatementBegin
CREATE TABLE hearings (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    case_id     UUID         NOT NULL,
    title       VARCHAR(255) NOT NULL,
    description TEXT,
    type        VARCHAR(50)  NOT NULL,
    location    VARCHAR(255) NOT NULL,
    scheduled_at TIMESTAMP   NOT NULL,
    created_at  TIMESTAMP    NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMP    NOT NULL DEFAULT NOW()
);

ALTER TABLE hearings
    ADD CONSTRAINT fk_hearings_case
    FOREIGN KEY (case_id)
    REFERENCES cases(id);

CREATE INDEX idx_hearings_case_id ON hearings (case_id);
CREATE INDEX idx_hearings_scheduled_at ON hearings (scheduled_at);
CREATE INDEX idx_hearings_type ON hearings (type);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS hearings;
-- +goose StatementEnd
