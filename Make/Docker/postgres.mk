.docker_postgres_up: .init_dot.env .init_dockerfile
	@docker compose -f ${DOCKER_YAML} --env-file dot.env up postgres -d --remove-orphans \
	  || { echo "[ERROR] failed to construct docker Postgres service" >&2; exit 1; }

.docker_postgres_stop: .init_dot.env
	@docker compose -f ${DOCKER_YAML} --env-file dot.env stop postgres \
	  || { echo "[ERROR] failed to down docker Postgres service" >&2; exit 1; }

.docker_postgres_rm: .init_dot.env
	@docker compose -f ${DOCKER_YAML} --env-file dot.env rm -sf postgres \
	  || { echo "[ERROR] failed to remove docker Postgres service" >&2; exit 1; }

.docker_postgres_purge: .init_dot.env .docker_postgres_rm
	@docker volume inspect appmail-badger_Postgres >/dev/null 2>&1 && ( \
			docker volume rm appmail-badger_Postgres >/dev/null 2>&1 && \
				 echo "successfully deleted volume Postgres" || \
				{ echo "[ERROR] failed to delete volume" >&2; exit 1; } )	

	@docker image prune -f -a --filter "label=com.docker.compose.project=${PROJECT_NAME}" \
	  || { echo "[ERROR] failed to remove docker ${PROJECT_NAME} images - $?" >&2; exit 1; }
