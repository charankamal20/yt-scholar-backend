-- +goose Up
CREATE TABLE playlist (
    playlist_id     VARCHAR(255) PRIMARY KEY,
    user_id         VARCHAR NOT NULL REFERENCES users(user_id) ON DELETE CASCADE,
    title           VARCHAR(255) NOT NULL,
    url             TEXT NOT NULL,
    thumbnail_url   TEXT,
    channel         VARCHAR(255),
    videos          JSON,
    updated_at  TIMESTAMP,
    created_at  TIMESTAMP DEFAULT now()
);


-- +goose Down
drop table if exists playlist;
