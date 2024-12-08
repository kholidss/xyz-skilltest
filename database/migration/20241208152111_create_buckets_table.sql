-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS buckets (
    id VARCHAR(36) NOT NULL PRIMARY KEY,
    file_name VARCHAR(255) NOT NULL,
    identifier_id VARCHAR(36) NOT NULL,
    identifier_name VARCHAR(255) NOT NULL,
    mimetype VARCHAR(100) NOT NULL,
    provider VARCHAR(100) NOT NULL,
    url TEXT NOT NULL,
    path TEXT NOT NULL,
    support_data JSON NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP NOT NULL,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP NOT NULL,
    is_deleted BOOLEAN DEFAULT FALSE,
    deleted_at TIMESTAMP NULL,
    KEY file_name(file_name),
    KEY identifier_id(identifier_id),
    KEY identifier_name(identifier_name),
    KEY mimetype(mimetype),
    KEY provider(provider),
    KEY created_at(created_at),
    KEY updated_at(updated_at),
    KEY is_deleted(is_deleted),
    KEY deleted_at(deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS buckets;
-- +goose StatementEnd
