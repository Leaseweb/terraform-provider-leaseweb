FROM golang:1.22.1-alpine3.19 AS godev
RUN apk add --no-cache \
        git=2.43.0-r0  \
        gpg=2.4.4-r0  \
        grep=3.11-r0 \
        make=4.4.1-r2 \
    && true
RUN go install golang.org/x/lint/golint@v0.0.0-20210508222113-6edffad5e616 &&\
    go install github.com/kisielk/errcheck@v1.7.0 &&\
    go install golang.org/x/tools/cmd/goimports@v0.19.0 &&\
    go install github.com/goreleaser/goreleaser@v1.24.0
ENV CGO_ENABLED=0
WORKDIR /src
