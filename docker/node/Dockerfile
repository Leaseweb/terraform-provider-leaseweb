FROM node:jod-slim@sha256:ab1ba25996833ac547afb1dd8d59f4d31bd2b29b5574bf20b07e57cb5726d274

RUN corepack enable \
    && corepack prepare pnpm@latest-9 --activate \
    && pnpm config set store-dir .pnpm-store
