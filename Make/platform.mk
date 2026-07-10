
.up_platform: .init_dot.env .init_dockerfile
	@docker compose -f ${DOCKER_YAML} --env-file dot.env up -d --remove-orphans \
	  || { echo "[ERROR] failed to construct docker appmail-badger plateform" >&2; exit 1; }

.stop_platform: .init_dot.env
	@docker compose -f ${DOCKER_YAML} --env-file dot.env stop \
	  || { echo "[ERROR] failed to remove docker appmail-badger plateform - $?" >&2; exit 1; };
	@echo "successfully delete plateform appmail-badger";

.purge_platform: .init_dot.env
	@docker compose -f ${DOCKER_YAML} --env-file dot.env down -v --remove-orphans \
	  || { echo "[ERROR] failed to purge docker appmail-badger plateform - $?" >&2; exit 1; }
	@docker image prune -f -a --filter "label=com.docker.compose.project=${PROJECT_NAME}" \
	  || { echo "[ERROR] failed to remove docker ${PROJECT_NAME} images - $?" >&2; exit 1; }
	@echo "successfully purge plateform appmail-badger";

.rm_platform: .init_dot.env
	@docker compose -f ${DOCKER_YAML} --env-file dot.env down \
	  || { echo "[ERROR] failed to remove docker appmail-badger plateform - $?" >&2; exit 1; };
	@echo "successfully delete plateform appmail-badger";
