-- +goose Up
ALTER TABLE playlist ALTER COLUMN thumbnail_url SET NOT NULL;
ALTER TABLE playlist ALTER COLUMN channel SET NOT NULL;

-- +goose Down
ALTER TABLE playlist ALTER COLUMN thumbnail_url DROP NOT NULL;
ALTER TABLE playlist ALTER COLUMN channel DROP NOT NULL;
