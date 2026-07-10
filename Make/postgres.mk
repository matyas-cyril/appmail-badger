.docker_postgres_up: .init_dot.env .init_dockerfile
	@docker compose -f ${DOCKER_YAML} --env-file dot.env up postgres -d --remove-orphans \
	  || { echo "[ERROR] failed to construct docker Postgres service" >&2; exit 1; }

.docker_postgres_stop: .init_dot.env
	@docker compose -f ${DOCKER_YAML} --env-file dot.env stop postgres \
	  || { echo "[ERROR] failed to down docker Postgres service" >&2; exit 1; }

.docker_postgres_rm: .init_dot.env
	@docker compose -f ${DOCKER_YAML} --env-file dot.env rm -sf postgres \
	  || { echo "[ERROR] failed to remove docker Postgres service" >&2; exit 1; }

.docker_postgres_purge: .init_dot.env
	@docker compose -f ${DOCKER_YAML} --env-file dot.env rm -v -sf postgres \
	  || { echo "[ERROR] failed to purge docker Postgres service - $?" >&2; exit 1; }
	@echo "successfully purge docker Postgres service";

	@docker volume rm appmail-badger_Postgres \
	  || { echo "[ERROR] failed to delete volume Postgres - $?" >&2; exit 1; }
	@echo "successfully deleted volume Postgres";
