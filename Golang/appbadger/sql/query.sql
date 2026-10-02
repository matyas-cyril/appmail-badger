
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
