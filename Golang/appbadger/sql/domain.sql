-- name: AddDomain :one
INSERT INTO Domains (
   domain
) VALUES (
    lower(trim(sqlc.arg('Domain')))
)
RETURNING domain AS Dom;
