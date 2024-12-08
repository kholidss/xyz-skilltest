-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS transaction_credits (
    id VARCHAR(36) NOT NULL PRIMARY KEY,
    user_id VARCHAR(36) NOT NULL,
    merchant_id VARCHAR(36) NOT NULL,
    transaction_id VARCHAR(36) NOT NULL,
    asset_name VARCHAR(255) NOT NULL,
    limit_amount BIGINT NOT NULL DEFAULT 0,
    otr_amount BIGINT NOT NULL DEFAULT 0,
    fee_amount BIGINT NOT NULL DEFAULT 0,
    installment_amount BIGINT NOT NULL DEFAULT 0,
    interest_amount BIGINT NOT NULL DEFAULT 0,
    interest_percentage INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMP NULL,
    KEY user_id(user_id),
    KEY merchant_id(merchant_id),
    KEY transaction_id(transaction_id),
    KEY asset_name(asset_name),
    KEY created_at(created_at),
    KEY updated_at(updated_at),
    KEY is_deleted(is_deleted),
    KEY deleted_at(deleted_at),
    CONSTRAINT fk_transaction_credits_user_id FOREIGN KEY (user_id) REFERENCES users(id),
    CONSTRAINT fk_transaction_credits_merchant_id FOREIGN KEY (merchant_id) REFERENCES merchants(id),
    CONSTRAINT fk_transaction_credits_transaction_id FOREIGN KEY (transaction_id) REFERENCES transactions(id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS transaction_credits;
-- +goose StatementEnd
