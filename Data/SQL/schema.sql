
----- CRÉATION DES TABLES -----

--- Table Groups
CREATE TABLE IF NOT EXISTS Groups (
	id_group BIGSERIAL PRIMARY KEY,
	name TEXT NOT NULL,
	description TEXT
);

--- Table Senders
CREATE TABLE IF NOT EXISTS Senders (
	id_sender BIGSERIAL PRIMARY KEY,
	sender VARCHAR(64) NOT NULL,
	description TEXT
);

--- Table Domains
CREATE TABLE IF NOT EXISTS Domains (
	id_domain BIGSERIAL PRIMARY KEY,
	domain VARCHAR(254) NOT NULL
);
CREATE UNIQUE INDEX IF NOT EXISTS domain ON Domains(domain);

--- Table Servers
CREATE TABLE IF NOT EXISTS Servers (
	id_server UUID PRIMARY KEY DEFAULT uuidv7(),
	name TEXT NOT NULL,
	ip INET NOT NULL,
	description TEXT,
	
	CONSTRAINT servers_name_unique UNIQUE (name),
	CONSTRAINT servers_ip_unique UNIQUE (ip)
);

--- Table Users
CREATE TABLE IF NOT EXISTS Users (
	id_user UUID PRIMARY KEY DEFAULT uuidv7(),
	username VARCHAR(64) NOT NULL,
	password VARCHAR(256) NOT NULL CHECK (password ~ '^\$argon2id\$v=\d+\$m=\d+,t=\d+,p=\d+\$[A-Za-z0-9+/=]+\$[A-Za-z0-9+/=]+$'),
	not_before TIMESTAMPTZ NOT NULL DEFAULT now(),
	expiration TIMESTAMPTZ,
    modified TIMESTAMPTZ NOT NULL DEFAULT now(),
	disabled BOOLEAN NOT NULL DEFAULT false,
	description TEXT,
	
	--On ne peut pas avoir la date d'expiration <= à la date d'utilisation
	CONSTRAINT issued_after_not_before CHECK (expiration IS NULL OR expiration > not_before)
);
CREATE UNIQUE INDEX IF NOT EXISTS users_username_idx ON Users(username);

--- Table Projects
CREATE TABLE IF NOT EXISTS Projects (
	id_project uuid PRIMARY KEY DEFAULT uuidv7(),
	name TEXT NOT NULL,
	id_sender BIGINT NOT NULL, -- S'assurer que chaque Project a exactement 1 Sender
	reply_to VARCHAR(64),
	rcpt_max SMALLINT DEFAULT 15 CHECK (rcpt_max >= 0 AND rcpt_max <= 1500),
    burst INTEGER NOT NULL DEFAULT 150 CHECK (burst >= 0 AND burst <= 65535), -- Nbr de mails que peut envoyer un utilisateur en 1 seul fois - mode rafale ;-)
	rate_limit VARCHAR(20) NOT NULL DEFAULT 0 CHECK (rate_limit ~ '^(?:\d+/\d+(?:s|m|h|d)|0)$'), -- 0 désactive le ratelimit
	size_max NUMERIC(20,0) NOT NULL DEFAULT 3145728 CHECK (size_max >= 0 AND size_max <= 18446744073709551615), -- Correspond à UINT64
	rcpt_pattern TEXT,
	enabled BOOLEAN NOT NULL DEFAULT false,
	description TEXT,
	id_group BIGINT,
	
	CONSTRAINT projects_id_group_fkey FOREIGN KEY (id_group) REFERENCES Groups(id_group) ON DELETE SET NULL
);

--- Table Auth
CREATE TABLE IF NOT EXISTS Auth (
	id_project UUID NOT NULL,
	id_user    UUID NOT NULL,

	PRIMARY KEY (id_project, id_user),
	FOREIGN KEY (id_project) REFERENCES Projects(id_project) ON DELETE CASCADE,
	FOREIGN KEY (id_user) REFERENCES Users(id_user) ON DELETE CASCADE
);

--- Table ProjectServers
CREATE TABLE IF NOT EXISTS ProjectsServers (
	id_project UUID NOT NULL,
	id_server  UUID NOT NULL,

	PRIMARY KEY (id_project, id_server),
	FOREIGN KEY (id_project) REFERENCES Projects(id_project) ON DELETE CASCADE,
	FOREIGN KEY (id_server) REFERENCES Servers(id_server) ON DELETE CASCADE
);

--- Table ProjectsDomains
CREATE TABLE IF NOT EXISTS ProjectsDomains (
	id_project UUID NOT NULL,
	id_domain BIGINT NOT NULL,

	PRIMARY KEY (id_project, id_domain),
	FOREIGN KEY (id_project) REFERENCES Projects(id_project) ON DELETE CASCADE,
	FOREIGN KEY (id_domain) REFERENCES Domains(id_domain) ON DELETE CASCADE
);

----- Fonctions ------
CREATE OR REPLACE FUNCTION enforce_ip_host_format()
RETURNS TRIGGER
LANGUAGE plpgsql AS $$
BEGIN
	-- Convertir l'IP en format d'hôte
	NEW.ip := host(NEW.ip);
	RETURN NEW;
END;
$$;

CREATE OR REPLACE FUNCTION set_users_modified()
RETURNS trigger
LANGUAGE plpgsql AS $$
BEGIN
    NEW.modified := now();
    RETURN NEW;
END;
$$;

----- Triggers -----
CREATE TRIGGER enforce_ip_host
	BEFORE INSERT OR UPDATE ON Servers
	FOR EACH ROW
	EXECUTE FUNCTION enforce_ip_host_format();

CREATE TRIGGER users_set_modified
    BEFORE INSERT OR UPDATE ON users
    FOR EACH ROW
    EXECUTE FUNCTION set_users_modified();

----- View -----

--- Liste des serveurs utilisés dans les projets
CREATE OR REPLACE VIEW v_used_servers AS
    SELECT DISTINCT s.id_server AS id_server, s.ip AS ip, s.name AS name
    FROM Servers s
    JOIN ProjectsServers ps ON ps.id_server = s.id_server
    JOIN Projects p ON p.id_project = ps.id_project
    WHERE p.enabled IS TRUE
    ORDER BY s.ip ASC;

--- Liste des serveurs non utilisés mais appartenant à des projets
CREATE OR REPLACE VIEW v_unused_servers AS
    SELECT DISTINCT s.id_server AS id_server, s.ip AS ip, s.name AS name
    FROM Servers s
    JOIN ProjectsServers ps ON ps.id_server = s.id_server
    JOIN Projects p ON p.id_project = ps.id_project
    WHERE p.enabled IS FALSE
    ORDER BY s.ip ASC;

--- Liste des serveurs sans projets -- Orphelins
CREATE OR REPLACE VIEW v_orphaned_servers AS
    SELECT DISTINCT s.id_server AS id_server, s.ip AS ip, s.name AS name, s.description AS description
    FROM Servers s
    FULL JOIN ProjectsServers ps ON ps.id_server = s.id_server
    WHERE ps.id_server IS NULL OR s.id_server IS NULL
    ORDER BY s.ip ASC;

--- Liste des serveurs authorisés
CREATE OR REPLACE VIEW v_clients_allowlist AS
    SELECT DISTINCT s.ip AS ip, s.name AS servername
    FROM Servers s
    JOIN ProjectsServers ps ON ps.id_server = s.id_server
    JOIN Projects p ON p.id_project = ps.id_project
    WHERE p.enabled IS TRUE;

--- Liste des comptes mails, sender @ domains actifs
CREATE OR REPLACE VIEW v_enabled_projects AS
    SELECT DISTINCT S.sender, D.domain AS domain, CONCAT(S.sender, '@', D.domain) AS mail, P.burst AS burst, P.rate_limit AS rate_limit, P.size_max AS size_max, P.reply_to AS reply_to, P.rcpt_max AS rcpt_max, P.rcpt_pattern AS rcpt_pattern, P.id_project AS id_project
    FROM ProjectsDomains PD
    INNER JOIN Domains D ON D.id_domain = PD.id_domain
    INNER JOIN Projects P ON P.id_project = PD.id_project AND P.enabled IS TRUE
    INNER JOIN Senders S ON S.id_sender = P.id_sender
    ORDER BY S.sender ASC;

--- Liste des Projets par Serveur
CREATE OR REPLACE VIEW v_projects_by_server AS
    SELECT SV.ip AS ip_server, SV.name AS server_name, S.sender AS sender, P.id_project AS id_project, P.enabled AS enabled_project
    FROM ProjectsServers PS
    INNER JOIN Projects P ON P.id_project = PS.id_project
    INNER JOIN Servers SV ON SV.id_server = PS.id_server
    INNER JOIN Senders S ON S.id_sender = P.id_sender
    ORDER BY SV.ip ASC;

--- Liste des utilisateurs autorisés
CREATE OR REPLACE VIEW v_logins AS
    SELECT id_user, username, password AS password_hash
    FROM users
    WHERE disabled = false
        AND not_before < now()
        AND (expiration IS NULL OR expiration > now());

--- Liste pour le login mismatch
CREATE OR REPLACE VIEW v_sender_login_maps AS
    SELECT DISTINCT LOWER(CONCAT(S.sender,'@',D.domain)) AS mail, U.username AS login
    FROM Auth A
    INNER JOIN Projects P ON P.id_project = A.id_project
    INNER JOIN Users U ON U.id_user = A.id_user
    INNER JOIN Senders S ON S.id_sender = P.id_sender
    INNER JOIN ProjectsDomains PD ON PD.id_project = P.id_project
    INNER JOIN Domains D ON D.id_domain = PD.id_domain
    WHERE P.enabled = true;

--- Liste des logins non actifs
CREATE OR REPLACE VIEW v_logins_inactive AS
    SELECT id_user, username, not_before, expiration, disabled
    FROM Users
    WHERE disabled IS TRUE OR
          expiration < NOW() OR
          not_before > NOW()
    ORDER BY username;

