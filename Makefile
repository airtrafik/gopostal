#############################
# Global vars
#############################
PROJECT_NAME := $(shell basename $(shell pwd))
PROJECT_VER  ?= $(shell git describe --tags --always --dirty | sed -e '/^v/s/^v\(.*\)$$/\1/g')

SRCDIR       ?= .
GO            = go

# The root module (from go.mod)
PROJECT_MODULE  ?= $(shell $(GO) list -m)

#############################
# Targets
#############################
all: build

# Humans running make:
build: clean lint test

# All clean commands
clean: cover-clean release-clean

# Import fragments (versions.mk must be first)
include build/versions.mk
include build/deps.mk
include build/lint.mk
include build/release.mk
include build/test.mk

help:
	@echo ""
	@echo "gopostal Makefile — v$(PROJECT_VER)"
	@echo ""
	@echo "Build"
	@echo "  make                          Full build (clean, lint, test)"
	@echo "  make clean                    Remove coverage and dist files"
	@echo ""
	@echo "Lint"
	@echo "  make lint                     Run all linters"
	@echo "  make lint-fix                 Auto-fix formatting and imports"
	@echo ""
	@echo "Test"
	@echo "  make test                     Run unit tests"
	@echo "  make test-unit                Run unit tests"
	@echo "  make cover-report             Generate HTML coverage report"
	@echo "  make cover-view               Open coverage report in browser"
	@echo ""
	@echo "Release"
	@echo "  make release-preview          Show next version and commits"
	@echo "  make release                  Interactive release with confirmation"
	@echo "  make release-auto             Non-interactive release (CI)"
	@echo "  make release-tag              Create tag manually   (VERSION=x.y.z)"
	@echo ""

.PHONY: all build clean help
