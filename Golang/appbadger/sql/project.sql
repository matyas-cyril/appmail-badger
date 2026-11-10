-- name: AddProject :one
INSERT INTO Projects (
    name, sender, reply_to, rcpt_max, burst, rate_limit, size_max, rcpt_pattern, enabled
) VALUES (
    sqlc.arg('Name'),
    sqlc.arg('Sender'),
    sqlc.narg('ReplyTo'),
    sqlc.narg('RcptMax'),
    sqlc.narg('Burst'),
    sqlc.narg('RateLimit'),
    sqlc.narg('SizeMax'),
    sqlc.narg('RcptPattern'),
    sqlc.arg('Enable')
)
RETURNING id_project AS uuid, name, sender, reply_to, rcpt_max, burst, rate_limit,size_max,rcpt_pattern, enabled;

-- name: AddDomainInProject :many
WITH req_domains AS (
    SELECT unnest(sqlc.arg('Domains')::text[]) AS dom
),
all_dom_match AS (
    SELECT NOT EXISTS (
        SELECT 1
        FROM req_domains R
        WHERE NOT EXISTS (
            SELECT 1
            FROM Domains D
            WHERE D.domain LIKE R.dom
        )
    ) AS ok
)
INSERT INTO ProjectsDomains(id_project, id_domain)
SELECT P.id_project, D.id_domain
FROM Projects AS P
CROSS JOIN Domains AS D
CROSS JOIN all_dom_match AS M
WHERE P.id_project = sqlc.arg('UUID')
  AND M.ok
  AND EXISTS (
      SELECT 1
      FROM req_domains R
      WHERE D.domain LIKE R.dom
  )
RETURNING id_project AS uuid_project, id_domain;

-- name: AddAuthInProject :many
WITH req_users AS (
    SELECT unnest(sqlc.arg('Users')::text[]) AS username
),
all_user_match AS (
    SELECT NOT EXISTS (
        SELECT 1
        FROM req_users R
        WHERE NOT EXISTS (
            SELECT 1
            FROM Users U
            WHERE U.username LIKE R.username
        )
    ) AS ok
)
INSERT INTO Auth(id_project, id_user)
SELECT P.id_project, U.id_user
FROM Projects AS P
CROSS JOIN Users AS U
CROSS JOIN all_user_match AS M
WHERE P.id_project = sqlc.arg('UUID')
  AND M.ok
  AND EXISTS (
      SELECT 1
      FROM req_users R
      WHERE U.username LIKE R.username
  )
RETURNING id_project AS uuid_project, id_user AS uuid_user;

-- name: AddServerInProject :many
WITH req_ips AS (
    SELECT unnest(sqlc.arg('IPs')::inet[]) AS IP
),
req_names AS (
    SELECT unnest(sqlc.arg('ServerNames')::text[]) AS Name
),
all_servers_match AS (
    SELECT
        NOT EXISTS (
            SELECT 1
            FROM req_ips AS RI
            WHERE NOT EXISTS (
                SELECT 1
                FROM Servers AS S
                WHERE S.ip = RI.ip
            )
        )
        AND NOT EXISTS (
            SELECT 1
            FROM req_names AS RN
            WHERE NOT EXISTS (
                SELECT 1
                FROM Servers AS S
                WHERE S.name = RN.name
            )
        ) AS ok
)
INSERT INTO ProjectsServers (id_project, id_server)
SELECT DISTINCT P.id_project, S.id_server
FROM Projects AS P
CROSS JOIN Servers AS S
CROSS JOIN all_servers_match AS M
WHERE P.id_project = sqlc.arg('UUID_Project')
  AND M.ok
  AND (
      EXISTS (SELECT 1 FROM req_ips AS RI WHERE RI.ip = S.ip)
      OR EXISTS (SELECT 1 FROM req_names AS RN WHERE RN.name = S.name)
  )
RETURNING id_project AS uuid_project, id_server AS uuid_server;
