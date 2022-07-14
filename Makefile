NAME ?= terraform-provider-leaseweb
VERSION ?= 0.0.1
GOOS ?= linux
GOARCH ?= amd64
BINARY = $(NAME)-$(VERSION)-$(GOOS)-$(GOARCH)

default: install

build:
	mkdir -p dist
	CGO_ENABLED=0 GOOS=$(GOOS) GOARCH=$(GOARCH) go build -a -tags netgo -ldflags '-w' -o dist/$(BINARY)

ci:
	golangci-lint run --disable-all -E gofmt -E whitespace -E errcheck

release:
	$(MAKE) build GOOS=darwin GOARCH=amd64
	$(MAKE) build GOOS=darwin GOARCH=arm64
	$(MAKE) build GOOS=freebsd GOARCH=amd64
	$(MAKE) build GOOS=linux GOARCH=amd64
	$(MAKE) build GOOS=windows GOARCH=amd64

install: build
	mkdir -p ~/.terraform.d/plugins/git.ocom.com/infra/leaseweb/$(VERSION)/$(GOOS)_$(GOARCH)
	mv $(BINARY) ~/.terraform.d/plugins/git.ocom.com/infra/leaseweb/$(VERSION)/$(GOOS)_$(GOARCH)/$(NAME)
