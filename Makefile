BIN     ?= tmpl
IMAGE   ?= ghcr.io/jeremybower/tmpl
SUFFIX	?=
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
# Build
#

build: ## Build command.
build:
	@if [ -z "${BIN}" ]; then echo "BIN is not set"; exit 1; fi
	@if [ -z "${VERSION}" ]; then echo "VERSION is not set"; exit 1; fi
	@mkdir -p bin
	@go build -ldflags="-X 'main.Version=${VERSION}'" -o bin/${BIN} .
.PHONY: build

#
# Clean
#

clean: ## Clean working directory.
clean:
	@rm -rf bin
	@rm -rf coverage
.PHONY: clean

#
# Image
#

define DOCKERFILE
FROM scratch

COPY bin/${BIN} /tmpl

ENTRYPOINT ["/tmpl"]
endef
export DOCKERFILE

image: ## Build Docker image.
image:
	@if [ -z "${BIN}" ]; then echo "BIN is not set"; exit 1; fi
	@if [ ! -f bin/${BIN} ]; then echo "bin/${BIN} does not exist"; exit 1; fi
	@if [ -z "${IMAGE}" ]; then echo "IMAGE is not set"; exit 1; fi
	@if [ -z "${VERSION}" ]; then echo "VERSION is not set"; exit 1; fi
	@echo "$$DOCKERFILE" > Dockerfile
	@docker build \
		--platform ${TARGETPLATFORM} \
		-t ${IMAGE}:latest${SUFFIX} \
		-t ${IMAGE}:${VERSION}${SUFFIX} \
		.
.PHONY: image

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
