-- +goose Up
-- +goose StatementBegin

CREATE TABLE IF NOT EXISTS transactions (
    id            BIGINT GENERATED ALWAYS AS IDENTITY PRIMARY KEY,

    wallet_id      BIGINT NOT NULL REFERENCES wallets(id),
    user_id        TEXT NOT NULL,

    source         TEXT NOT NULL, -- onchain|cex|manual

    tx_hash        TEXT NULL,
    log_index      INT  NULL,

    external_id    TEXT NULL,

    asset_ref      TEXT NOT NULL,
    qty            NUMERIC(36,18) NOT NULL,

    price_usd      NUMERIC(36,18) NULL,
    fee_asset_ref  TEXT NULL,
    fee_qty        NUMERIC(36,18) NULL,

    occurred_at    TIMESTAMPTZ NOT NULL,

    meta           JSONB NULL,

    created_at     TIMESTAMPTZ NOT NULL DEFAULT now()
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin

DROP TABLE IF EXISTS transactions;

-- +goose StatementEnd
