-- name: CreateUser :one
INSERT INTO users (email, password_hash)
VALUES (
    $1,
    $2
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * from users where email = $1;
-- name: GetAllUsers :many
SELECT * from users;
-- name: DeleteAllUsers :exec
DELETE from users;