# Terraform Leaseweb Provider

[![Go Reference](https://pkg.go.dev/badge/github.com/leaseweb/terraform-provider-leaseweb.svg)](https://pkg.go.dev/github.com/leaseweb/terraform-provider-leaseweb)

A Terraform provider to manage Leaseweb resources.

## Requirements

- [Terraform](https://developer.hashicorp.com/terraform/downloads) >= 1.7
- [Go](https://golang.org/doc/install) >= 1.21
- [Node](https://nodejs.org) >= 20.13
- [pnPM](https://pnpm.io/) >= 9.1

All requirements are also satisfied by the included [docker-compose.yml](docker-compose.yml).

## Building The Provider

1. Clone the repository
1. Enter the repository directory
1. Build the provider using the Go `install` command:

```shell
go install
```

## Adding Dependencies

This provider uses [Go modules](https://github.com/golang/go/wiki/Modules).
Please see the Go documentation for the most up-to-date information about using Go modules.

To add a new dependency `github.com/author/dependency` to your Terraform provider:

```shell
go get github.com/author/dependency
go mod tidy
```

Then commit the changes to `go.mod` and `go.sum`.

## Developing the Provider

If you wish to work on the provider, you'll first need [Go](http://www.golang.org) installed on your machine (see [Requirements](#requirements) above).

To compile the provider, run `go install`. This will build the provider and put the provider binary in the `$GOPATH/bin` directory.

To generate or update documentation, run `go generate`.

In order to run the full suite of Acceptance tests, run `make testacc`.

_Note:_ Acceptance tests create real resources, and often cost money to run.

```shell
make testacc
```

## Run mock server

To test against API specifications a mock server can be run

```shell
docker-compose up -d
```

Make sure to use `localhost` in your requests as caddy does not know what to do with `127.0.0.1`

```shell

curl -i http://localhost:8080/publicCloud/v1/instances --header 'x-lsw-auth: tralala'
```

## Linting

Files are automatically linted via git hooks on commit and on push. To enable the git hooks run

```bash
pnpm husky
```
