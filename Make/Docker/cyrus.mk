
.docker_cyrus_up: .init_dot.env .init_dockerfile
	@docker compose -f ${DOCKER_YAML} --env-file dot.env up cyrus -d --remove-orphans \
	  || { echo "[ERROR] failed to construct docker cyrus service" >&2; exit 1; }

.docker_cyrus_stop: .init_dot.env
	@docker compose -f ${DOCKER_YAML} --env-file dot.env stop cyrus \
	  || { echo "[ERROR] failed to down docker Cyrus service" >&2; exit 1; }

.docker_cyrus_rm: .init_dot.env
	@docker compose -f ${DOCKER_YAML} --env-file dot.env rm -sf cyrus \
	  || { echo "[ERROR] failed to remove docker Cyrus service" >&2; exit 1; }

.docker_cyrus_purge: .init_dot.env .docker_cyrus_rm
	@docker images --filter=reference='${PROJECT_NAME}/cyrus:*' --format="{{.Repository}}:{{.Tag}}" | \
		while read -r image; do \
			docker rmi -f "$$image" >/dev/null 2>&1 || \
				{ echo "[ERROR] failed to delete image $$image" >&2; exit 1; }; \
		done
