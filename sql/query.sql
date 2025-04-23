-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1
)
RETURNING *;

-- name: ExistUser :one

SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);

-- name: WipeUsers :exec
DELETE FROM users;

-- CHIRP

-- name: CreateChirp :one
INSERT INTO chirp (id, created_at, updated_at, body, user_id)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING *;
