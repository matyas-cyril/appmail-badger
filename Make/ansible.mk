
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
ansible-rm: .rm_venv

smtp-deploy: .deploy_smtp
smtp-deploy-reset: .docker_smtp_rm .deploy_smtp
smtp-deploy-from-scratch: .docker_smtp_purge .deploy_smtp