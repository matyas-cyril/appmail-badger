
# Attendre que la socket PG soit dispo
.waiting_postgres:
	@cpt=0; \
	while ! docker exec postgres_$(PROJECT_NAME) test -S '/var/run/postgresql/.s.PGSQL.5432'; do \
        if [ $$cpt -ge 10 ]; then \
            echo "[ERROR] PostgreSQL did not start in time"; \
            exit 1; \
        fi; \
		 echo "[WARN] waiting postgres..."; \
		sleep 1; \
		cpt=$$((cpt+1)); \
	done

sql-init: .docker_postgres_up .waiting_postgres
	@docker exec -i postgres_$(PROJECT_NAME) psql -U $(PG_USER) -d $(BDD) < $(SQL_INIT) \
	  || { echo "[ERROR] failed to init database" >&2; exit 1; }

sql-drop: .docker_postgres_up .waiting_postgres
	@docker exec -i postgres_$(PROJECT_NAME) psql -U $(PG_USER) -d $(BDD) < $(SQL_DROP) \
	  || { echo "[ERROR] failed to drop database" >&2; exit 1; }

sql-data: .docker_postgres_up .waiting_postgres
	@docker exec -i postgres_$(PROJECT_NAME) psql -U $(PG_USER) -d $(BDD) < $(SQL_DATA) \
	  || { echo "[ERROR] failed to drop database" >&2; exit 1; }
