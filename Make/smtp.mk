
.docker_smtp_up: .init_dot.env .init_dockerfile
	@docker compose -f ${DOCKER_YAML} --env-file dot.env up smtp_* -d --remove-orphans \
	  || { echo "[ERROR] failed to construct docker SMTP service" >&2; exit 1; }

.docker_smtp_rm: .init_dot.env
	@docker compose -f ${DOCKER_YAML} --env-file dot.env ps --services | grep '^smtp_' | \
	 xargs -r docker compose -f ${DOCKER_YAML} --env-file dot.env rm -sf \
	  || { echo "[ERROR] failed to remove docker SMTP service" >&2; exit 1; }

.docker_smtp_stop: .init_dot.env
	@docker compose -f ${DOCKER_YAML} --env-file dot.env ps --services | grep '^smtp_' | \
	 xargs -r docker compose -f ${DOCKER_YAML} --env-file dot.env stop \
	  || { echo "[ERROR] failed to stop docker SMTP service" >&2; exit 1; }
