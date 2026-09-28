# Toute la plateforme
.smtp_full_install:
	ANSIBLE_CONFIG=Ansible/ansible.cfg \
		$(VENV_DIR)/bin/ansible-playbook Ansible/smtp.yml --tags install-debian,install-sasl,install-memcached,install-rspamd,install-postfix,install-clamav,install-amavis

.smtp_full_update:
	ANSIBLE_CONFIG=Ansible/ansible.cfg \
		$(VENV_DIR)/bin/ansible-playbook Ansible/smtp.yml --tags update-debian,update-sasl,update-memcached,update-rspamd,update-postfix,update-clamav,update-amavis

.smtp_full_upgrade: .smtp_debian_upgrade

# Debian
.smtp_debian_install:
	ANSIBLE_CONFIG=Ansible/ansible.cfg \
		$(VENV_DIR)/bin/ansible-playbook Ansible/smtp.yml --tags install-debian

.smtp_debian_update:
	ANSIBLE_CONFIG=Ansible/ansible.cfg \
		$(VENV_DIR)/bin/ansible-playbook Ansible/smtp.yml --tags update-debian

.smtp_debian_upgrade:
	ANSIBLE_CONFIG=Ansible/ansible.cfg \
		$(VENV_DIR)/bin/ansible-playbook Ansible/smtp.yml --tags upgrade-debian


# Saslauthd
.smtp_saslauthd_install:
	ANSIBLE_CONFIG=Ansible/ansible.cfg \
		$(VENV_DIR)/bin/ansible-playbook Ansible/smtp.yml --tags install-sasl

.smtp_saslauthd_update:
	ANSIBLE_CONFIG=Ansible/ansible.cfg \
		$(VENV_DIR)/bin/ansible-playbook Ansible/smtp.yml --tags update-sasl

# Postix
.smtp_postfix_install:
	ANSIBLE_CONFIG=Ansible/ansible.cfg \
		$(VENV_DIR)/bin/ansible-playbook Ansible/smtp.yml --tags install-postfix

.smtp_postfix_update:
	ANSIBLE_CONFIG=Ansible/ansible.cfg \
		$(VENV_DIR)/bin/ansible-playbook Ansible/smtp.yml --tags update-postfix

# Antivirus
.smtp_antivirus_install:
	ANSIBLE_CONFIG=Ansible/ansible.cfg \
		$(VENV_DIR)/bin/ansible-playbook Ansible/smtp.yml --tags install-clamav,install-amavis

.smtp_antivirus_update:
	ANSIBLE_CONFIG=Ansible/ansible.cfg \
		$(VENV_DIR)/bin/ansible-playbook Ansible/smtp.yml --tags update-clamav,update-amavis

# Clamav
.smtp_clamav_install:
	ANSIBLE_CONFIG=Ansible/ansible.cfg \
		$(VENV_DIR)/bin/ansible-playbook Ansible/smtp.yml --tags install-clamav

.smtp_clamav_update:
	ANSIBLE_CONFIG=Ansible/ansible.cfg \
		$(VENV_DIR)/bin/ansible-playbook Ansible/smtp.yml --tags update-clamav

# Amavis
.smtp_amavis_install:
	ANSIBLE_CONFIG=Ansible/ansible.cfg \
		$(VENV_DIR)/bin/ansible-playbook Ansible/smtp.yml --tags install-amavis

.smtp_amavis_update:
	ANSIBLE_CONFIG=Ansible/ansible.cfg \
		$(VENV_DIR)/bin/ansible-playbook Ansible/smtp.yml --tags update-amavis

# Memcached
.smtp_memcached_install:
	ANSIBLE_CONFIG=Ansible/ansible.cfg \
		$(VENV_DIR)/bin/ansible-playbook Ansible/smtp.yml --tags install-memcached

.smtp_memcached_update:
	ANSIBLE_CONFIG=Ansible/ansible.cfg \
		$(VENV_DIR)/bin/ansible-playbook Ansible/smtp.yml --tags update-memcached

# RSpamD
.smtp_rspamd_install:
	ANSIBLE_CONFIG=Ansible/ansible.cfg \
		$(VENV_DIR)/bin/ansible-playbook Ansible/smtp.yml --tags install-rspamd

.smtp_rspamd_update:
	ANSIBLE_CONFIG=Ansible/ansible.cfg \
		$(VENV_DIR)/bin/ansible-playbook Ansible/smtp.yml --tags update-rspamd
