#
# Makefile fragment for Testing
#
# Uses `go run pkg@version` to execute tools without requiring pre-installation.
# Tool versions are defined in build/versions.mk
#
# Note: This is a CGo library requiring libpostal to be installed.
# Tests cannot run without the native library and its data files.
#

GO           ?= go

COVERAGE_DIR ?= ./coverage
COVERMODE    ?= atomic
SRCDIR       ?= .
GO_PKGS      ?= $(shell $(GO) list ./... | grep -v -e "/vendor/" -e "/example")

test: test-unit

test-unit: deps-only
	@echo "=== $(PROJECT_NAME) === [ test-unit        ]: running unit tests..."
	@mkdir -p $(COVERAGE_DIR)
	@$(GO) run $(GOTESTSUM_PKG) -f testname --junitfile $(COVERAGE_DIR)/unit.xml -- -v -covermode=$(COVERMODE) -coverprofile $(COVERAGE_DIR)/unit.tmp $(GO_PKGS)

#
# Coverage
#
cover-clean:
	@echo "=== $(PROJECT_NAME) === [ cover-clean      ]: removing coverage files..."
	@rm -rfv $(COVERAGE_DIR)/*

cover-report:
	@echo "=== $(PROJECT_NAME) === [ cover-report     ]: generating coverage results..."
	@mkdir -p $(COVERAGE_DIR)
	@echo 'mode: $(COVERMODE)' > $(COVERAGE_DIR)/coverage.out
	@cat $(COVERAGE_DIR)/*.tmp | grep -v 'mode: $(COVERMODE)' >> $(COVERAGE_DIR)/coverage.out || true
	@$(GO) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "=== $(PROJECT_NAME) === [ cover-report     ]:     $(COVERAGE_DIR)/coverage.html"

cover-view: cover-report
	@$(GO) tool cover -html=$(COVERAGE_DIR)/coverage.out

.PHONY: test test-unit cover-clean cover-report cover-view
