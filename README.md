# appmail-badger

# 1. Make

## 1.1 Commandes de base

``` bash
# Générer la plateforme, puis déployer les recettes Ansible et le schéma de BDD
make all

# Correspond à make all + ajout d'un jeu de données
make demo

# Arrêter et supprimer les images, même les orphelines (correspond à make docker-plateform-purge)
make purge
```

## 1.2 Déploiement Docker

### 1.2.1 Ensemble de la plateforme


```bash
# Générer la plateforme
make docker-platform-up

# Mettre en pause la plateforme
make docker-platform-stop

# Arrêter et supprimer les containers
make docker-platform-rm

# Arrêter et supprimer les images, même les orphelines (correspond à make purge)
make docker-platform-purge

```
### 1.2.2 Service PostgreSQL

```bash
# Générer le service Postgres
make docker-postgres-up

# Mettre en pause le service Postgres
make docker-postgres-stop

# Arrêter et supprimer le container
make docker-postgres-rm

# Arrêter et supprimer le container, le volume et l'image du service Postgres
make docker-postgres-purge
```

### 1.2.3 Service SMTP

```bash
# Générer le ou les services SMTP
make docker-smtp-up

# Mettre en pause le ou les services SMTP
make docker-smtp-stop

# Arrêter et supprimer le ou les containers SMTP
make docker-smtp-rm

# Arrêter et supprimer le container et l'image du service SMTP
make docker-smtp-purge
```

### 1.2.4 Service Cyrus

```bash
# Générer le ou les services Cyrus
make docker-cyrus-up

# Mettre en pause le ou les services Cyrus
make docker-cyrus-stop

# Arrêter et supprimer le ou les containers Cyrus
make docker-cyrus-rm

# Arrêter et supprimer le container et l'image du service Cyrus
make docker-cyrus-purge
```

### 1.2.5 Service Redis

```bash
# Générer le ou les services Redis
make docker-redis-up

# Mettre en pause le ou les services Redis
make docker-redis-stop

# Arrêter et supprimer le ou les containers Redis
make docker-redis-rm

# Arrêter et supprimer le container et l'image du service Redis
make docker-redis-purge
```

### 1.2.6 Service Webmail

```bash
# Générer le ou les services Webmail
make docker-webmail-up

# Mettre en pause le ou les services Webmail
make docker-webmail-stop

# Arrêter et supprimer le ou les containers Webmail
make docker-webmail-rm

# Arrêter et supprimer le container et l'image du service Webmail
make docker-webmail-purge
```

### 1.2.7 Service Adminer

Le service Adminer dépend du service Postgres.  
Un **make adminer-up** va créer le service Postgres si non présent.  

```bash
# Générer le ou les services Adminer
make docker-adminer-up

# Mettre en pause le ou les services Adminer
make docker-adminer-stop

# Arrêter et supprimer le ou les containers Adminer
make docker-adminer-rm

# Arrêter et supprimer le container et l'image du service Adminer
make docker-adminer-purge
```

## 1.3 Base de données

L'initialisation de la base de données nécessite la présence du service Postgres à minima.  
Pour une exploitation graphique il faut ajouter le service Adminer.  

```bash
# Intialiser la BDD
make sql-init

# Supprimer les tables, views, données...
make sql-drop

# Ajouter des données 
make sql-data
```
