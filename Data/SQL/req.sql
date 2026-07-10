--- Obtenir les Servers actifs uniquement
SELECT DISTINCT s.ip AS ip, s.name AS name
FROM Servers s
JOIN ProjectsServers ps ON ps.id_server = s.id_server
JOIN Projects p ON p.id_project = ps.id_project
WHERE p.enabled IS TRUE
ORDER BY s.ip;



