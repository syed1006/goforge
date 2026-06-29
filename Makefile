BINARY      := goforge
PKG         := github.com/syed1006/goforge
CMD         := ./cmd/$(BINARY)
BIN_DIR     := bin
VERSION     ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS     := -s -w -X $(PKG)/internal/cli.version=$(VERSION)

.PHONY: all build install run test lint fmt tidy clean

all: build

build:
	@mkdir -p $(BIN_DIR)
	go build -trimpath -ldflags '$(LDFLAGS)' -o $(BIN_DIR)/$(BINARY) $(CMD)

install:
	go install -trimpath -ldflags '$(LDFLAGS)' $(CMD)

run:
	go run $(CMD) $(ARGS)

test:
	go test ./... -race -count=1

lint:
	@command -v golangci-lint >/dev/null 2>&1 || { echo "golangci-lint not installed"; exit 1; }
	golangci-lint run ./...

fmt:
	gofmt -s -w .
	@command -v goimports >/dev/null 2>&1 && goimports -w . || true

tidy:
	go mod tidy

clean:
	rm -rf $(BIN_DIR)
