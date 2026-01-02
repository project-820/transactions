-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS wallet_sync_tasks (
    wallet_id        BIGINT PRIMARY KEY REFERENCES wallets(id) ON DELETE CASCADE,
    user_id          TEXT NOT NULL,

    status           SMALLINT NOT NULL DEFAULT 1, -- 1=pending,2=running,3=idle,4=failed
    run_after        TIMESTAMPTZ NOT NULL DEFAULT now(),

    locked_until     TIMESTAMPTZ NULL,

    attempts         INT NOT NULL DEFAULT 0,
    last_error       TEXT NULL,

    last_started_at  TIMESTAMPTZ NULL,
    last_finished_at TIMESTAMPTZ NULL,

    created_at       TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at       TIMESTAMPTZ NOT NULL DEFAULT now()
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS wallet_sync_tasks;
-- +goose StatementEnd
