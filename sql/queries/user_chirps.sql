-- name: CreateUserChirp :one
INSERT INTO user_chirps (created_at, updated_at, body, user_id)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: GetAllChirps :many
SELECT * FROM user_chirps
ORDER BY created_at;

-- name: GetChirp :one
SELECT * FROM user_chirps
WHERE id = $1;