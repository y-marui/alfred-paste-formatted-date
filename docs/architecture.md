# Architecture

## Overview

An Alfred Workflow (Go): `cmd/paste-formatted-date-alfred` is the single
universal (amd64+arm64) binary `workflow/info.plist` invokes. Its Script
Filter node passes the query following the `date` keyword as `$1`;
`internal/datecmd.Dispatch` parses it, resolves the target date via
`internal/dateresolve`, formats it via `internal/datefmt`, and prints Alfred
Script Filter JSON via `internal/scriptfilter`. Selecting a result copies its
value to the clipboard and auto-pastes it through Alfred's own native
Clipboard Output node — no script is involved in that step.
`scripts/build-workflow.sh` packages the binary with `workflow/info.plist`
and `workflow/icon.png` into a `.alfredworkflow`.

This structure — a thin `cmd/` entry point over independently testable
`internal/` packages, a single dispatch function instead of a generic
command-router abstraction, Script Filter JSON via a small `scriptfilter`
package — deliberately matches
[y-marui/alfred-clean-invisible-text](https://github.com/y-marui/alfred-clean-invisible-text)
and [y-marui/alfred-password-generator](https://github.com/y-marui/alfred-password-generator),
this author's other Alfred Workflows already implemented in Go. This
workflow itself was originally a Python implementation
(`src/alfred`/`src/app`, following the `alfred-workflow-template` scaffold);
see `CHANGELOG.md`'s `[Unreleased]` entry for what changed and why in that
rewrite.

## Entry Points

- `cmd/paste-formatted-date-alfred` — a single command, no subcommands. The
  query it receives (e.g. `""`, `"-2d ISO"`, `"config reset"`, `"help"`)
  determines behavior — see `internal/datecmd`'s package doc comment for the
  full command list.

One Alfred trigger reaches it: the `date` keyword, wired in
`workflow/info.plist`.

## Directory Structure

| Directory | Role |
|---|---|
| `cmd/paste-formatted-date-alfred/` | The binary Alfred invokes; recovers panics into a Script Filter error item and writes the response |
| `internal/datecmd/` | Query dispatch (`date` / `config` / `help`) — builds the Alfred result rows |
| `internal/dateresolve/` | Parses a query into a target date and remaining format filter (relative offsets, direct dates), unit tested independently of Alfred |
| `internal/datefmt/` | The list of selectable date formats and their rendering |
| `internal/wfconfig/` | Persistent key/value store for the `config` command, backed by a JSON file in Alfred's workflow data directory |
| `internal/scriptfilter/` | Alfred Script Filter JSON response types |
| `workflow/` | `info.plist` (the Alfred object graph), `icon.png` |
| `scripts/build-workflow.sh` | Builds the universal binary and packages `workflow/` into `dist/*.alfredworkflow` |
| `scripts/extract-changelog.sh` | Extracts one version's notes from `CHANGELOG.md` for GitHub Releases |
| `docs/` | Architecture and Configuration Builder reference |
| `docs/dev-charter/` | Shared dev-charter (`git subtree`) |

## Query Parsing

`internal/datecmd.Dispatch` splits the query into `<command> <rest>` on the
first whitespace run:

```
""                →  command=""       →  date command, today, all formats
"-2d ISO"         →  command="-2d"    →  unrecognized → date command, full query as args
"config reset"    →  command="config" →  config command, args="reset"
"help"            →  command="help"   →  help command
```

Only `config`, `help`, and the explicit literal `date` are registered
commands; anything else falls back to the date command with the whole
trimmed query as its args (so a bare filter like `"ISO"` or a relative
offset like `"-2d"` works without a command prefix).

## Key Dependencies

None. Every `internal/` package uses only the Go standard library
(`time`, `regexp`, `encoding/json`, etc.).

## Alfred Configuration Builder (`userconfigurationconfig`)

Alfred 5 の Configuration Builder は `info.plist` の `userconfigurationconfig` キーで定義する。
利用可能な全型・各キーの詳細は [`docs/configuration-builder.md`](configuration-builder.md) を参照。

This workflow currently declares no Configuration Builder variables — the
Python predecessor's `Use uv` and `Log Level` settings were scaffold
leftovers from `alfred-workflow-template` with no equivalent need in the Go
binary (no interpreter to select, no file logging).
