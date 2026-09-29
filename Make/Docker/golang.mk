
.docker_golang_up: .init_dot.env .init_dockerfile
	@docker compose -f ${DOCKER_YAML} --env-file dot.env up golang -d --remove-orphans \
	  || { echo "[ERROR] failed to construct docker Adminer service" >&2; exit 1; }

.docker_golang_stop: .init_dot.env
	@docker compose -f ${DOCKER_YAML} --env-file dot.env stop golang \
	  || { echo "[ERROR] failed to down docker Adminer service" >&2; exit 1; }

.docker_golang_rm: .init_dot.env
	@docker compose -f ${DOCKER_YAML} --env-file dot.env rm -sf golang \
	  || { echo "[ERROR] failed to remove docker Adminer service" >&2; exit 1; }

.docker_golang_purge: .init_dot.env .docker_golang_rm
	@docker image prune -f -a --filter "label=com.docker.compose.project=${PROJECT_NAME}" \
	  || { echo "[ERROR] failed to remove Adminer image" >&2; exit 1; }
