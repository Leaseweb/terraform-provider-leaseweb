FROM node:jod-slim@sha256:4b44c32c9f3118d60977d0dde5f758f63c4f9eac8ddee4275277239ec600950f

RUN corepack enable \
    && corepack prepare pnpm@latest-9 --activate \
    && pnpm config set store-dir .pnpm-store
