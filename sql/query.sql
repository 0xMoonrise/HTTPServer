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

-- name: CreateRefreshToken :exec
INSERT INTO token (token, created_at, updated_at, user_id, expire_at, revoke_at)
VALUES (
    $1,
    NOW(),
    NOW(),
	$2,
    $3,
    $4
);

-- name: GetUserByRefreshToken :one
SELECT 
	user_id
FROM token
WHERE token=$1;

-- name: UpdateRevokeToken :exec
UPDATE token
SET token = $1,
	revoke_at  = NOW(),
	updated_at = NOW()
WHERE user_id = $2;

-- name: ChangePassAndEmail :exec
UPDATE users
SET hashed_password = $1,
email = $2
WHERE id = $3;

-- name: DeleteChirp :exec
DELETE
FROM chirp
WHERE id=$1
AND   user_id=$2;
