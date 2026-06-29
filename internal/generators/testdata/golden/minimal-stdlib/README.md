# demo

Generated with [goforge](https://github.com/syed1006/goforge).

## Stack

- **Go**: 1.23
- **HTTP framework**: stdlib
- **Database**: none

## Getting started

```sh
go mod tidy
make run   # one-shot run
```

## Layout

```
cmd/demo/    application entrypoint
internal/                  private application code
```

## Make targets

```
make run       # run locally
make test      # run unit tests
make build     # build a static binary
make lint      # golangci-lint run
make fmt       # gofmt + goimports
make tidy      # go mod tidy
```
