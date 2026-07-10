
.docker_webmail_up: .init_dot.env .init_dockerfile
	@docker compose -f ${DOCKER_YAML} --env-file dot.env up webmail -d --remove-orphans \
	  || { echo "[ERROR] failed to construct docker Webmail service" >&2; exit 1; }

.docker_webmail_stop: .init_dot.env
	@docker compose -f ${DOCKER_YAML} --env-file dot.env stop webmail \
	  || { echo "[ERROR] failed to down docker Webmail service" >&2; exit 1; }

.docker_webmail_rm: .init_dot.env
	@docker compose -f ${DOCKER_YAML} --env-file dot.env rm -sf webmail \
	  || { echo "[ERROR] failed to remove docker Webmail service" >&2; exit 1; }

.docker_webmail_purge: .init_dot.env .docker_webmail_rm
	@docker image prune -f -a --filter "label=com.docker.compose.project=${PROJECT_NAME}" \
	  || { echo "[ERROR] failed to remove Webmail image" >&2; exit 1; }
