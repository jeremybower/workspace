#
# Help
#

.PHONY: help
help: ## Show this help message.
	@echo
	@echo 'usage: make [target]'
	@echo
	@echo 'targets:'
	@echo
	@egrep '^(.+)\:\ ##\ (.+)' ${MAKEFILE_LIST} | column -t -c 2 -s ':#'
	@echo

#
# Build Targets
#

OUT ?= bin/tmpl
VERSION ?= v0.0.0

.PHONY: build
build: ## Build the command.
build:
	@mkdir -p bin
	@go build -ldflags="-X 'main.Version=${VERSION}'" -o ${OUT} .

#
# Clean
#

.PHONY: clean
clean: ## Clean the working directory.
clean:
	@rm -rf bin
	@rm -rf coverage

#
# Test Targets
#

.PHONY: test
test: ## Run the tests.
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
