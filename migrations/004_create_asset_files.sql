-- +goose Up
CREATE TABLE IF NOT EXISTS asset_files(
    id UUID PRIMARY KEY,
    asset_id UUID NOT NULL,

    role TEXT NOT NULL,
    texture_type TEXT NULL,

    original_name TEXT NOT NULL,
    mime_type TEXT NOT NULL,
    extension TEXT NOT NULL,

    size BIGINT NOT NULL,
    storage_key TEXT NOT NULL UNIQUE,
    checksum TEXT NOT NULL,

    created_at TIMESTAMPTZ NOT NULL,

    CONSTRAINT asset_files_asset_fk
        FOREIGN KEY (asset_id)
        REFERENCES assets(id)
        ON DELETE CASCADE,

    CONSTRAINT asset_files_role_check
        CHECK (role IN ('main', 'texture', 'attachment', 'preview_source')),

    CONSTRAINT asset_file_size_check
        CHECK (size >= 0),

    CONSTRAINT asset_files_texture_type_presence_check
        CHECK (
            (role = 'texture' AND texture_type IS NOT NULL)
            OR
            (role <> 'texture' AND texture_type IS NULL)    
        ),

    CONSTRAINT asset_files_texture_type_check
        CHECK (texture_type IS NULL
        OR texture_type IN (
            'base_color',
            'normal',
            'roughness',
            'metallic',
            'ao',
            'emissive',
            'opacity',
            'height',
            'orm',
            'other'
        ))
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_asset_files_one_main_per_asset
    ON asset_files(asset_id)
    WHERE role = 'main';

CREATE INDEX IF NOT EXISTS idx_asset_files_asset_id ON asset_files(asset_id);

-- +goose Down
DROP INDEX IF EXISTS idx_asset_files_one_main_per_asset;
DROP INDEX IF EXISTS idx_asset_files_asset_id;

DROP TABLE IF EXISTS asset_files;
