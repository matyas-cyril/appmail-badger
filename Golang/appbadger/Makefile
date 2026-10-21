
GO       ?= go
APP      ?= appbadger-cli
MAIN     ?= *.go
BIN_DIR  ?= BIN
BINARY   := $(BIN_DIR)/$(APP)
BUILD_DATE ?= $(shell date -u +%Y-%m-%d-%H:%M:%S)

# Optimisations
BUILD_FLAGS ?= -trimpath
LDFLAGS     ?= -s -w

.PHONY: all build run clean

all: build

build:
	@mkdir -p $(BIN_DIR)
	$(GO) build $(BUILD_FLAGS) -ldflags="$(LDFLAGS) -X main.Build_DATE=$(BUILD_DATE)" -o $(BINARY) $(MAIN)

run: build
	./$(BINARY) $(ARGS)

clean:
	@rm -rf $(BIN_DIR)
	$(GO) clean
