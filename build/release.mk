#
# Makefile fragment for Releases
#
# Semantic versioning based on conventional commits using svu.
# Creates GitHub releases with auto-generated notes via gh CLI.
# No goreleaser, no changelog files — this is a library.
#
# Targets:
#   release-preview  - Show next version and commits without making changes
#   release          - Interactive release with confirmation prompt
#   release-auto     - Non-interactive release (for CI pipelines)
#   release-tag      - Create and push a release tag manually
#   release-clean    - Clean distribution directory
#

GO           ?= go
DIST_DIR     ?= ./dist

# Tool commands using centralized versions from versions.mk
SVU_CMD          := $(GO) run $(SVU_PKG)

# Default branch for releases
DEFAULT_BRANCH   ?= main

# Preview the next release without making changes
release-preview:
	@echo "=== $(PROJECT_NAME) === [ release-preview  ]: Calculating next version from conventional commits..."
	@echo ""
	@CURRENT=$$($(SVU_CMD) current --strip-prefix 2>/dev/null || echo "0.0.0"); \
	NEXT=$$($(SVU_CMD) next --strip-prefix 2>/dev/null || echo "0.1.0"); \
	echo "  Current version: v$$CURRENT"; \
	echo "  Next version:    v$$NEXT"; \
	echo ""; \
	if [ "$$CURRENT" = "$$NEXT" ]; then \
		echo "  Status: No version bump needed (no feat/fix commits since last tag)"; \
	else \
		echo "  Commits since v$$CURRENT:"; \
		git log --oneline v$$CURRENT..HEAD 2>/dev/null | head -20 | sed 's/^/    /'; \
		TOTAL=$$(git rev-list --count v$$CURRENT..HEAD 2>/dev/null || echo "0"); \
		if [ "$$TOTAL" -gt 20 ]; then \
			echo "    ... and $$((TOTAL - 20)) more commits"; \
		fi; \
	fi
	@echo ""

# Interactive release with confirmation prompt
release: release-check
	@echo "=== $(PROJECT_NAME) === [ release          ]: Starting release process..."
	@CURRENT=$$($(SVU_CMD) current --strip-prefix 2>/dev/null || echo "0.0.0"); \
	NEXT=$$($(SVU_CMD) next --strip-prefix 2>/dev/null || echo "0.1.0"); \
	if [ "$$CURRENT" = "$$NEXT" ]; then \
		echo ""; \
		echo "  No version bump needed - no feat/fix commits since v$$CURRENT"; \
		echo "  To force a release, use: make release-tag VERSION=x.y.z"; \
		exit 0; \
	fi; \
	echo ""; \
	echo "  Current version: v$$CURRENT"; \
	echo "  Next version:    v$$NEXT"; \
	echo ""; \
	echo "  This will:"; \
	echo "    1. Create and push tag v$$NEXT"; \
	echo "    2. Create GitHub release with auto-generated notes"; \
	echo ""; \
	read -p "  Proceed with release? [y/N] " confirm; \
	if [ "$$confirm" != "y" ] && [ "$$confirm" != "Y" ]; then \
		echo "  Release cancelled."; \
		exit 0; \
	fi; \
	$(MAKE) release-execute VERSION=$$NEXT

# Non-interactive release for CI pipelines
release-auto: release-check
	@echo "=== $(PROJECT_NAME) === [ release-auto     ]: Starting automated release..."
	@CURRENT=$$($(SVU_CMD) current --strip-prefix 2>/dev/null || echo "0.0.0"); \
	NEXT=$$($(SVU_CMD) next --strip-prefix 2>/dev/null || echo "0.1.0"); \
	if [ "$$CURRENT" = "$$NEXT" ]; then \
		echo "  No version bump needed - no feat/fix commits since v$$CURRENT"; \
		exit 0; \
	fi; \
	echo "  Releasing v$$NEXT (current: v$$CURRENT)"; \
	$(MAKE) release-execute VERSION=$$NEXT

# Pre-flight checks for release
release-check:
	@echo "=== $(PROJECT_NAME) === [ release-check    ]: Running pre-flight checks..."
	@BRANCH=$$(git rev-parse --abbrev-ref HEAD); \
	if [ "$$BRANCH" != "$(DEFAULT_BRANCH)" ]; then \
		echo "  ERROR: Must be on $(DEFAULT_BRANCH) branch (currently on $$BRANCH)"; \
		exit 1; \
	fi
	@if [ -n "$$(git status --porcelain)" ]; then \
		echo "  ERROR: Working directory has uncommitted changes"; \
		git status --short; \
		exit 1; \
	fi
	@GIT_USER=$$(git config user.name); \
	if [ -z "$$GIT_USER" ]; then \
		echo "  ERROR: git user.name not configured"; \
		exit 1; \
	fi
	@GIT_EMAIL=$$(git config user.email); \
	if [ -z "$$GIT_EMAIL" ]; then \
		echo "  ERROR: git user.email not configured"; \
		exit 1; \
	fi
	@if ! gh auth token >/dev/null 2>&1; then \
		echo "  ERROR: GitHub CLI not authenticated (run 'gh auth login')"; \
		exit 1; \
	fi
	@echo "  All checks passed"

# Execute the release (called by release and release-auto)
release-execute:
	@if [ -z "$(VERSION)" ]; then \
		echo "  ERROR: VERSION not specified"; \
		exit 1; \
	fi
	@echo "=== $(PROJECT_NAME) === [ release-execute  ]: Creating tag v$(VERSION)..."
	@git tag v$(VERSION)
	@git push --no-verify origin HEAD:$(DEFAULT_BRANCH) --tags
	@echo "=== $(PROJECT_NAME) === [ release-execute  ]: Creating GitHub release..."
	@gh release create v$(VERSION) --generate-notes --latest
	@echo ""
	@echo "=== $(PROJECT_NAME) === [ release-execute  ]: Release v$(VERSION) complete!"

# Create and push a new release tag manually (bypasses svu calculation)
release-tag:
	@if [ -z "$(VERSION)" ]; then \
		echo "ERROR: VERSION not specified. Usage: make release-tag VERSION=0.1.0"; \
		exit 1; \
	fi
	@echo "=== $(PROJECT_NAME) === [ release-tag      ]: Creating tag v$(VERSION)..."
	@git tag -a v$(VERSION) -m "Release v$(VERSION)"
	@git push origin v$(VERSION)
	@echo "=== $(PROJECT_NAME) === [ release-tag      ]: Tag v$(VERSION) created and pushed"

# Clean distribution directory
release-clean:
	@echo "=== $(PROJECT_NAME) === [ release-clean    ]: Cleaning distribution files..."
	@rm -rf $(DIST_DIR)

.PHONY: release release-auto release-check release-execute release-preview \
        release-tag release-clean
