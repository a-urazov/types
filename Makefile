# Go parameters
GOCMD=go
GOLINT=golangci-lint

.PHONY: help lint lint-fix test clean

##@ Development

help: ## Display this help
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n"} /^[a-zA-Z_0-9-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

lint: ## Run golangci-lint
	$(GOLINT) run ./...

lint-fix: ## Run golangci-lint and apply fixes
	$(GOLINT) run --fix ./...

test: ## Run tests
	$(GOCMD) test ./... -v

clean: ## Clean build artifacts
	$(GOCMD) clean