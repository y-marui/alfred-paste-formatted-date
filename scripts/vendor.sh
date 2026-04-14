#!/usr/bin/env bash
# Install runtime dependencies into workflow/vendor/
# Run this after adding packages to requirements.txt.
#
# Inherits use_uv from the environment (set by Makefile or caller):
#   use_uv=0 (default) → pip3
#   use_uv=1           → uv pip
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
VENDOR_DIR="$REPO_ROOT/workflow/vendor"

echo "→ Installing dependencies into $VENDOR_DIR"
mkdir -p "$VENDOR_DIR"

# Clear existing vendor dir to avoid stale packages
rm -rf "${VENDOR_DIR:?}"/*

if [[ ! -f "$REPO_ROOT/requirements.txt" ]]; then
  echo "  No requirements.txt found - skipping vendor install"
  exit 0
fi

if [[ "${use_uv:-0}" == "1" ]]; then
  uv pip install \
    --requirement "$REPO_ROOT/requirements.txt" \
    --target "$VENDOR_DIR"
else
  env PIP_REQUIRE_VIRTUALENV=false pip3 install \
    --quiet \
    --requirement "$REPO_ROOT/requirements.txt" \
    --target "$VENDOR_DIR" \
    --upgrade
fi

echo "✓ Vendor install complete"
