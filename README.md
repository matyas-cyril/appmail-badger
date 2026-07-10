# appmail-badger

# 1. Make

## 1.1 Ensemble de la plateforme

Les commandes suivantes permettent d'interagir avec l'ensemble des services.  

```bash
# Générer la plateforme
make platform-up

# Mettre en pause la plateforme
make platform-stop

# Arrêter et supprimer les containers
make platform-rm

# Arrêter et supprimer les images, même les orphelines
make platform-purge
ou
make purge
```
## 1.2 Service PostgreSQL

```bash
# Générer le service Postgres
make postgres-up

# Mettre en pause le service Postgres
make postgres-stop

# Arrêter et supprimer le container
make postgres-rm

# Arrêter et supprimer le container, le volume et l'image du service Postgres
make postgres-purge
```

## 1.3 Service SMTP

```bash
# Générer le ou les services SMTP
make smtp-up

# Mettre en pause le ou les services SMTP
make smtp-stop

# Arrêter et supprimer le ou les containers SMTP
make smtp-rm

# Arrêter et supprimer le container et l'image du service SMTP
make smtp-purge
```

## 1.4 Service Cyrus

```bash
# Générer le ou les services Cyrus
make cyrus-up

# Mettre en pause le ou les services Cyrus
make cyrus-stop

# Arrêter et supprimer le ou les containers Cyrus
make cyrus-rm

# Arrêter et supprimer le container et l'image du service Cyrus
make cyrus-purge
```

## 1.4 Service Redis

```bash
# Générer le ou les services Redis
make redis-up

# Mettre en pause le ou les services Redis
make redis-stop

# Arrêter et supprimer le ou les containers Redis
make redis-rm

# Arrêter et supprimer le container et l'image du service Redis
make redis-purge
```

## 1.5 Service Webmail

```bash
# Générer le ou les services Webmail
make webmail-up

# Mettre en pause le ou les services Webmail
make webmail-stop

# Arrêter et supprimer le ou les containers Webmail
make webmail-rm

# Arrêter et supprimer le container et l'image du service Webmail
make webmail-purge
```

## 1.6 Service Adminer

Le service Adminer dépend du service Postgres.  
Un **make adminer-up** va créer le service Postgres si non présent.  

```bash
# Générer le ou les services Adminer
make adminer-up

# Mettre en pause le ou les services Adminer
make adminer-stop

# Arrêter et supprimer le ou les containers Adminer
make adminer-rm

# Arrêter et supprimer le container et l'image du service Adminer
make adminer-purge
```

## 1.7 Base de données

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