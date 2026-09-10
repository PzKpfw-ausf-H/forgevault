-- +goose Up
CREATE TABLE IF NOT EXISTS assets (
    id UUID PRIMARY KEY,
    title TEXT NOT NULL,
    description TEXT NOT NULL DEFAULT '',

    author_name TEXT NOT NULL DEFAULT '',
    uploaded_by UUID NOT NULL,

    type TEXT NOT NULL,
    processing_status TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT asset_uploaded_by_fk
        FOREIGN KEY (uploaded_by)
        REFERENCES users(id),

    CONSTRAINT asset_type_check
        CHECK (type IN ('3d', '2d', 'audio', 'vfx', 'doc')),

    CONSTRAINT asset_processing_status_check
        CHECK (processing_status IN ('pending', 'processing', 'ready', 'failed'))
);

CREATE INDEX IF NOT EXISTS idx_assets_uploaded_by ON assets(uploaded_by);
CREATE INDEX IF NOT EXISTS idx_assets_type ON assets(type);

-- +goose Down
DROP INDEX IF EXISTS idx_assets_uploaded_by;
DROP INDEX IF EXISTS idx_assets_type;

DROP TABLE IF EXISTS assets;