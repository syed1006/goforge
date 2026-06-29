# goforge

`goforge` is an opinionated Go project scaffolder. It generates a production-leaning
project layout with the integrations you actually want: an HTTP framework, a
database driver, gRPC, GraphQL, hot-reload, linting, Docker, and CI — all wired up
and ready to `go run`.

It is inspired by [`go-blueprint`](https://github.com/melkeydev/go-blueprint) but
extends it across the stack and is built as a pluggable generator pipeline so new
features can be added without touching the core.

## Highlights

- Interactive prompts (powered by [`huh`](https://github.com/charmbracelet/huh)) or
  fully flag-driven for CI use.
- Choice of HTTP framework: `stdlib`, `chi`, `gin`, `fiber`, `echo`.
- Choice of database driver: `postgres`, `mysql`, `sqlite`, `mongo`, `redis`, or
  none.
- Optional gRPC server (with `buf` config and a sample service).
- Optional GraphQL server via `gqlgen`.
- Optional hot-reload via [`air`](https://github.com/air-verse/air).
- Optional linting via `golangci-lint` with a sensible default config and a
  `pre-commit` hook.
- Optional Docker + `docker-compose` for the selected database.
- Optional GitHub Actions CI (lint + test + build).

## Install

```sh
go install github.com/syed1006/goforge/cmd/goforge@latest
```

Or, from a checkout:

```sh
make install
```

## Usage

Interactive:

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
  --hot-reload \
  --lint \
  --docker \
  --ci \
  --no-interactive
```

## Architecture

`goforge` is a pipeline of independent generators. Each generator declares which
files it produces; the orchestrator resolves dependencies, renders embedded
templates, and writes the result to disk. Adding a new feature is a matter of
implementing the `Generator` interface and registering it.
