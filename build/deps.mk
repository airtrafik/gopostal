#
# Makefile fragment for installing deps
#
# Tools are executed via `go run pkg@version` which downloads on-demand.
# Tool versions are centralized in build/versions.mk
#

GO           ?= go

deps: deps-only

deps-only:
	@echo "=== $(PROJECT_NAME) === [ deps             ]: Installing package dependencies required by the project..."
	@$(GO) mod tidy
	@$(GO) mod download

.PHONY: deps deps-only
