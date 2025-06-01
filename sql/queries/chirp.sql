-- name: CreateChirp :one
INSERT INTO chirps (body, user_id)
VALUES (
    $1,
    $2
)
RETURNING *;

-- name: GetChirpById :one
SELECT * from chirps where id = $1;
-- name: GetAllChirpByUser :many
SELECT * from chirps where user_id = $1;
-- name: DeleteAllChirps :exec
DELETE from chirps;