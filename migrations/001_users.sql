-- +goose Up
CREATE EXTENSION IF NOT EXISTS citext;

CREATE TABLE IF NOT EXISTS users (
    id uuid primary key,
    email CITEXT not null unique,
    password_hash text not null,
    created_at timestamptz not null
);

-- +goose Down
DROP TABLE IF EXISTS users;