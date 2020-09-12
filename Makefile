export GO111MODULE=on

TIMESTAMP := $(shell date '+%m%d%H%M%Y.%S')
RELEASE_TAG   ?= $(TIMESTAMP)

# Default Go linker flags.
GO_LDFLAGS ?= -ldflags="-s -w -X main.Version=v${RELEASE_TAG}"

.PHONY: ci
ci: clean build lint test

.PHONY: build
build:
	@mkdir -p dist
	go build -o ./dist/ddf

# Run all the linters
lint:
	golangci-lint run ./...
	misspell -error **/*
.PHONY: lint

.PHONY: test
test:
	go test -timeout=30s $(TESTARGS) ./...

.PHONY: cover
cover:
	@$(MAKE) test TESTARGS="-coverprofile=coverage.out"
	@go tool cover -html=coverage.out

# Install all the build and lint dependencies
.PHONY: setup
setup:
	go mod download
	go generate -v ./...

.PHONY: clean
clean:
	@rm -rf ./dist
