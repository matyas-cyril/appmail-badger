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

platform-up: .init_dot.env .up_platform
platform-stop: .stop_platform
platform-rm: .rm_platform
platform-purge: .purge_platform

