FROM golang:1.20-alpine3.19 AS godev
RUN apk add --no-cache \
        git=2.43.0-r0  \
        gpg=2.4.4-r0  \
        grep=3.11-r0 \
        make=4.4.1-r2 \
    && true
RUN go install golang.org/x/lint/golint@latest
RUN go install github.com/kisielk/errcheck@latest
RUN go install golang.org/x/tools/cmd/goimports@latest
RUN go install github.com/goreleaser/goreleaser@latest
ENV CGO_ENABLED=0
WORKDIR /src
