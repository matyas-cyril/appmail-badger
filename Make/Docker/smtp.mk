
.docker_smtp_up: .init_dot.env .init_dockerfile
	@docker compose -f ${DOCKER_YAML} --env-file dot.env up smtp_* -d --remove-orphans \
	  || { echo "[ERROR] failed to construct docker SMTP service" >&2; exit 1; }

.docker_smtp_stop: .init_dot.env
	@docker compose -f ${DOCKER_YAML} --env-file dot.env ps --services | grep '^smtp_' | \
	 xargs -r docker compose -f ${DOCKER_YAML} --env-file dot.env stop \
	  || { echo "[ERROR] failed to stop docker SMTP service" >&2; exit 1; }

.docker_smtp_rm: .init_dot.env
	@docker compose -f ${DOCKER_YAML} --env-file dot.env ps --services | grep '^smtp_' | \
	 xargs -r docker compose -f ${DOCKER_YAML} --env-file dot.env rm -sf \
	  || { echo "[ERROR] failed to remove docker SMTP service" >&2; exit 1; }

.docker_smtp_purge: .init_dot.env .docker_smtp_rm
	@docker images --filter=reference='${PROJECT_NAME}/smtp:*' --format="{{.Repository}}:{{.Tag}}" | \
		while read -r image; do \
			docker rmi -f "$$image" >/dev/null 2>&1 || \
				{ echo "[ERROR] failed to delete image $$image" >&2; exit 1; }; \
		done
