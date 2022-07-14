FROM golang:1.18-alpine AS gobase
RUN apk add --no-cache \
        git && true
WORKDIR /src

FROM gobase as goci
ENV GOLANGCI_LINT_CACHE=/tmp
RUN apk add --no-cache \
        gcc musl-dev make \
    && true
COPY . /src
RUN go install golang.org/x/lint/golint@latest
RUN go install golang.org/x/tools/cmd/goimports@latest
RUN go install github.com/golangci/golangci-lint/cmd/golangci-lint@v1.46.2
WORKDIR /src

FROM gobase as gobuilder
COPY --from=gobase /src /src
