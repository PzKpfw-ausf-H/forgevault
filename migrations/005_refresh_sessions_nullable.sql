-- +goose Up
ALTER TABLE refresh_sessions
  ALTER COLUMN revoked_at DROP NOT NULL,
  ALTER COLUMN replaced_by DROP NOT NULL;

-- +goose Down
ALTER TABLE refresh_sessions
  ALTER COLUMN revoked_at SET NOT NULL,
  ALTER COLUMN replaced_by SET NOT NULL;