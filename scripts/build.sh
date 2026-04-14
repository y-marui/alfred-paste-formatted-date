#!/usr/bin/env bash
# Build the .alfredworkflow package.
#
# Steps:
#   1. Read version from pyproject.toml
#   2. Sync version into workflow/info.plist
#   3. Copy src/ into the build dir
#   4. Ensure vendor/ is up to date
#   5. Zip into dist/<name>-<version>.alfredworkflow
#
# Inherits use_uv from the environment (set by Makefile or caller):
#   use_uv=0 (default) → python3 / pip3
#   use_uv=1           → uv run python / uv pip
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
WORKFLOW_DIR="$REPO_ROOT/workflow"
SRC_DIR="$REPO_ROOT/src"
DIST_DIR="$REPO_ROOT/dist"
BUILD_DIR="$REPO_ROOT/.build"

# Select Python interpreter based on use_uv flag
if [[ "${use_uv:-0}" == "1" ]]; then
  PYTHON="uv run python"
  PIP_INSTALL="uv pip install"
else
  PYTHON="python3"
  PIP_INSTALL="env PIP_REQUIRE_VIRTUALENV=false pip3 install --quiet --upgrade"
fi

# ---------------------------------------------------------------------------
# 1. Read version from pyproject.toml
#    Delegated to scripts/read_version.py to avoid heredoc/quote conflicts.
# ---------------------------------------------------------------------------
VERSION=$($PYTHON scripts/read_version.py)
echo "→ Version: $VERSION"

# ---------------------------------------------------------------------------
# 2. Sync version into workflow/info.plist
#    Use PlistBuddy for a surgical edit — avoids round-tripping the entire
#    file through plistlib, which can corrupt Alfred Configuration Builder data.
# ---------------------------------------------------------------------------
/usr/libexec/PlistBuddy -c "Set :version ${VERSION}" "$WORKFLOW_DIR/info.plist"
echo "  Synced info.plist version → ${VERSION}"

# ---------------------------------------------------------------------------
# 3. Prepare build directory
# ---------------------------------------------------------------------------
echo "→ Preparing build directory"
rm -rf "$BUILD_DIR"
cp -r "$WORKFLOW_DIR/" "$BUILD_DIR/"

# Copy src/ into the build dir (available as src/ at workflow root)
echo "→ Copying src/"
cp -r "$SRC_DIR/" "$BUILD_DIR/src/"

# ---------------------------------------------------------------------------
# 4. Install vendor dependencies
# ---------------------------------------------------------------------------
echo "→ Installing vendor dependencies"
mkdir -p "$BUILD_DIR/vendor"

if [[ -f "$REPO_ROOT/requirements.txt" ]]; then
  $PIP_INSTALL \
    --requirement "$REPO_ROOT/requirements.txt" \
    --target "$BUILD_DIR/vendor"
else
  echo "  No requirements.txt - skipping"
fi

# Remove development artifacts to reduce package size
find "$BUILD_DIR" -name "*.pyc" -delete
find "$BUILD_DIR" -name "__pycache__" -type d -print0 | xargs -0 rm -rf
find "$BUILD_DIR" -name "*.egg-info" -print0 | xargs -0 rm -rf
find "$BUILD_DIR" -name "*.dist-info" -print0 | xargs -0 rm -rf

# ---------------------------------------------------------------------------
# 5. Package into .alfredworkflow
# ---------------------------------------------------------------------------
mkdir -p "$DIST_DIR"

# Read workflow name from info.plist for the output filename
WORKFLOW_NAME=$($PYTHON - "$BUILD_DIR/info.plist" <<'EOF'
import sys, plistlib, pathlib, re
with pathlib.Path(sys.argv[1]).open("rb") as f:
    data = plistlib.load(f)
name = data.get("name", "workflow")
# Sanitize for filename
print(re.sub(r"[^\w\-.]", "-", name).lower())
EOF
)

OUTPUT="$DIST_DIR/${WORKFLOW_NAME}-${VERSION}.alfredworkflow"

rm -f "$OUTPUT"  # ensure a clean zip (zip -r updates rather than replaces)
echo "→ Packaging: $OUTPUT"
(cd "$BUILD_DIR" && zip -r "$OUTPUT" . -x "*.DS_Store" -x ".git/*" --quiet)

echo "✓ Build complete: $OUTPUT"
