-- +goose Up
CREATE TABLE IF NOT EXISTS asset_tags (
    asset_id uuid not null references assets(id) ON DELETE CASCADE,
    tag text not null,
    primary key (asset_id, tag)
);

CREATE INDEX IF NOT EXISTS idx_asset_tags_tag ON asset_tags(tag);
CREATE INDEX IF NOT EXISTS idx_asset_tags_asset_id ON asset_tags(asset_id);
-- +goose Down
DROP INDEX IF EXISTS idx_asset_tags_tag;
DROP INDEX IF EXISTS idx_asset_tags_asset_id;

DROP TABLE IF EXISTS asset_tags;

