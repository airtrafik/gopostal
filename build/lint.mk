
#
# Makefile fragment for Linting
#
# Uses `go run pkg@version` to execute tools without requiring pre-installation.
# Tool versions are defined in build/versions.mk
#

GO           ?= go
GOFMT        ?= gofmt

EXCLUDEDIR      ?= .git
SRCDIR          ?= .
GO_PKGS         ?= $(shell ${GO} list ./... | grep -v -e "/vendor/" -e "/example")
FILES           ?= $(shell find ${SRCDIR} -type f \( -name '*.go' -o -name '*.yml' -o -name '*.json' -o -name '*.mk' \))
GO_FILES        ?= $(shell find $(SRCDIR) -type f -name "*.go" | grep -v -e ".git/" -e '/vendor/' -e '/example/')
PROJECT_MODULE  ?= $(shell $(GO) list -m)

lint: deps-only spell-check gofmt golangci goimports
lint-fix: deps-only gofmt-fix goimports

#
# Check spelling on all the files, not just source code
#
spell-check: deps-only
	@echo "=== $(PROJECT_NAME) === [ spell-check      ]: Checking for spelling mistakes..."
	@$(GO) run $(MISSPELL_PKG) -source text $(FILES)

spell-check-fix: deps-only
	@echo "=== $(PROJECT_NAME) === [ spell-check-fix  ]: Fixing spelling mistakes..."
	@$(GO) run $(MISSPELL_PKG) -source text -w $(FILES)

gofmt:
	@echo "=== $(PROJECT_NAME) === [ gofmt            ]: Checking file format with $(GOFMT)..."
	@$(GOFMT) -e -l -s -d $(GO_FILES)

gofmt-fix:
	@echo "=== $(PROJECT_NAME) === [ gofmt-fix        ]: Fixing file format with $(GOFMT)..."
	@$(GOFMT) -e -l -s -w $(GO_FILES)

goimports: deps-only
	@echo "=== $(PROJECT_NAME) === [ goimports        ]: Checking/fixing imports..."
	@$(GO) run $(GOIMPORTS_PKG) -l -w -local $(PROJECT_MODULE) $(GO_FILES)

golangci: deps-only
	@echo "=== $(PROJECT_NAME) === [ golangci-lint    ]: Linting..."
	@$(GO) run $(GOLANGCI_PKG) run --allow-serial-runners

outdated: deps-only
	@echo "=== $(PROJECT_NAME) === [ outdated         ]: Finding outdated deps..."
	@$(GO) list -u -m -json -mod=mod all | $(GO) run $(GO_MOD_OUTDATED_PKG) -direct -update

.PHONY: lint spell-check spell-check-fix gofmt gofmt-fix lint-fix goimports golangci outdated
