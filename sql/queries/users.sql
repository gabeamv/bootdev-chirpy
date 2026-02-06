-- name: CreateUser :one
INSERT INTO users (created_at, updated_at, email, hashed_password)
VALUES (
    $1,
    $2,
    $3,
    $4
)
RETURNING *;

-- name: DeleteAllUsers :exec
DELETE FROM users;

-- name: GetUserByEmail :one
SELECT * FROM users
WHERE email = $1;

-- name: GetUserByRefreshToken :one
SELECT users.* FROM users
INNER JOIN refresh_tokens ON users.id = refresh_tokens.user_id
WHERE token = $1;

-- name: UpdateUserEmailPassword :one
UPDATE users
SET email = $2, hashed_password = $3
WHERE id = $1
RETURNING id, created_at, updated_at, email;

-- name: UpdateToChirpyRedByID :exec
UPDATE users
SET is_chirpy_red = TRUE
WHERE id = $1;