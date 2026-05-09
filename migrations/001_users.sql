-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id uuid primary key,
    email text not null unique,
    password_hash text not null,
    created_at timestampz not null
);

-- +goose Down
DROP TABLE IF EXISTS users;