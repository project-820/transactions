-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS outbox_messages (
    id            BIGSERIAL   PRIMARY KEY,
    subject       TEXT        NOT NULL,
    event_name    TEXT        NOT NULL,
    status        TEXT        NOT NULL DEFAULT 'new',
    payload       BYTEA       NOT NULL,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    processed_at  TIMESTAMPTZ
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS outbox_messages;

-- +goose StatementEnd
