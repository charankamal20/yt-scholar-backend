-- name: CreateUser :one
INSERT INTO users (user_id, email, name, profile_pic, updated_at, created_at)
VALUES ($1, $2, $3, $4, NOW(), NOW())
RETURNING *;


-- name: GetUserByEmail :one
SELECT * FROM users WHERE email = $1;


-- name: GetUserById :one
SELECT * FROM users WHERE user_id = $1;
