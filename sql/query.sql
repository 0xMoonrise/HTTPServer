-- name: CreateUser :one
INSERT INTO users (id, created_at, updated_at, email, hashed_password)
VALUES (
    gen_random_uuid(),
    NOW(),
    NOW(),
    $1,
    $2
)
RETURNING id, created_at, updated_at, email;

-- name: ExistUser :one
SELECT EXISTS(SELECT 1 FROM users WHERE email = $1);

-- name: ExistUserById :one
SELECT EXISTS(SELECT 1 FROM users WHERE id = $1);

-- name: ExistChirpById :one
SELECT EXISTS(SELECT 1 FROM chirp WHERE id = $1);

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

-- name: GetChirps :many
SELECT * FROM chirp;

-- name: GetChirpById :one
SELECT 
	id, 
	created_at, 
	updated_at, 
	body, 
	user_id 
FROM chirp 
WHERE id = $1; 

-- name: GetUserPassword :one

SELECT
	hashed_password
FROM users
WHERE email = $1;

-- name: GetUserByEmail :one
SELECT
	id,
	created_at,
	updated_at,
	email
FROM users
WHERE email = $1;
