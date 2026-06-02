-- +goose Up
-- +goose StatementBegin
CREATE TABLE timeline_events (
    id          UUID         PRIMARY KEY DEFAULT gen_random_uuid(),
    case_id     UUID         NOT NULL,
    type        VARCHAR(50)  NOT NULL,
    description TEXT         NOT NULL,
    metadata    JSONB,
    created_at  TIMESTAMP    NOT NULL DEFAULT NOW()
);

ALTER TABLE timeline_events
    ADD CONSTRAINT fk_timeline_events_case
    FOREIGN KEY (case_id)
    REFERENCES cases(id);

CREATE INDEX idx_timeline_case_id ON timeline_events (case_id);
CREATE INDEX idx_timeline_created_at ON timeline_events (created_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS timeline_events;
-- +goose StatementEnd
