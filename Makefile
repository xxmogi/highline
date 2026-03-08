BINARY_NAME = highline
MAIN        = ./cmd/highline
GOFLAGS     ?=
DIST_DIR    = dist

PLATFORMS = \
	linux/amd64 \
	linux/arm64 \
	darwin/amd64 \
	darwin/arm64 \
	windows/amd64

.PHONY: build test lint clean install dist help $(PLATFORMS)

## build: Build the binary for the current platform
build:
	go build $(GOFLAGS) -o $(BINARY_NAME) $(MAIN)

## test: Run all tests
test:
	go test ./...

## lint: Run go vet
lint:
	go vet ./...

## clean: Remove build artifacts
clean:
	rm -f $(BINARY_NAME)
	rm -rf $(DIST_DIR)

## install: Install the binary to GOPATH/bin
install:
	go install $(MAIN)

## dist: Cross-compile for all platforms into dist/
dist: $(PLATFORMS)

$(PLATFORMS):
	$(eval OS   := $(word 1, $(subst /, ,$@)))
	$(eval ARCH := $(word 2, $(subst /, ,$@)))
	$(eval EXT  := $(if $(filter windows,$(OS)),.exe,))
	$(eval OUT  := $(DIST_DIR)/$(BINARY_NAME)-$(OS)-$(ARCH)$(EXT))
	@mkdir -p $(DIST_DIR)
	GOOS=$(OS) GOARCH=$(ARCH) go build $(GOFLAGS) -o $(OUT) $(MAIN)
	@echo "  built $(OUT)"

## help: Show available targets
help:
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/^## //'
