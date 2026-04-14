#!/usr/bin/env bash
# Local development runner - simulates Alfred's Script Filter execution.
#
# Usage:
#   ./scripts/dev.sh "search foo"
#   ./scripts/dev.sh "config"
#   ./scripts/dev.sh ""
#
# Inherits use_uv from the environment (set by Makefile or caller).
# When not set, auto-detects uv: ON if uv is found, OFF otherwise.
#   use_uv=1 (default when uv found) → uv run python
#   use_uv=0                         → python3
#
# Output is pretty-printed JSON if `jq` is available.
set -euo pipefail

REPO_ROOT="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ENTRY="$REPO_ROOT/workflow/scripts/entry.py"

QUERY="${1:-}"

# Auto-detect uv when use_uv is not explicitly set
if [[ -z "${use_uv+set}" ]]; then
  command -v uv >/dev/null 2>&1 && use_uv=1 || use_uv=0
fi

# Make src/ importable without `pip install -e .`
# In the packaged workflow, entry.py adds workflow_root/src/ to sys.path.
# During development, workflow/src/ does not exist – the real source is at
# repo_root/src/, so we set PYTHONPATH explicitly.
export PYTHONPATH="$REPO_ROOT/src${PYTHONPATH:+:$PYTHONPATH}"

# Set environment variables that Alfred normally provides
export alfred_workflow_bundleid="${alfred_workflow_bundleid:-com.example.dev}"
export alfred_workflow_cache="${alfred_workflow_cache:-/tmp/alfred-dev-cache}"
export alfred_workflow_data="${alfred_workflow_data:-/tmp/alfred-dev-data}"
export alfred_workflow_version="${alfred_workflow_version:-0.0.0-dev}"
export alfred_workflow_uid="${alfred_workflow_uid:-com.example.dev}"
export alfred_workflow_name="${alfred_workflow_name:-Workflow Dev}"

mkdir -p "$alfred_workflow_cache" "$alfred_workflow_data"

# Select Python command
if [[ "$use_uv" == "1" ]]; then
  PYTHON_CMD="uv run python"
else
  PYTHON_CMD="python3"
fi

echo "─────────────────────────────────────"
echo "  Alfred Script Filter Simulator"
echo "  Query: \"$QUERY\""
if [[ "$use_uv" == "1" ]]; then
  echo "  Python: uv venv"
else
  echo "  Python: $(python3 --version 2>&1)"
fi
echo "─────────────────────────────────────"

if command -v jq &>/dev/null; then
  $PYTHON_CMD "$ENTRY" "$QUERY" | jq .
else
  $PYTHON_CMD "$ENTRY" "$QUERY"
fi
