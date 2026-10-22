
-- name: AddLogin :one
INSERT INTO Users (
    username, password, not_before, expiration, disabled, description
) VALUES (
    sqlc.arg('Login'),
    sqlc.arg('Password'),
    sqlc.narg('NotBefore'),
    sqlc.narg('Expiration'),
    sqlc.narg('Disabled'),
    sqlc.narg('Description')
)
RETURNING id_user AS uuid, username AS login, not_before AS nbf, expiration AS exp, disabled, description;

-- name: GetLogin :one
SELECT id_user AS uuid, username AS login, 
    COALESCE(FLOOR(EXTRACT(EPOCH FROM not_before))::BIGINT, 0) AS nbf, 
    COALESCE(FLOOR(EXTRACT(EPOCH FROM expiration))::BIGINT, 0) AS exp, 
    disabled, 
    description
FROM Users
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
FROM Users
ORDER BY username;

-- name: DeleteLogin :one
DELETE FROM Users
WHERE username = $1
RETURNING id_user AS uuid, username AS login;

-- name: DeleteLoginByUUID :one
DELETE FROM Users
WHERE id_user = $1
RETURNING id_user AS uuid, username AS login;

-- name: GetPassword :one
SELECT password_hash
FROM v_logins
WHERE username = $1
LIMIT 1;

-- name: EnableLogin :one
UPDATE Users
SET disabled = '0'
WHERE username = $1
RETURNING id_user AS uuid, username AS login, disabled;

-- name: EnableLoginByUUID :one
UPDATE Users
SET disabled = '0'
WHERE id_user = $1
RETURNING id_user AS uuid, username AS login, disabled;

-- name: DisableLogin :one
UPDATE Users
SET disabled = '1' 
WHERE username = $1
RETURNING id_user AS uuid, username AS login, disabled;

-- name: DisableLoginByUUID :one
UPDATE Users
SET disabled = '1' 
WHERE id_user = $1
RETURNING id_user AS uuid, username AS login, disabled;

-- name: ModifyLogin :one
UPDATE Users
SET
    password = COALESCE(sqlc.narg('Password'), password),
    not_before = COALESCE(sqlc.narg('NotBefore'), not_before),
    expiration = COALESCE(sqlc.narg('Expiration'), expiration),
    disabled = COALESCE(sqlc.narg('Disable'), disabled),
    description = COALESCE(sqlc.narg('Description'), description)
WHERE username = sqlc.arg('Login') OR id_user = sqlc.arg('UUID')
RETURNING id_user AS uuid, username AS login, not_before AS nbf, expiration AS exp, modified AS mod, disabled, description;

-- name: SetExpiration :one
UPDATE Users
SET
    expiration = COALESCE(sqlc.narg('Expiration'), expiration)
WHERE username = sqlc.arg('Login') OR id_user = sqlc.arg('UUID')
RETURNING id_user AS uuid, username AS login, expiration AS exp;

-- name: ClearExpiration :one
UPDATE Users
SET 
    expiration = NULL
WHERE username = sqlc.arg('Login') OR id_user = sqlc.arg('UUID')
RETURNING id_user AS uuid, username AS login, expiration AS exp;
