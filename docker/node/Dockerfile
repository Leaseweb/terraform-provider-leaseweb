FROM node:iron-slim@sha256:5bba9bfa8c663a8bb5f014d39ea97550946b355b49308c2d1305d706489e8e14

RUN corepack enable \
    && corepack prepare pnpm@latest-9 --activate \
    && pnpm config set store-dir .pnpm-store
