
-include Make/Docker/platform.mk
-include Make/Docker/postgres.mk
-include Make/Docker/smtp.mk
-include Make/Docker/cyrus.mk
-include Make/Docker/redis.mk
-include Make/Docker/webmail.mk
-include Make/Docker/adminer.mk

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
		for dir in Docker/DEBIAN; do \
			sed "\|\#\ \%COPY_GPG_CMD\%|d" ./$$dir/Dockerfile.template | sed "\|\#\ \%COPY_SRC_CMD\%|d" - > "./$$dir/Dockerfile"; \
			rm -rf "./Docker/$$dir/sources"; \
		done \
	fi

platform-up: .up_platform
platform-stop: .stop_platform
platform-rm: .rm_platform
platform-purge: .purge_platform

postgres-up: .docker_postgres_up
postgres-stop: .docker_postgres_stop
postgres-rm: .docker_postgres_rm
postgres-purge: .docker_postgres_purge

smtp-up: .docker_smtp_up
smtp-stop: .docker_smtp_stop
smtp-rm: .docker_smtp_rm
smtp-purge: .docker_smtp_purge

cyrus-up: .docker_cyrus_up
cyrus-stop: .docker_cyrus_stop
cyrus-rm: .docker_cyrus_rm
cyrus-purge: .docker_cyrus_purge

redis-up: .docker_redis_up
redis-stop: .docker_redis_stop
redis-rm: .docker_redis_rm
redis-purge: .docker_redis_purge

webmail-up: .docker_webmail_up
webmail-stop: .docker_webmail_stop
webmail-rm: .docker_webmail_rm
webmail-purge: .docker_webmail_purge

adminer-up: .docker_adminer_up
adminer-stop: .docker_adminer_stop
adminer-rm: .docker_adminer_rm
adminer-purge: .docker_adminer_purge
