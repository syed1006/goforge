# demo

Generated with [goforge](https://github.com/syed1006/goforge).

## Stack

- **Go**: 1.23
- **HTTP framework**: chi
- **Database**: redis
- **Lint**: golangci-lint + pre-commit
- **Docker**: enabled
- **CI**: GitHub Actions

## Getting started

```sh
go mod tidy
make run   # one-shot run
```

## Layout

```
cmd/demo/    application entrypoint
internal/                  private application code
internal/storage/          database wiring (redis)
```

## Make targets

```
make run       # run locally
make test      # run unit tests
make build     # build a static binary
make lint      # golangci-lint run
make fmt       # gofmt + goimports
make tidy      # go mod tidy
make docker    # docker build
make compose   # docker compose up
```
