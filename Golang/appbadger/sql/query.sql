
-- name: GetLogin :one
SELECT id_user AS uuid, username AS login, 
    COALESCE(FLOOR(EXTRACT(EPOCH FROM not_before))::BIGINT, 0) AS nbf, 
    COALESCE(FLOOR(EXTRACT(EPOCH FROM expiration))::BIGINT, 0) AS exp, 
    disabled, 
    description
FROM users
WHERE username = $1
LIMIT 1;

-- name: GetLoginByUUID :one
SELECT id_user AS uuid, username AS login, 
    COALESCE(FLOOR(EXTRACT(EPOCH FROM not_before))::BIGINT, 0) AS nbf, 
    COALESCE(FLOOR(EXTRACT(EPOCH FROM expiration))::BIGINT, 0) AS exp, 
    disabled, 
    description
FROM users
WHERE id_user = $1
LIMIT 1;

-- name: GetLogins :many
SELECT id_user AS uuid, username AS login, 
    COALESCE(FLOOR(EXTRACT(EPOCH FROM not_before))::BIGINT, 0) AS nbf, 
    COALESCE(FLOOR(EXTRACT(EPOCH FROM expiration))::BIGINT, 0) AS exp, 
    disabled, 
    description
FROM users
ORDER BY username;

-- name: AddLogin :one
INSERT INTO users (
    username, password, not_before, expiration, disabled, description
) VALUES (
    $1, $2, $3, $4, $5, $6
)
RETURNING *;

-- name: DeleteLogin :one
DELETE FROM users
WHERE username = $1
RETURNING id_user, username;

-- name: DeleteLoginByUUID :one
DELETE FROM users
WHERE id_user = $1
RETURNING id_user, username;

-- name: GetPassword :one
SELECT password_hash
FROM v_logins
WHERE username = $1
LIMIT 1;

-- name: EnableLogin :one
UPDATE users
SET disabled = '0'
WHERE username = $1
RETURNING id_user, username, disabled;

-- name: EnableLoginByUUID :one
UPDATE users
SET disabled = '0'
WHERE id_user = $1
RETURNING id_user, username, disabled;

-- name: DisableLogin :one
UPDATE users
SET disabled = '1' 
WHERE username = $1
RETURNING id_user, username, disabled;

-- name: DisableLoginByUUID :one
UPDATE users
SET disabled = '1' 
WHERE id_user = $1
RETURNING id_user, username, disabled;

-- name: ModifyLogin :one
UPDATE users
SET
    password = COALESCE(sqlc.narg('Password'), password),
    not_before = COALESCE(sqlc.narg('NotBefore'), not_before),
    expiration = COALESCE(sqlc.narg('Expiration'), expiration),
    description = COALESCE(sqlc.narg('Description'), description)
WHERE username = sqlc.arg('Login') OR id_user = sqlc.arg('UUID')
RETURNING *;
