-- name: AddNewPlaylist :exec
INSERT INTO playlist (
    playlist_id, user_id, title, url, thumbnail_url, channel, videos, updated_at
) VALUES (
    $1, $2, $3, $4, $5, $6, $7, now()
);

-- name: GetPlaylistById :one
SELECT * FROM playlist WHERE playlist_id = $1;

-- name: GetAllUserPlaylists :many
SELECT * FROM playlist WHERE user_id = $1;

-- name: GetPlaylistForUser :one
SELECT * FROM playlist WHERE playlist_id = $1 and user_id = $2;


-- name: UpdatePlaylistForUser :exec
UPDATE playlist
SET
    title = $2,
    url = $3,
    thumbnail_url = $4,
    channel = $5,
    videos = $6,
    updated_at = now()
WHERE playlist_id = $1 and user_id = $7;


-- name: DeletePlaylistForUser :exec
delete from playlist where playlist_id = $1 and user_id = $2;
