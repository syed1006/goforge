BINARY      := demo
PKG         := example.com/demo
CMD         := ./cmd/$(BINARY)
BIN_DIR     := bin
VERSION     ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS     := -s -w -X $(PKG)/internal/buildinfo.Version=$(VERSION)

.PHONY: all build install run test lint fmt tidy clean docker compose

all: build

build:
	@mkdir -p $(BIN_DIR)
	go build -trimpath -ldflags '$(LDFLAGS)' -o $(BIN_DIR)/$(BINARY) $(CMD)

install:
	go install -trimpath -ldflags '$(LDFLAGS)' $(CMD)

run:
	go run $(CMD)

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

docker:
	docker build -t $(BINARY):$(VERSION) .

compose:
	docker compose up --build

