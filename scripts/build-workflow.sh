#!/usr/bin/env bash
# Build the .alfredworkflow package.
#
# Steps:
#   1. Build cmd/paste-formatted-date-alfred as a universal (amd64+arm64)
#      binary via lipo, so the bundle runs natively on both Intel and
#      Apple Silicon.
#   2. Copy workflow/ (info.plist, icon.png) into the build dir.
#   3. Zip into dist/<name>-<version>.alfredworkflow.
set -euo pipefail

REPO_ROOT=$(git rev-parse --show-toplevel)
WORKFLOW_DIR="${REPO_ROOT}/workflow"
DIST_DIR="${REPO_ROOT}/dist"
BUILD_DIR="${REPO_ROOT}/.build"

echo "→ Preparing build directory"
rm -rf "$BUILD_DIR"
cp -r "$WORKFLOW_DIR/" "$BUILD_DIR/"

echo "→ Building universal binary (amd64 + arm64)"
GOOS=darwin GOARCH=amd64 go build -o "${BUILD_DIR}/paste-formatted-date-alfred-amd64" ./cmd/paste-formatted-date-alfred
GOOS=darwin GOARCH=arm64 go build -o "${BUILD_DIR}/paste-formatted-date-alfred-arm64" ./cmd/paste-formatted-date-alfred
lipo -create -output "${BUILD_DIR}/paste-formatted-date-alfred" \
  "${BUILD_DIR}/paste-formatted-date-alfred-amd64" \
  "${BUILD_DIR}/paste-formatted-date-alfred-arm64"
rm "${BUILD_DIR}/paste-formatted-date-alfred-amd64" "${BUILD_DIR}/paste-formatted-date-alfred-arm64"
chmod +x "${BUILD_DIR}/paste-formatted-date-alfred"
lipo -info "${BUILD_DIR}/paste-formatted-date-alfred"

VERSION=$(/usr/libexec/PlistBuddy -c "Print :version" "${BUILD_DIR}/info.plist")
WORKFLOW_NAME=$(/usr/libexec/PlistBuddy -c "Print :name" "${BUILD_DIR}/info.plist" | tr '[:upper:] ' '[:lower:]-')

mkdir -p "$DIST_DIR"
OUTPUT="${DIST_DIR}/${WORKFLOW_NAME}-${VERSION}.alfredworkflow"
rm -f "$OUTPUT" # ensure a clean zip (zip -r updates rather than replaces)

echo "→ Packaging: ${OUTPUT}"
(cd "$BUILD_DIR" && zip -r "$OUTPUT" . -x "*.DS_Store" --quiet)

echo "✓ Build complete: ${OUTPUT}"
