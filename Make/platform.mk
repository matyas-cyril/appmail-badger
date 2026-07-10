
.up_platform: .init_dot.env .init_dockerfile
	@docker compose -f ${DOCKER_YAML} --env-file dot.env up -d --remove-orphans \
	  || { echo "[ERROR] failed to construct docker appmail-badger plateform" >&2; exit 1; }

.stop_platform: .init_dot.env
	@docker compose -f ${DOCKER_YAML} --env-file dot.env stop \
	  || { echo "[ERROR] failed to remove plateform - $?" >&2; exit 1; };
	@echo "successfully delete plateform";

.purge_platform: .init_dot.env
	@docker compose -f ${DOCKER_YAML} --env-file dot.env down -v --remove-orphans --rmi all \
	  || { echo "[ERROR] failed to purge plateform - $?" >&2; exit 1; }
	@echo "successfully purge plateform";

.rm_platform: .init_dot.env
	@docker compose -f ${DOCKER_YAML} --env-file dot.env down \
	  || { echo "[ERROR] failed to remove plateform - $?" >&2; exit 1; };
	@echo "successfully delete plateform";
