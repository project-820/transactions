-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS wallets (
    id          BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,
    user_id     TEXT NOT NULL,
    chain       TEXT NOT NULL,
    address     TEXT NOT NULL,
    label       TEXT,

    status      SMALLINT NOT NULL DEFAULT 1, -- 1=active, 2=disabled

    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS wallets;

-- +goose StatementEnd
