-- +goose Up
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    user_name VARCHAR(30) NOT NULL,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    password_hash VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL UNIQUE
);

-- +goose Down
DROP TABLE IF EXISTS users;