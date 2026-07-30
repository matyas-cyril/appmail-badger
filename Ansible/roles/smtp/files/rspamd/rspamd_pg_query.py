#!/usr/bin/env python3

import sys
import json
import os
import configparser
import psycopg2
from psycopg2 import pool

# Fichier de conf par défaut
CONFIG_FILE = "/etc/rspamd/local.d/postgres.conf"

def load_db_config(filename=CONFIG_FILE):
    """Charge les identifiants depuis un fichier clé = valeur sans section"""
    if not os.path.exists(filename):
        raise FileNotFoundError(f"Le fichier de configuration {filename} est introuvable.")
        
    config = configparser.ConfigParser()
    
    # Astuce : On injecte une section virtuelle [global] pour lire
    with open(filename, "r", encoding="utf-8") as f:
        config_string = "[global]\n" + f.read()
        
    config.read_string(config_string)
    return dict(config.items("global"))

try:
    # Chargement des identifiants
    db_params = load_db_config()

    # Initialisation du pool
    connection_pool = psycopg2.pool.SimpleConnectionPool(
        1, 5,
        dbname=db_params.get("dbname",""),
        user=db_params.get("user",""),
        password=db_params.get("password",""),
        host=db_params.get("host","127.0.0.1"),
        port=int(db_params.get("port", 5432)),
        connect_timeout=int(db_params.get("timout", 3)),
        options="-c statement_timeout=2000"
    )

except Exception as e:
    print(json.dumps({"error": f"Erreur d'initialisation du pool: {e}"}), file=sys.stderr)
    sys.exit(1)

def get_sender_limit(email):
    """Récupère la limite de taille pour un expéditeur de manière sécurisée"""
    conn = None
    try:
        # Récupération de la connexion depuis le pool
        conn = connection_pool.getconn()
        
        # Le bloc 'with' gère automatiquement la fermeture du curseur et les transactions
        with conn.cursor() as cursor:
            cursor.execute("SELECT size_max FROM v_enabled_projects WHERE mail = %s LIMIT 1",(email,))
            result = cursor.fetchone()
            return result[0] if result else None

    except Exception as e:
        print(json.dumps({"error": str(e)}), file=sys.stderr)
        return None
        
    finally:
        # S'assure que la connexion est TOUJOURS renvoyée au pool, même en cas d'erreur
        if conn is not None:
            connection_pool.putconn(conn)

if __name__ == "__main__":

    if len(sys.argv) < 2:
        print("Usage: rspamd_pg_query.py <email>", file=sys.stderr)
        sys.exit(1)
    
    email = sys.argv[1]
    limit = get_sender_limit(email)
    
    if limit is not None:
        print(limit)
    else:
        print("-1")  # -1 = pas de limite trouvée

    sys.exit(0)
