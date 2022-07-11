FROM golang:alpine AS godev
RUN apk add --no-cache \
        git \
    && true
RUN go install golang.org/x/lint/golint@latest
RUN go install golang.org/x/tools/cmd/goimports@latest
WORKDIR /src
