# Créer un fichier vide pour éviter les if dans les déclarations suivantes
.init_dot.env:
	@if [ ! -f "dot.env" ]; then touch dot.env; fi

# Initialiser le Dockerfile à partir du template
.init_dockerfile:
	@if [ -d "./sources" ] && [ -n "$(COPY_SRC_CMD)" ] && [ -n "$(COPY_GPG_CMD)" ]; then \
		for dir in Docker/DEBIAN; do \
			cp -rf ./sources/ "./$$dir/sources"; \
			sed -e "s|\#\ \%COPY_GPG_CMD\%|$(COPY_GPG_CMD)|g" ./$$dir/Dockerfile.template | sed -e "s|\#\ \%COPY_SRC_CMD\%|$(COPY_SRC_CMD)|g" - > "./Docker/$$dir/Dockerfile"; \
		done \
	else \
		for dir in Docker/DEBIAN; do \
			sed "\|\#\ \%COPY_GPG_CMD\%|d" ./$$dir/Dockerfile.template | sed "\|\#\ \%COPY_SRC_CMD\%|d" - > "./$$dir/Dockerfile"; \
			rm -rf "./Docker/$$dir/sources"; \
		done \
	fi

### PLATEFORM ###
.up_platform: .init_dot.env .init_dockerfile
	@docker compose -f ${DOCKER_YAML} --env-file dot.env up -d --remove-orphans \
	  || { echo "[ERROR] failed to construct docker printers plateform" >&2; exit 1; }

.stop_platform: .init_dot.env
	@docker compose -f ${DOCKER_YAML} --env-file dot.env stop \
	  || { echo "[ERROR] failed to remove plateform - $?" >&2; exit 1; };
	@echo "successfully delete plateform";

.purge_platform: .init_dot.env
	@docker compose -f ${DOCKER_YAML} --env-file dot.env down -v --remove-orphans --rmi all \
	  || { echo "[ERROR] failed to purge plateform - $?" >&2; exit 1; }
	@echo "successfully purge plateform";

.rm_platform: .init_dot.env
	@docker compose -f ${DOCKER_YAML} --env-file dot.env down \
	  || { echo "[ERROR] failed to remove plateform - $?" >&2; exit 1; };
	@echo "successfully delete plateform";
###
