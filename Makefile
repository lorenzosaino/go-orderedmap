SHELL = /bin/bash -euo pipefail

# Go binary for non-container targets
GO ?= go

# Go version for container targets
GO_VERSION ?= latest

# Container image for container targets
# To use Go tip, set CONTAINER=aleksi/golang-tip:master
CONTAINER ?= golang:$(GO_VERSION)

PKG = github.com/lorenzosaino/go-orderedmap
DOCKER_RUN_FLAGS = --rm -it -v $$(pwd):/go/src/$(PKG) -w /go/src/$(PKG)

export GO111MODULE=on

all: fmt-check lint vet staticcheck test ## Run all checks and tests

.PHONY: mod-upgrade
mod-upgrade: ## Upgrade all vendored dependencies
	$(GO) get -d -u -t ./...
	$(GO) mod tidy
	$(GO) mod vendor

.PHONY: mod-update
mod-update: ## Ensure all used dependencies are tracked in go.{mod|sum} and vendored
	$(GO) mod tidy
	$(GO) mod vendor

.PHONY: fmt-check
fmt-check: ## Validate that all source files pass "go fmt"
	exit $(shell $(GO) fmt ./... | wc -l)

.PHONY: lint
lint: ## Run go lint
	@[ -x "$(shell which golint)" ] || $(GO) install ./vendor/golang.org/x/lint/golint 2>/dev/null || $(GO) get -u golang.org/x/lint/golint
	@# We need to explicitly exclude ./vendor because of https://github.com/golang/lint/issues/320
	golint -set_exit_status $(shell $(GO)  list ./... | grep -v '/vendor/' | grep -v 'internal/list')

.PHONY: vet
vet: ## Run go vet
	$(GO) vet ./...

.PHONY: staticcheck
staticcheck: ## Run staticcheck
	@[ -x "$(shell which staticcheck)" ] || $(GO) install ./vendor/honnef.co/go/tools/cmd/staticcheck 2>/dev/null || $(GO) get -u honnef.co/go/tools/cmd/staticcheck
	staticcheck ./...

.PHONY: test
test: ## Run all tests
	$(GO) test -race ./...

.PHONY: container-shell
container-shell: ## Open a shell on a Docker container
	docker run $(DOCKER_RUN_FLAGS) $(CONTAINER) /bin/bash

.PHONY: container-%
container-%: ## Run any target of this Makefile in a Docker container
	docker run $(DOCKER_RUN_FLAGS) $(CONTAINER) make $*

.PHONY: help
help: ## Print help
	@(grep -E '^[a-zA-Z0-9_%-]+:.*?## .*$$' Makefile || true )| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
