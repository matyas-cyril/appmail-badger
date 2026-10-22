-- name: AddServer :one
INSERT INTO Servers (
    name, ip, description
) VALUES (
    lower(trim(sqlc.arg('Name'))),
    sqlc.arg('IP'),
    sqlc.narg('Description')
)
RETURNING id_server AS uuid, name, ip , description;
