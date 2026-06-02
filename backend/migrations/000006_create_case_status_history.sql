-- +goose Up
-- +goose StatementBegin
CREATE TABLE case_status_history (
    id         UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    case_id    UUID         NOT NULL,
    old_status VARCHAR(30),
    new_status VARCHAR(30) NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT NOW()
);

ALTER TABLE case_status_history
    ADD CONSTRAINT fk_case_status_history_case
    FOREIGN KEY (case_id)
    REFERENCES cases(id);

CREATE INDEX idx_case_status_history_case_id ON case_status_history (case_id);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS case_status_history;
-- +goose StatementEnd
