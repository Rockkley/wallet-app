-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS wallets
(
    id      UUID PRIMARY KEY,
    balance BIGINT NOT NULL DEFAULT 0
);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE wallets;
-- +goose StatementEnd
