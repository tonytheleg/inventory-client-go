GOHOSTOS:=$(shell go env GOHOSTOS)
GOPATH:=$(shell go env GOPATH)

ifeq ($(DOCKER),)
DOCKER:=$(shell command -v podman || command -v docker)
endif

TITLE:="Kessel Asset Inventory go client"
ifeq ($(VERSION),)
VERSION:=$(shell git describe --tags --always)
endif
INVENTORY_SCHEMA_VERSION=0.11.0


.PHONY: lint
# run go linter with the repositories lint config
lint:
	@echo "Linting code."
	@$(DOCKER) run -t --rm -v $(PWD):/app -w /app golangci/golangci-lint golangci-lint run -v