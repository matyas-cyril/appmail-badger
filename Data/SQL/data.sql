--- Ajouter des utilisateurs -- Mot de passe à test
INSERT INTO Users (username, password) VALUES ('cyril',  '$argon2id$v=19$m=16,t=2,p=1$TUpyNzVYMmdKYUhFOWZmRQ$+H+PB+1BtS2g8qfC2BG9NQ');
INSERT INTO Users (username, password) VALUES ('samy',   '$argon2id$v=19$m=16,t=2,p=1$eG83WmlmSWFPUVEyV05HTA$kzMu+LbLtK0olKKevvk3ug');
INSERT INTO Users (username, password) VALUES ('benoit', '$argon2id$v=19$m=16,t=2,p=1$MjZvOFhyWkNZZWdPSm43bQ$XSlVJitDiO1pp+K8LmYZDw');

--- Ajout des senders
INSERT INTO Senders (sender) VALUES ('sans-reponse');
INSERT INTO Senders (sender) VALUES ('no-reply');
INSERT INTO Senders (sender) VALUES ('photocopieurs');

--- Ajout de serveurs
INSERT INTO Servers (name, ip) VALUES ('server1.loc', '192.168.10.1');
INSERT INTO Servers (name, ip) VALUES ('server2.loc', '192.168.10.2');
INSERT INTO Servers (name, ip) VALUES ('server3.loc', '192.168.10.3');
INSERT INTO Servers (name, ip) VALUES ('server4.loc', '192.168.10.4');
INSERT INTO Servers (name, ip) VALUES ('server5.loc', '192.168.10.5');

--- Ajout des domaines
INSERT INTO Domains (domain) VALUES ('email.test');
INSERT INTO Domains (domain) VALUES ('programm.test');

--- Ajout des projets
INSERT INTO Projects(name, id_sender) VALUES ('Photocopieurs direction',3);

--- Affecter des utilisateurs à un projet
INSERT INTO Auth (id_project, id_user)
SELECT P.id_project AS id_project, U.id_user AS id_user
FROM Projects AS P
CROSS JOIN Users AS U
WHERE P.name LIKE 'Photocopieurs direction'
  AND ( U.username LIKE 'cyril' OR U.username LIKE 'samy' );

--- Affecter des serveurs à un projet
INSERT INTO ProjectsServers(id_project, id_server)
SELECT P.id_project AS id_project, S.id_server AS id_server
FROM Projects AS P
CROSS JOIN Servers AS S
WHERE P.name LIKE 'Photocopieurs direction'
  AND ( S.ip = '192.168.10.1'
        OR S.name = 'server4.loc'
        OR S.name = 'server2.loc' );

--- Ajouter des domaines à un projet
INSERT INTO ProjectsDomains(id_project, id_domain)
SELECT P.id_project AS id_project, D.id_domain AS id_domain
FROM Projects AS P
CROSS JOIN Domains AS D
WHERE P.name LIKE 'Photocopieurs direction'
  AND ( domain = 'email.test' );

--- NE PAS OUBLIER D'ACTIVER LE PROJET

