FROM --platform=$BUILDPLATFORM node:22.22.0-alpine3.23@sha256:e4bf2a82ad0a4037d28035ae71529873c069b13eb0455466ae0bc13363826e34 AS web
WORKDIR /src
RUN npm install --global --ignore-scripts pnpm@10.13.1
COPY package.json pnpm-workspace.yaml pnpm-lock.yaml .npmrc ./
COPY apps/web/package.json apps/web/package.json
RUN --mount=type=cache,target=/root/.local/share/pnpm/store pnpm install --frozen-lockfile --ignore-scripts
COPY apps/web/ apps/web/
RUN pnpm --dir apps/web build
FROM --platform=$BUILDPLATFORM golang:1.26.11-alpine3.23@sha256:23f198931d4bd83b4f5a613bd75656257ae8bf832dcbc2db80e1050b5219ca2a AS build
WORKDIR /src
COPY apps/api/go.mod apps/api/go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod go mod download
COPY apps/api/ ./
ARG TARGETOS
ARG TARGETARCH
ARG VERSION=dev
RUN --mount=type=cache,target=/go/pkg/mod --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 GOOS=$TARGETOS GOARCH=$TARGETARCH go build -trimpath -buildvcs=false -ldflags="-s -w" -o /out/spa-server ./cmd/spa-server
FROM gcr.io/distroless/static-debian13:nonroot@sha256:8ce6dd87e9823c5b48dcfb2c8604714ccb351ce59607327bcc4957448da279b5
COPY --from=build --chown=nonroot:nonroot /out/spa-server /app/spa-server
COPY --from=web --chown=nonroot:nonroot /src/apps/web/dist /srv
USER nonroot:nonroot
EXPOSE 8080
HEALTHCHECK --interval=30s --timeout=6s CMD ["/app/spa-server", "healthcheck"]
ENTRYPOINT ["/app/spa-server"]
