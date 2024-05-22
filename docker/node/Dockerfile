FROM node:iron-slim@sha256:cffed8cd39d6a380434e6d08116d188c53e70611175cd5ec7700f93f32a935a6

RUN corepack enable \
    && corepack prepare pnpm@latest-9 --activate \
    && pnpm config set store-dir .pnpm-store
