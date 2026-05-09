-- +goose Up
CREATE TABLE IF NOT EXISTS assets (
    id uuid primary key,
    title text not null,
    description text not null default '',
    type text not null,
    author_id uuid not null references users(id),
    created_at timestampz not null,
    updated_at timestampz not null
);

CREATE INDEX IF NOT EXISTS idx_assets_type ON assets(type);
CREATE INDEX IF NOT EXISTS idx_assets_author_id ON assets(author_id);
CREATE INDEX IF NOT EXISTS idx_assets_created_at ON assets(created_at);

-- +goose Down
DROP INDEX IF EXISTS idx_assets_type;
DROP INDEX IF EXISTS idx_assets_author_id;
DROP INDEX IF EXISTS idx_assets_created_at;

DROP TABLE IF EXISTS assets;