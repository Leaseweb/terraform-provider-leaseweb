FROM golang:1.18-alpine AS godev
RUN apk add --no-cache \
        git \
        make \
    && true
RUN go install golang.org/x/lint/golint@latest
RUN go install golang.org/x/tools/cmd/goimports@latest
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.46.2
WORKDIR /src
