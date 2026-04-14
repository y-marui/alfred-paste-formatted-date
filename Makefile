.PHONY: help install hooks lint format typecheck test test-cov vendor build deploy release run clean update-charter

# ---------------------------------------------------------------------------
# Python environment selector
#
#   use_uv=1  (default when uv is found) → uv managed virtual environment
#   use_uv=0  (default when uv is absent) → global python3 / pip3
#
# Auto-detected from PATH; override explicitly if needed:
#   make test use_uv=0   # force pip3 even if uv is installed
#   make test use_uv=1   # force uv
# ---------------------------------------------------------------------------
use_uv ?= $(shell command -v uv >/dev/null 2>&1 && echo 1 || echo 0)

ifeq ($(use_uv),1)
  PYTHON := uv run python
  RUN    := uv run
else
  PYTHON := python3
  RUN    :=
endif

# Export so child shell scripts (vendor.sh, build.sh, dev.sh) inherit the flag
export use_uv

# Default target
help:
	@echo "Alfred Workflow Template - Development Commands"
	@echo ""
	@echo "  make install     Install dev dependencies"
	@echo "  make hooks       Install pre-commit hooks"
	@echo "  make lint        Run ruff linter"
	@echo "  make format      Auto-format with ruff"
	@echo "  make typecheck   Run mypy type checker"
	@echo "  make test        Run tests"
	@echo "  make test-cov    Run tests with coverage report"
	@echo "  make vendor      Install runtime deps into workflow/vendor/"
	@echo "  make build       Build .alfredworkflow package"
	@echo "  make deploy      Build and install into Alfred"
	@echo "  make release     Create GitHub Release (requires git tag)"
	@echo "  make run Q=''    Simulate Alfred with query Q"
	@echo "  make clean       Remove build artifacts"
	@echo "  make update-charter  Pull latest dev-charter via git subtree"
	@echo ""
	@echo "  use_uv=1 (default when uv found) → uv managed virtual environment"
	@echo "  use_uv=0 (default when uv absent) → global python3 / pip3"
	@echo ""

# ---------------------------------------------------------------------------
# Setup
# ---------------------------------------------------------------------------
install:
ifeq ($(use_uv),1)
	uv sync --extra dev
	@$(MAKE) hooks
else
	pip3 install --quiet -e ".[dev]"
	@command -v pre-commit >/dev/null 2>&1 && $(MAKE) hooks || true
endif

hooks:
ifeq ($(use_uv),1)
	uv run pre-commit install
else
	pip3 install --quiet pre-commit
	pre-commit install
endif

# ---------------------------------------------------------------------------
# Code quality
# ---------------------------------------------------------------------------
lint:
	$(RUN) ruff check src/ tests/

format:
	$(RUN) ruff format src/ tests/
	$(RUN) ruff check --fix src/ tests/

typecheck:
	$(RUN) mypy src/

# ---------------------------------------------------------------------------
# Tests
# ---------------------------------------------------------------------------
test:
	$(RUN) pytest

test-cov:
	$(RUN) pytest --cov=src --cov-report=term-missing --cov-report=html

# ---------------------------------------------------------------------------
# Build
# ---------------------------------------------------------------------------
vendor:
	./scripts/vendor.sh

build: vendor
	./scripts/build.sh

deploy: build
	@open dist/*.alfredworkflow

release:
	./scripts/release.sh

# ---------------------------------------------------------------------------
# Local dev
# ---------------------------------------------------------------------------
run:
	@./scripts/dev.sh "$(Q)"

# ---------------------------------------------------------------------------
# Maintenance
# ---------------------------------------------------------------------------
clean:
	rm -rf .build dist/ .coverage htmlcov/ .mypy_cache/ .pytest_cache/
	find . -name "__pycache__" -type d -exec rm -rf {} + 2>/dev/null || true
	find . -name "*.pyc" -delete

# ---------------------------------------------------------------------------
# Dev Charter
# ---------------------------------------------------------------------------
update-charter:
	git subtree pull --prefix=docs/dev-charter dev-charter main --squash
