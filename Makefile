NAME ?= terraform-provider-leaseweb
VERSION ?= 0.0.1
GOOS ?= linux
GOARCH ?= amd64
BINARY = $(NAME)-$(VERSION)-$(GOOS)-$(GOARCH)

.PHONY: help
help:
	@LC_ALL=C $(MAKE) -pRrq -f $(lastword $(MAKEFILE_LIST)) : 2>/dev/null | awk -v RS= -F: '/^# File/,/^# Finished Make data base/ {if ($$1 !~ "^[#.]") {print $$1}}' | sort | egrep -v -e '^[^[:alnum:]]' -e '^$@$$'

.PHONY: build
build:
	mkdir -p dist
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -a -tags netgo -ldflags '-w' -o dist/$(BINARY)

.PHONY: lint
lint: lint-eol lint-spaces lint-tabs lint-go

.PHONY: lint-eol
lint-eol:
	@echo "==> Validating unix style line endings of files:"
	@! git ls-files | xargs grep --files-with-matches --recursive --exclude Makefile '' || ( echo '[ERROR] Above files have CRLF line endings' && exit 1 )
	@echo All files have valid line endings

.PHONY: lint-spaces
lint-spaces:
	@echo "==> Validating trailing whitespaces in files:"
	@! git ls-files | grep -v '^docs/' | xargs grep --files-with-matches --recursive --extended-regexp ' +$$' || ( echo '[ERROR] Above files have trailing whitespace' && exit 1 )
	@echo No files have trailing whitespaces

.PHONY: lint-tabs
lint-tabs:
	@echo "==> Validating literal tab characters in files:"
	@! git ls-files '*.go' | xargs grep --files-with-matches --recursive --extended-regexp '^ +' || ( echo '[ERROR] Above go files use literal tabs' && exit 1 )
	@echo All files use spaces

.PHONY: lint-go
lint-go:
	golint -set_exit_status ./...
	go vet -v ./...
	errcheck -exclude errcheck_excludes.txt ./...

.PHONY: ci
ci: lint

.PHONY: doc
doc:
	tfplugindocs generate --rendered-provider-name Leaseweb

.PHONY: format
format:
	go fmt ./...
	terraform fmt -recursive examples/

.PHONY: release
release:
	$(MAKE) build GOOS=darwin GOARCH=amd64
	$(MAKE) build GOOS=darwin GOARCH=arm64
	$(MAKE) build GOOS=freebsd GOARCH=amd64
	$(MAKE) build GOOS=linux GOARCH=amd64
	$(MAKE) build GOOS=windows GOARCH=amd64

.PHONY: install
install: build
	mkdir -p ~/.terraform.d/plugins/git.ocom.com/infra/leaseweb/$(VERSION)/$(GOOS)_$(GOARCH)
	mv dist/$(BINARY) ~/.terraform.d/plugins/git.ocom.com/infra/leaseweb/$(VERSION)/$(GOOS)_$(GOARCH)/$(NAME)
