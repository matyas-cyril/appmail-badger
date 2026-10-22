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

-- name: GetDomain :one
SELECT domain AS Dom
FROM Domains
WHERE domain = lower(trim(sqlc.arg('Domain')))
LIMIT 1;

-- name: GetDomains :many
SELECT domain AS Dom
FROM Domains
ORDER BY domain ASC;
