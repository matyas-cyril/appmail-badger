
-include Make/platform.mk
-include Make/postgres.mk
-include Make/smtp.mk
-include Make/cyrus.mk
-include Make/redis.mk
-include Make/webmail.mk

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
