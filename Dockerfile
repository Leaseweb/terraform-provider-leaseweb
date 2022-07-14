FROM golang:1.18-alpine AS godev
RUN apk add --no-cache \
        git \
        grep \
        make \
    && true
RUN go install golang.org/x/lint/golint@latest
RUN go install github.com/kisielk/errcheck@latest
RUN go install golang.org/x/tools/cmd/goimports@latest
WORKDIR /src
