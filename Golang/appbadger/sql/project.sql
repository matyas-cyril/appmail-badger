-- name: AddProject :one
INSERT INTO Projects (
    name, sender, enabled
) VALUES (
    sqlc.arg('Name'),
    sqlc.narg('Sender'),
    sqlc.arg('Enable')
)
RETURNING id_project AS uuid, name, sender, reply_to, rcpt_max, burst, rate_limit,size_max,rcpt_pattern, enabled;

-- name: AddDomainInProject :many
INSERT INTO ProjectsDomains(id_project, id_domain)
SELECT P.id_project AS id_project, D.id_domain AS id_domain
FROM Projects AS P
CROSS JOIN Domains AS D
WHERE P.id_project = sqlc.arg('UUID')
  AND D.domain = ANY(sqlc.arg('Domains')::text[])
RETURNING *;
