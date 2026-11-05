-- name: AddServer :one
INSERT INTO Servers (
    name, ip, description
) VALUES (
    lower(trim(sqlc.arg('Name'))),
    sqlc.arg('IP'),
    sqlc.narg('Description')
)
RETURNING id_server AS uuid, name, ip , description;

-- name: DelServer :one
DELETE FROM Servers
WHERE id_server = sqlc.arg('UUID') OR
      name = sqlc.arg('Name') OR
      ip = sqlc.arg('IP')
RETURNING id_server AS uuid, name, ip;

-- name: GetServers :many
SELECT id_server AS uuid, name, ip, description
FROM Servers
ORDER BY ip;

-- name: GetServer :one
SELECT id_server AS uuid, name, ip, description
FROM Servers
WHERE id_server = sqlc.arg('UUID') OR
      name = sqlc.arg('Name') OR
      ip = sqlc.arg('IP')
LIMIT 1;

-- name: ModifyServer :one
UPDATE Servers
SET
    name = COALESCE(lower(trim(sqlc.arg('Name'))), name),
    ip = COALESCE(sqlc.narg('IP'), ip),
    description = COALESCE(sqlc.narg('Description'), description)
WHERE id_server = sqlc.arg('UUID')
RETURNING id_server AS uuid, name, ip, description;
