-- name: CreateUserChirp :one
INSERT INTO user_chirps (created_at, updated_at, body, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;