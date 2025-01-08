FROM node:jod-slim@sha256:0d2188b57802df16244f28bcad8ab8dcad0bc2bc0ed364e7e5550f8b11e48570

RUN corepack enable \
    && corepack prepare pnpm@latest-9 --activate \
    && pnpm config set store-dir .pnpm-store
