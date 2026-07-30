#!/usr/bin/env python3

import sys
import json
import psycopg2
from psycopg2 import pool


try:
    connection_pool = psycopg2.pool.SimpleConnectionPool(
        1, 5,
        dbname="AppMail",
        user="Badger",
        password="Password",
        host="postgres",
        port=5432,
        connect_timeout=3,  # Échoue si la connexion prend plus de 3 secondes
        options="-c statement_timeout=2000"  # Échoue si la requête SQL prend plus de 2000ms (2s)
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
