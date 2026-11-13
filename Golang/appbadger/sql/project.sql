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

-- name: DelProjet :one
DELETE FROM Projects
WHERE id_project = sqlc.arg('UUID')
RETURNING id_project AS uuid, name, sender, enabled;

-- name: GetProject :many
SELECT
    P.id_project AS uuid, P.name AS name, P.sender AS sender,
    P.reply_to AS reply_to, P.rcpt_max AS rcpt_max, P.burst AS burst,
    P.rate_limit AS rate_limit, P.size_max AS size_max, P.rcpt_pattern AS rcpt_pattern,
    P.enabled AS enable,
    ARRAY(
        SELECT S.ip
        FROM ProjectsServers AS PS
        JOIN Servers AS S ON S.id_server = PS.id_server
        WHERE PS.id_project = P.id_project
        ORDER BY S.ip
    ) AS ips,
    ARRAY(
        SELECT S.name
        FROM ProjectsServers AS PS
        JOIN Servers AS S ON S.id_server = PS.id_server
        WHERE PS.id_project = P.id_project
        ORDER BY S.ip
    ) AS name_servers,
    ARRAY(
        SELECT D.domain
        FROM ProjectsDomains AS PD
        JOIN Domains AS D ON D.id_domain = PD.id_domain
        WHERE PD.id_project = P.id_project
        ORDER BY D.id_domain
    ) AS domains
FROM Projects AS P
WHERE P.id_project = sqlc.arg('UUID') OR
      P.name = sqlc.arg('Name') OR
      P.sender = sqlc.arg('Sender') OR
      P.reply_to = sqlc.arg('ReplyTo')
GROUP BY P.id_project;

-- name: GetProjects :many
SELECT
    P.id_project AS uuid, P.name AS name, P.sender AS sender,
    P.reply_to AS reply_to, P.rcpt_max AS rcpt_max, P.burst AS burst,
    P.rate_limit AS rate_limit, P.size_max AS size_max, P.rcpt_pattern AS rcpt_pattern,
    P.enabled AS enable,
    ARRAY(
        SELECT S.ip
        FROM ProjectsServers AS PS
        JOIN Servers AS S ON S.id_server = PS.id_server
        WHERE PS.id_project = P.id_project
        ORDER BY S.ip
    ) AS ips,
    ARRAY(
        SELECT S.name
        FROM ProjectsServers AS PS
        JOIN Servers AS S ON S.id_server = PS.id_server
        WHERE PS.id_project = P.id_project
        ORDER BY S.ip
    ) AS name_servers,
    ARRAY(
        SELECT D.domain
        FROM ProjectsDomains AS PD
        JOIN Domains AS D ON D.id_domain = PD.id_domain
        WHERE PD.id_project = P.id_project
        ORDER BY D.id_domain
    ) AS domains
FROM Projects AS P
GROUP BY P.id_project;
