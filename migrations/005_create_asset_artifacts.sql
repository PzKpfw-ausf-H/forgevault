-- +goose Up
CREATE TABLE IF NOT EXISTS asset_artifacts(
    id UUID PRIMARY KEY
    asset_id UUID NOT NULL,
    type TEXT NOT NULL,
    mime_type TEXT NOT NULL,

    size BIGINT NOT NULL,

    checksum TEXT NOT NULL,
    storage_key TEXT NOT NULL UNIQUE,

    created_at TIMESTAMPZ NOT NULL,

    CONSTRAINT asset_artifacts_asset_fk
        FOREIGN KEY (asset_id)
        REFERENCES assets(id)
        ON DELETE CASCADE,

    CONSTRAINT asset_artifacts_type_check
        CHECK (type IN ('viewer_model', 'thumbnail', 'preview_image')),
    
    CONSTRAINT asset_artifact_size_check
        CHECK (size >= 0),
);

CREATE INDEX IF NOT EXISTS idx_asset_artifacts_asset_id ON asset_artifacts(asset_id);

-- +goose Down
DROP INDEX IF EXISTS idx_asset_artifacts_asset_id;

DROP TABLE IF EXISTS asset_artifacts;
