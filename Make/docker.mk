
-include Make/Docker/platform.mk
-include Make/Docker/postgres.mk
-include Make/Docker/smtp.mk
-include Make/Docker/cyrus.mk
-include Make/Docker/redis.mk
-include Make/Docker/webmail.mk
-include Make/Docker/adminer.mk
-include Make/Docker/golang.mk

# Créer un fichier vide pour éviter les if dans les déclarations suivantes
.init_dot.env:
	@if [ ! -f "dot.env" ]; then touch dot.env; fi

# Initialiser le Dockerfile à partir du template
.init_dockerfile:
	@if [ -d "./sources" ] && [ -n "$(COPY_SRC_CMD)" ] && [ -n "$(COPY_GPG_CMD)" ]; then \
		for dir in Docker/DEBIAN; do \
			cp -rf ./sources/ "./$$dir/sources"; \
			sed -e "s|\#\ \%COPY_GPG_CMD\%|$(COPY_GPG_CMD)|g" ./$$dir/Dockerfile.template | sed -e "s|\#\ \%COPY_SRC_CMD\%|$(COPY_SRC_CMD)|g" - > "./Docker/$$dir/Dockerfile"; \
		done \
	else \
		for dir in Docker/DEBIAN Docker/GOLANG; do \
			sed "\|\#\ \%COPY_GPG_CMD\%|d" ./$$dir/Dockerfile.template | sed "\|\#\ \%COPY_SRC_CMD\%|d" - > "./$$dir/Dockerfile"; \
			rm -rf "./Docker/$$dir/sources"; \
		done \
	fi

docker-platform-up: .up_platform
docker-platform-stop: .stop_platform
docker-platform-rm: .rm_platform
docker-platform-purge: .purge_platform

docker-postgres-up: .docker_postgres_up
docker-postgres-stop: .docker_postgres_stop
docker-postgres-rm: .docker_postgres_rm
docker-postgres-purge: .docker_postgres_purge

docker-smtp-up: .docker_smtp_up
docker-smtp-stop: .docker_smtp_stop
docker-smtp-rm: .docker_smtp_rm
docker-smtp-purge: .docker_smtp_purge

docker-cyrus-up: .docker_cyrus_up
docker-cyrus-stop: .docker_cyrus_stop
docker-cyrus-rm: .docker_cyrus_rm
docker-cyrus-purge: .docker_cyrus_purge

docker-redis-up: .docker_redis_up
docker-redis-stop: .docker_redis_stop
docker-redis-rm: .docker_redis_rm
docker-redis-purge: .docker_redis_purge

docker-webmail-up: .docker_webmail_up
docker-webmail-stop: .docker_webmail_stop
docker-webmail-rm: .docker_webmail_rm
docker-webmail-purge: .docker_webmail_purge

docker-adminer-up: .docker_adminer_up
docker-adminer-stop: .docker_adminer_stop
docker-adminer-rm: .docker_adminer_rm
docker-adminer-purge: .docker_adminer_purge

docker-golang-up: .docker_golang_up
docker-golang-stop: .docker_golang_stop
docker-golang-rm: .docker_golang_rm
docker-golang-purge: .docker_golang_purge
