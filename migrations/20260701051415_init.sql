-- +goose Up
CREATE TABLE refresh_tokens(
    value VARCHAR(100) NOT NULL UNIQUE, --TODO: какой размер тут?
    user_id INTEGER NOT NULL, 
    device_id VARCHAR(100), 
    expired_at TIMESTAMP,
    revoked BOOLEAN NOT NULL DEFAULT FALSE,
    used BOOLEAN NOT NULL DEFAULT FALSE,
    session_id VARCHAR(100) NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE refresh_tokens;