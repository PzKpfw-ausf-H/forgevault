-- +goose Up
CREATE TABLE IF NOT EXISTS refresh_sessions (
    id uuid primary key,
    user_id uuid not null references users(id) ON DELETE CASCADE,
    token_hash text not null unique,
    expires_at timestamptz not null,
    created_at timestamptz not null,
    revoked_at timestamptz not null,
    replaced_by uuid null
);

CREATE INDEX IF NOT EXISTS idx_refresh_sessions_user_id ON refresh_sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_refresh_sessions_expires_at ON refresh_sessions(expires_at);
-- +goose Down
DROP INDEX IF EXISTS idx_refresh_sessions_user_id;
DROP INDEX IF EXISTS idx_refresh_sessions_expires_at;

DROP TABLE IF EXISTS refresh_sessions;