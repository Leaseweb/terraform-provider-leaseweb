FROM golang:1.20-alpine AS godev
RUN apk add --no-cache \
        git \
        gpg \
        grep \
        make \
    && true
RUN go install golang.org/x/lint/golint@latest
RUN go install github.com/kisielk/errcheck@latest
RUN go install golang.org/x/tools/cmd/goimports@latest
RUN go install github.com/goreleaser/goreleaser@latest
ENV CGO_ENABLED=0
WORKDIR /src
