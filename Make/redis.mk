.docker_redis_up: .init_dot.env .init_dockerfile
	@docker compose -f ${DOCKER_YAML} --env-file dot.env up redis -d --remove-orphans \
	  || { echo "[ERROR] failed to construct docker Redis service" >&2; exit 1; }

.docker_redis_stop: .init_dot.env
	@docker compose -f ${DOCKER_YAML} --env-file dot.env stop redis \
	  || { echo "[ERROR] failed to down docker Redis service" >&2; exit 1; }

.docker_redis_rm: .init_dot.env
	@docker compose -f ${DOCKER_YAML} --env-file dot.env rm -sf redis \
	  || { echo "[ERROR] failed to remove docker Redis service" >&2; exit 1; }

.docker_redis_purge: .init_dot.env .docker_redis_rm
	@docker image prune -f -a --filter "label=com.docker.compose.project=${PROJECT_NAME}" \
	  || { echo "[ERROR] failed to remove Redis image" >&2; exit 1; }
