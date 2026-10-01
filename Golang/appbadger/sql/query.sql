
-- name: GetLogin :one
SELECT *
FROM users
WHERE username = $1
LIMIT 1;

-- name: GetLoginByUUID :one
SELECT *
FROM users
WHERE id_user = $1
LIMIT 1;