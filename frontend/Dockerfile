FROM node:24-alpine AS build

ARG VITE_API_URL
ENV VITE_API_URL=${VITE_API_URL}

RUN corepack enable

WORKDIR /app

COPY frontend/pnpm-lock.yaml frontend/pnpm-workspace.yaml frontend/package.json ./
COPY frontend/apps/web/package.json ./apps/web/
COPY frontend/packages/config/package.json ./packages/config/
COPY frontend/packages/contracts/package.json ./packages/contracts/
COPY frontend/packages/env/package.json ./packages/env/
COPY frontend/packages/ui/package.json ./packages/ui/

RUN pnpm install --frozen-lockfile

COPY frontend/ ./

RUN test -n "$VITE_API_URL" || (echo "VITE_API_URL build arg is required" && exit 1)
RUN pnpm -F web build

FROM caddy:2-alpine

COPY frontend/Caddyfile /etc/caddy/Caddyfile
COPY --from=build /app/apps/web/dist /srv

EXPOSE 8080
