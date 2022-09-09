.PHONY: help
help:
	@awk -F: '/^[a-zA-Z0-9_-]+:.*$$/{print "make " $$1}' $(MAKEFILE_LIST)

.PHONY: init
init:
	terraform init

.PHONY: update-plugin
update-plugin:
	sed -i.bak -e '/hashes/{N;N;d;}' .terraform.lock.hcl && rm -f .terraform.lock.hcl.bak
	rm -rf .terraform/providers/git.ocom.com
	$(MAKE) init

.PHONY: plan
plan:
	terraform plan

.PHONY: destroy
destroy:
	rm -rf terraform.tfstate terraform.tfstate.backup .terraform.lock.hcl .terraform
