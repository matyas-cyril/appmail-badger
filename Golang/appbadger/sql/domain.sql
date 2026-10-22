-- name: AddDomain :one
INSERT INTO Domains (
   domain
) VALUES (
    lower(trim(sqlc.arg('Domain')))
)
RETURNING domain AS Dom;

-- name: DelDomain :one
DELETE FROM Domains
WHERE domain = lower(trim(sqlc.arg('Domain')))
RETURNING domain AS Dom;
