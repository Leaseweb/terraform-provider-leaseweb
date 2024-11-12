FROM node:iron-slim@sha256:4eec2c9e4468edc690a31597b5bf7aecaf10e2dffd3fd0521670b6f602c44002

RUN corepack enable \
    && corepack prepare pnpm@latest-9 --activate \
    && pnpm config set store-dir .pnpm-store
