DROP VIEW IF EXISTS v_clients_allowlist;
DROP VIEW IF EXISTS v_enabled_projects;
DROP VIEW IF EXISTS v_orphaned_servers;
DROP VIEW IF EXISTS v_projects_by_server;
DROP VIEW IF EXISTS v_unused_servers;
DROP VIEW IF EXISTS v_used_servers;
DROP VIEW IF EXISTS v_logins;
DROP VIEW IF EXISTS v_sender_login_maps;

DROP TABLE IF EXISTS ProjectsDomains;
DROP TABLE IF EXISTS ProjectsServers;
DROP TABLE IF EXISTS Servers;
DROP TABLE IF EXISTS Senders;
DROP TABLE IF EXISTS Auth;
DROP TABLE IF EXISTS Domains;
DROP TABLE IF EXISTS Projects;
DROP TABLE IF EXISTS Users;
DROP TABLE IF EXISTS Groups;

DROP FUNCTION IF EXISTS enforce_ip_host_format;

DROP TRIGGER IF EXISTS users_set_modified ON Users;
