SHELL=/bin/bash
PYTHON=python3
ANSIBLE_VERSION=10.7.0
VENV_DIR=venv

DOCKER_YAML=Docker/appmail-badger.yml

# Inclure uniquement le fichier si existe
sinclude dot.env

-include Make/docker.mk

.PHONY: all 

all: platform-up

purge: .purge_platform

platform-up: .up_platform
platform-stop: .stop_platform
platform-rm: .rm_platform
platform-purge: .purge_platform

postgres-up: .docker_postgres_up
postgres-stop: .docker_postgres_stop
postgres-rm: .docker_postgres_rm
postgres-purge: .docker_postgres_purge

smtp-up: .docker_smtp_up
smtp-stop: .docker_smtp_stop
smtp-rm: .docker_smtp_rm
