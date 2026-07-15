
.deploy_smtp:
	ANSIBLE_CONFIG=Ansible/ansible.cfg \
		$(VENV_DIR)/bin/ansible-playbook Ansible/smtp.yml --tags install