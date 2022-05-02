TAG_COMMIT := $(shell git rev-list --abbrev-commit --tags --max-count=1)
# `2>/dev/null` suppress errors and `|| true` suppress the error codes.
TAG := $(shell git describe --abbrev=0 --tags ${TAG_COMMIT} 2>/dev/null || true)
# here we strip the version prefix
VERSION := $(TAG:v%=%)


help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

release: ## Build a release version
	@go build -tags osusergo,netgo -o certbot-govh
	@mkdir releases
	@tar -xjf releases/certbot-govh-$(VERSION) certbot-govh certbot-govh-auth.sh certbot-govh-cleanup.sh
	@echo "Prod release done"

.PHONY: help
.DEFAULT_GOAL := help