FROM node:iron-slim@sha256:ec35a66c9a0a275b027debde05247c081f8b2f0c43d7399d3a6ad5660cee2f6a

RUN corepack enable \
    && corepack prepare pnpm@latest-9 --activate \
    && pnpm config set store-dir .pnpm-store
