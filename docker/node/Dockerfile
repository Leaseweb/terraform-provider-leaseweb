FROM node:jod-slim@sha256:bca6c602d41be03724f54c8113bb2b76e31e4b135055f19fd0b383a1ca307f39

RUN corepack enable \
    && corepack prepare pnpm@latest-9 --activate \
    && pnpm config set store-dir .pnpm-store
