-- +goose Up
CREATE TABLE IF NOT EXISTS asset_tags (
    asset_id UUID NOT NULL REFERENCES assets(id) ON DELETE CASCADE,
    tag TEXT NOT NULL,
    PRIMARY KEY (asset_id, tag)
);

CREATE INDEX IF NOT EXISTS idx_asset_tags_tag ON asset_tags(tag); 

-- +goose Down
DROP INDEX IF EXISTS idx_asset_tags_tag;

DROP TABLE IF EXISTS asset_tags;