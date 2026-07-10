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