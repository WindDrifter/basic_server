-- name: CreateChirp :one
INSERT INTO chirps (body, user_id)
VALUES (
    $1,
    $2
)
RETURNING *;

-- name: GetChirpById :one
SELECT * from chirps where id = $1;
-- name: GetAllChirpsByUser :many
SELECT * from chirps where user_id = $1 order by created_at ASC;
-- name: GetAllChirps :many
SELECT * from chirps order by created_at ASC;
-- name: DeleteAllChirps :exec
DELETE from chirps;