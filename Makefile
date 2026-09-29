SHELL=/bin/bash
PYTHON=python3
ANSIBLE_VERSION=10.7.0
VENV_DIR=venv

DOCKER_YAML=Docker/appmail-badger.yml

# Inclure uniquement le fichier si existe
sinclude dot.env

# Attribuer une valeur par défaut
PROJECT_NAME ?= appmail-badger

PG_USER ?= Badger

BDD ?= AppMail
SQL_INIT ?= Data/SQL/init.sql
SQL_DROP ?= Data/SQL/drop.sql
SQL_DATA ?= Data/SQL/data.sql

-include Make/ansible.mk
-include Make/docker.mk
-include Make/sample.mk
-include Make/golang.mk

.PHONY: all

all: ansible-init docker-platform-up sql-init smtp-full-install

demo: all sql-data

purge: .purge_platform ansible-delete
