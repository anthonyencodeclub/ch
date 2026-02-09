SHELL := /bin/bash

.DEFAULT_GOAL := build

.PHONY: build ch help fmt lint test ci tools

BIN_DIR := $(CURDIR)/bin
BIN := $(BIN_DIR)/ch
CMD := ./cmd/ch

VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
COMMIT := $(shell git rev-parse --short=12 HEAD 2>/dev/null || echo "")
DATE := $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -X github.com/openclaw/ch/internal/cmd.version=$(VERSION) -X github.com/openclaw/ch/internal/cmd.commit=$(COMMIT) -X github.com/openclaw/ch/internal/cmd.date=$(DATE)

TOOLS_DIR := $(CURDIR)/.tools
GOLANGCI_LINT := $(TOOLS_DIR)/golangci-lint

ifneq ($(filter ch,$(MAKECMDGOALS)),)
RUN_ARGS := $(wordlist 2,$(words $(MAKECMDGOALS)),$(MAKECMDGOALS))
$(eval $(RUN_ARGS):;@:)
endif

build:
	@mkdir -p $(BIN_DIR)
	@go build -ldflags "$(LDFLAGS)" -o $(BIN) $(CMD)

ch: build
	@if [ -n "$(RUN_ARGS)" ]; then \
		$(BIN) $(RUN_ARGS); \
	elif [ -z "$(ARGS)" ]; then \
		$(BIN) --help; \
	else \
		$(BIN) $(ARGS); \
	fi

help: build
	@$(BIN) --help

tools:
	@mkdir -p $(TOOLS_DIR)
	@GOBIN=$(TOOLS_DIR) go install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@latest

fmt:
	@gofmt -s -w .

lint: tools
	@$(GOLANGCI_LINT) run

test:
	@go test ./...

ci: fmt lint test
