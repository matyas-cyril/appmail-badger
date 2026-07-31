#!/usr/bin/env python3

# MTK: 31 juillet 2026

import sys
import json
import os
import re
import argparse
import configparser
import psycopg2
from psycopg2 import pool

# Fichier de conf par défaut
CONFIG_FILE = "/etc/rspamd/local.d/postgres.conf"

def is_valid_email(email):
    """Valide le format syntaxique d'une adresse email"""
    pattern = r'^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$'
    return re.match(pattern, email) is not None

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
    result = {
        "success": False,
        "error": f"Erreur d'initialisation du pool: {e}"
    }
    print(json.dumps(result))
    sys.exit(1)

def get_sender_ratelimit(email):
    """Obtenir le ratelimit pour un sender"""

    conn = None
    try:
        conn = connection_pool.getconn()
        
        with conn.cursor() as cursor:
            cursor.execute("SELECT burst, rate_limit FROM v_enabled_projects WHERE mail = %s LIMIT 2", (email,))
            results = cursor.fetchall()

            match len(results):

                case 1:

                    burst = 0

                    # On cast la burn de retour
                    try:
                        burst = int(str(results[0][0]))

                    except ValueError:
                        return {
                            "success": False,
                            "email": email,
                            "error": "Failed to cast burst value to integer"
                        }

                    # Traitement du ratelimit
                    ratelimit = str(results[0][1]).strip()

                    if ratelimit == "0":
                        ratelimit = 0

                    return {
                        "success": True,
                        "email": email,
                        "burst": burst,
                        "ratelimit": ratelimit
                    }           

                case 2: # Ne doit pas arriver....
                    return {
                        "success": False,
                        "email": email,
                        "error": "Return 2 results instead of 1"
                   }                   

                case _:
                    return {
                        "success": False,
                        "email": email,
                        "error": "No data available"
                    }

    except Exception as e:
        return {
            "success": False,
            "email": email,
            "error": str(e)
        }
        
    finally:
        if conn is not None:
            connection_pool.putconn(conn)

def get_sender_maxsize(email):
    """Obtenir la taille max autorisée pour une message pour un sender"""

    conn = None
    try:
        conn = connection_pool.getconn()
        
        with conn.cursor() as cursor:
            cursor.execute("SELECT size_max FROM v_enabled_projects WHERE mail = %s LIMIT 2", (email,))
            results = cursor.fetchall()

            match len(results):

                case 1:

                    maxsize = 0

                    # On cast la valeur de retour
                    try:
                        maxsize = int(str(results[0][0]))

                    except ValueError:
                        return {
                            "success": False,
                            "email": email,
                            "error": "Failed to cast size_max value to integer"
                        }

                    return {
                        "success": True,
                        "email": email,
                        "maxsize": maxsize
                    }           

                case 2: # Ne doit pas arriver....
                    return {
                        "success": False,
                        "email": email,
                        "error": "Return 2 results instead of 1"
                   }                   

                case _:
                    return {
                        "success": False,
                        "email": email,
                        "message": "No size available"
                    }

    except Exception as e:
        return {
            "success": False,
            "email": email,
            "error": str(e)
        }
        
    finally:
        if conn is not None:
            connection_pool.putconn(conn)

if __name__ == "__main__":

    parser = argparse.ArgumentParser(
        description="Permet les requêtes pour rspamd",
        add_help=True
    )
    
    # Récupérer les paramètres de cmd
    group = parser.add_mutually_exclusive_group(required=True)
    group.add_argument(
        "-ratelimit",
        metavar="<email>",
        help="Obtenir la limite de débit pour un email"
    )
    group.add_argument(
        "-maxsize",
        metavar="<email>",
        help="Obtenir la taille maximale pour un email"
    )
    
    args = parser.parse_args()

    # Pour pouvoir aiguiller.. Y a surement plus propre
    if args.ratelimit:
        email = args.ratelimit
        option = "ratelimit"
    else:
        email = args.maxsize
        option = "maxsize"
    
    # Validation du format email
    if not is_valid_email(email):
        result = {
            "success": False,
            "email": email,
            "error": "Email syntax invalid"
        }
        print(json.dumps(result))
        sys.exit(1)
    
    # Routage vers la bonne requête
    if option == "ratelimit":
        result = get_sender_ratelimit(email)
    else:
        result = get_sender_maxsize(email)
    
    print(json.dumps(result))
    sys.exit(0)
