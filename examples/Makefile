.PHONY: help
help:
	@awk -F: '/^[a-zA-Z0-9_-]+:.*$$/{print "make " $$1}' $(MAKEFILE_LIST)

init:
	terraform init

update-plugin:
	sed -i.bak -e '/hashes/{N;N;d;}' .terraform.lock.hcl && rm -f .terraform.lock.hcl.bak
	rm -rf .terraform/providers/git.ocom.com
	$(MAKE) init

plan:
	terraform plan

destroy:
	rm -rf terraform.tfstate terraform.tfstate.backup .terraform.lock.hcl .terraform
