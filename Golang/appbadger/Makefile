
GO       ?= go
APP      ?= appbadger-cli
MAIN     ?= *.go
BIN_DIR  ?= BIN
BINARY   := $(BIN_DIR)/$(APP)
BUILD_DATE ?= $(shell date -u +%Y-%m-%d)
BUILD_TIME ?= $(shell date -u +%H:%M:%S)

BUILD_VERSION ?= "0.0.0"

# Optimisations
BUILD_FLAGS ?= -trimpath
LDFLAGS     ?= -s -w

.PHONY: all build run clean

all: build

build:
	@mkdir -p $(BIN_DIR)
	$(GO) build $(BUILD_FLAGS) \
		-ldflags="$(LDFLAGS) \
			-X main.BUILD_VERSION=$(BUILD_VERSION) \
			-X main.BUILD_DATE=$(BUILD_DATE) \
			-X main.BUILD_TIME=$(BUILD_TIME)" \
		-o $(BINARY) $(MAIN)

clean:
	@rm -rf $(BIN_DIR)
	$(GO) clean
