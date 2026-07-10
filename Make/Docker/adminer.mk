
.docker_adminer_up: .init_dot.env .init_dockerfile
	@docker compose -f ${DOCKER_YAML} --env-file dot.env up adminer -d --remove-orphans \
	  || { echo "[ERROR] failed to construct docker Adminer service" >&2; exit 1; }

.docker_adminer_stop: .init_dot.env
	@docker compose -f ${DOCKER_YAML} --env-file dot.env stop adminer \
	  || { echo "[ERROR] failed to down docker Adminer service" >&2; exit 1; }

.docker_adminer_rm: .init_dot.env
	@docker compose -f ${DOCKER_YAML} --env-file dot.env rm -sf adminer \
	  || { echo "[ERROR] failed to remove docker Adminer service" >&2; exit 1; }

.docker_adminer_purge: .init_dot.env .docker_adminer_rm
	@docker image prune -f -a --filter "label=com.docker.compose.project=${PROJECT_NAME}" \
	  || { echo "[ERROR] failed to remove Adminer image" >&2; exit 1; }
