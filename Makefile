SUFFIX	?=
TAG     ?= latest
VERSION ?= v0.0.0

#
# Help
#

help: ## Show this help message.
	@echo
	@echo 'usage: make [target]'
	@echo
	@echo 'targets:'
	@echo
	@egrep '^(.+)\:\ ##\ (.+)' ${MAKEFILE_LIST} | column -t -c 2 -s ':#'
	@echo
.PHONY: help

#
# Init
#

init: ## Initialize the project.
init: docker-compose.local.yml
.PHONY: init

docker-compose.local.yml: ## Generate docker-compose.local.yml.
docker-compose.local.yml:
	@touch docker-compose.local.yml

#
# Build
#

build: ## Build the command for the host platform.
	@mkdir -p bin
	@echo "Building bin/tmpl${SUFFIX} with version ${VERSION}..."
	@CGO_ENABLED=0 go build -ldflags="-X 'main.Version=${VERSION}'" -o bin/tmpl${SUFFIX} .
	@file bin/tmpl${SUFFIX}
.PHONY: build

dist: ## Build commands for all supported platforms.
	@GOOS=darwin  GOARCH=amd64  SUFFIX=-darwin-amd64  make --no-print-directory build
	@GOOS=darwin  GOARCH=arm64  SUFFIX=-darwin-arm64  make --no-print-directory build
	@GOOS=linux   GOARCH=amd64  SUFFIX=-linux-amd64   make --no-print-directory build
	@GOOS=linux   GOARCH=arm64  SUFFIX=-linux-arm64   make --no-print-directory build
.PHONY: dist

#
# Clean
#

clean: ## Clean working directory.
	@rm -rf bin
	@rm -rf coverage
.PHONY: clean

#
# Image
#

define DOCKERFILE
FROM scratch

ARG TARGETOS
ARG TARGETARCH

COPY --chmod=0755 bin/tmpl-$${TARGETOS}-$${TARGETARCH} /tmpl

ENTRYPOINT ["/tmpl"]
endef
export DOCKERFILE

image: ## Build Docker image.
image: | build Dockerfile
	@docker build -t ghcr.io/jeremybower/tmpl:${TAG} .
.PHONY: image

Dockerfile: ## Generate Dockerfile.
	@echo "$$DOCKERFILE" > Dockerfile

#
# Test
#

test: ## Run tests.
test:
	@mkdir -p coverage
	@GOEXPERIMENT=nocoverageredesign go test \
		-timeout 120s \
		-cover \
		-covermode=atomic \
		-coverprofile coverage/coverage.out \
		-count=1 \
		-failfast \
		./...
	@go tool cover \
		-html=coverage/coverage.out \
		-o coverage/coverage.html
.PHONY: test

#
# Modules
#

update: ## Update the go modules.
	@go get -u ./...
	@go mod tidy
.PHONY: update
