-- +goose Up
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    role TEXT NOT NULL DEFAULT 'student',
    created_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT users_role_check
        CHECK (role IN ('student', 'employee', 'admin', 'superadmin'))
);

-- +goose Down
DROP TABLE IF EXISTS users;