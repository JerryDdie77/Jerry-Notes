-- +goose Up
CREATE TABLE pending_users (
    email         VARCHAR(50) PRIMARY KEY,
    user_name     VARCHAR(30)  NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    code          VARCHAR(6)   NOT NULL,
    created_at    TIMESTAMP    NOT NULL DEFAULT NOW()
);


-- +goose Down
DROP TABLE IF EXISTS pending_users;