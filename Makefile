export GO111MODULE=on

TIMESTAMP := $(shell date '+%m%d%H%M%Y.%S')
RELEASE_TAG   ?= $(TIMESTAMP)

# Default Go linker flags.
GO_LDFLAGS ?= -ldflags="-s -w -X main.Version=v${RELEASE_TAG}"

# Binary name.
DDF_OSX := ./bin/ddf-osx

.PHONY: all
all: clean lint $(DDF_OSX) test

$(DDF_OSX):
	GOOS=darwin GOARCH=amd64 go build $(GO_LDFLAGS) $(BUILDARGS) -o $@ .

.PHONY: vendor
vendor:
	go mod tidy
	go mod vendor

.PHONY: lint
lint:
	@ golangci-lint run --fast

.PHONY: test
test:
	go test -mod=vendor -timeout=30s $(TESTARGS) ./...

.PHONY: cover
cover:
	@$(MAKE) test TESTARGS="-coverprofile=coverage.out"
	@go tool cover -html=coverage.out
	@rm -f coverage.out

.PHONY: clean
clean:
	@rm -rf ./bin

.PHONY: package
package: all
	zip -j bin/ddf-osx.zip $(DDF_OSX)
	shasum -a 256 bin/ddf-osx.zip > bin/ddf-osx.sha256