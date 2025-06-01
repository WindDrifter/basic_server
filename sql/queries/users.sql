-- name: CreateUser :one
INSERT INTO users (email)
VALUES (
    $1
)
RETURNING *;

-- name: GetUserByEmail :one
SELECT * from users where email = $1;
-- name: GetAllUsers :many
SELECT * from users;
-- name: DeleteAllUsers :exec
DELETE from users;