# goforge

`goforge` is an opinionated Go project scaffolder. It generates a production-leaning
project layout with the integrations you actually want — HTTP framework, database
driver, gRPC, GraphQL, hot reload, linting, Docker, CI — all wired up and ready
to `go run`.

It is inspired by [`go-blueprint`](https://github.com/melkeydev/go-blueprint) but
extends the stack and is built as a pluggable generator pipeline so new features
can be added without touching the core.

## Features

- **HTTP framework**: `stdlib`, `chi`, `gin`, `fiber`, or `echo` — every variant
  ships with structured `slog` logging, `/healthz` and `/readyz` probes, and
  graceful shutdown.
- **Database driver**: `postgres` (pgx pool), `mysql` (database/sql), `sqlite`
  (pure-Go modernc), `mongo`, `redis`, or `none`. Wired into `main.go` with a
  bounded context for both connect and close.
- **gRPC**: optional `internal/grpcsrv` server with health-check, reflection,
  and a slog interceptor; `buf.yaml` + `buf.gen.yaml` + a sample
  `proto/<service>/v1/ping.proto` so `buf generate` works out of the box.
- **GraphQL**: optional gqlgen integration. The scaffolder runs
  `gqlgen generate` as a post-step, mounts `/graphql` + `/playground` in the
  selected framework, and supplies a framework-agnostic handler in
  `internal/graphql`.
- **OpenTelemetry**: optional tracing via `internal/telemetry` with an
  OTLP/HTTP exporter, W3C TraceContext + Baggage propagators, and the
  framework-native instrumentation middleware wired in. When gRPC is also
  enabled, the gRPC server gets `otelgrpc.StatsHandler` automatically.
- **Prometheus metrics**: optional `/metrics` endpoint with the default
  registry (Go runtime + process collectors).
- **Hot reload**: `.air.toml` tuned for the project binary.
- **Lint**: `.golangci.yml` with a curated linter set + goimports configured
  for the project module, `.pre-commit-config.yaml`, and `.editorconfig`.
- **Docker**: cache-aware multi-stage `Dockerfile` landing on
  `gcr.io/distroless/static:nonroot`, plus a `docker-compose.yml` that wires
  the app to the selected database service.
- **CI**: `.github/workflows/ci.yml` with test (race), lint (when enabled),
  and build jobs — automatically running `gqlgen generate` when needed.

## Install

```sh
go install github.com/syed1006/goforge/cmd/goforge@latest
```

Or, from a checkout:

```sh
make install
```

## Usage

Interactive (everything you don't pass via flags is prompted):

```sh
goforge new
```

Flag-driven (no prompts):

```sh
goforge new myapi \
  --module github.com/me/myapi \
  --framework gin \
  --database postgres \
  --grpc \
  --graphql \
  --otel \
  --metrics \
  --hot-reload \
  --lint \
  --docker \
  --ci \
  --no-interactive
```

Other useful commands:

```sh
goforge list frameworks    # show supported HTTP frameworks
goforge list databases     # show supported database drivers
goforge version
```

## Architecture

```
cmd/goforge/                   entrypoint
internal/cli/                  Cobra commands (root, new, list, version)
internal/prompt/               huh-driven interactive form
internal/config/               resolved Config + validation + enums
internal/generator/            Generator interface + Manifest + Registry
internal/scaffold/             orchestrator (Run)
internal/template/             text/template wrapper over an fs.FS
internal/fsutil/               file writer with error/overwrite/skip modes
internal/templates/files/      embedded scaffold templates
internal/generators/
  base/                        README, Makefile, .gitignore, .env, config
  framework/                   stdlib | chi | gin | fiber | echo
  database/                    postgres | mysql | sqlite | mongo | redis
  grpc/                        gRPC server + buf + sample proto
  graphql/                     gqlgen + handler + schema
  otel/                        internal/telemetry + framework middleware deps
  metrics/                     internal/metrics + /metrics route
  hotreload/                   .air.toml
  lint/                        .golangci.yml + pre-commit + editorconfig
  docker/                      Dockerfile + docker-compose + .dockerignore
  ci/                          GitHub Actions workflow
```

Each feature is a `Generator` that opts in via `Applies(cfg)` and writes files
through a shared `Context` (renderer + writer + module manifest). The
orchestrator runs every applicable generator, then resolves manifest
dependencies via `go get`, optionally runs `gqlgen generate`, and finishes
with `go mod tidy` + `gofmt`.

Adding a new feature is a matter of:

1. Drop your templates under `internal/templates/files/<feature>/`.
2. Implement the `Generator` interface in `internal/generators/<feature>/`.
3. Register it in `internal/generators/registry.go`.

## Develop

```sh
make build           # build the CLI
make test            # run the test suite
make run ARGS="new"  # run goforge from source
```
