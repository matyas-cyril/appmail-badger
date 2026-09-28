
-include Make/Ansible/smtp.mk

.make_venv:
	@if [ ! -d $(VENV_DIR) ] || [ ! -f $(VENV_DIR)/bin/activate ]; then \
		$(PYTHON) -m venv $(VENV_DIR) \
			|| { echo "virtual environment failed in '$(VENV_DIR)'"; exit 2; }; \
		echo "virtual environment created in '$(VENV_DIR)'"; \
	else \
		echo "virtual environment '$(VENV_DIR)' already exist"; \
	fi

.install_ansible: .make_venv
	@echo "starting install of Ansible $(ANSIBLE_VERSION)";
	@$(VENV_DIR)/bin/pip3 install ansible==$(ANSIBLE_VERSION) \
		|| { echo "ansible installation failed"; exit 2; }; \
	echo "successfully installed Ansible $(ANSIBLE_VERSION)";

.rm_venv:
	@rm -rf $(VENV_DIR) || { echo "failed to remove '$(VENV_DIR)'"; exit 2; }; \
	echo "successfully removed '$(VENV_DIR)'";

ansible-init: .install_ansible
ansible-delete: .rm_venv

# Déployer toute la plateforme
smtp-full-install: .smtp_full_install
smtp-full-update: .smtp_full_update
smtp-full-upgrade: .smtp_full_upgrade
smtp-os-ready: docker_smtp_rm .smtp_debian_install

smtp-saslauthd-install: .smtp_saslauthd_install
smtp-saslauthd-update: .smtp_saslauthd_update

smtp-postfix-install: .smtp_postfix_install
smtp-postfix-update: .smtp_postfix_update
smtp-postfix-data: .smtp_postfix_update_data

smtp-antivirus-install: .smtp_antivirus_install
smtp-antivirus-update: .smtp_antivirus_update

smtp-clamav-install: .smtp_clamav_install
smtp-clamav-update: .smtp_clamav_update

smtp-amavis-install: .smtp_amavis_install
smtp-amavis-update: .smtp_amavis_update

smtp-debian-install: .smtp_debian_install
smtp-debian-update: .smtp_debian_update
smtp-debian-upgrade: .smtp_debian_upgrade

smtp-memcached-install: .smtp_memcached_install
smtp-memcached-update: .smtp_memcached_update

smtp-rspamd-install: .smtp_rspamd_install
smtp-rspamd-update: .smtp_rspamd_update
