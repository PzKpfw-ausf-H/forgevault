-- +goose Up
CREATE TABLE IF NOT EXISTS asset_files (
    id uuid primary key,
    asset_id uuid not null references assets(id) on delete cascade,
    version int not null,
    filename text not null,
    size_bytes bigint not null,
    content_type text not null,
    storage_key text not null,
    checksum text null,
    created_at timestamptz not null
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_asset_files_asset_version ON asset_files(asset_id, version);
CREATE INDEX IF NOT EXISTS idx_asset_files_asset_id ON asset_files(asset_id);
-- +goose Down
DROP INDEX IF EXISTS idx_asset_files_asset_version;
DROP INDEX IF EXISTS idx_asset_files_asset_id;

DROP TABLE IF EXISTS asset_files;