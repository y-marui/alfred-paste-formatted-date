# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### Changed

- Rewrote the implementation from Python to Go, matching the `cmd/`+`internal/`
  architecture already used in the sibling `alfred-clean-invisible-text` and
  `alfred-password-generator` workflows. No user-facing behavior change: the
  `date`/`config`/`help` commands, date formats, and relative/direct date
  parsing are ported 1:1 (see the previous Python implementation at
  [`src/alfred`/`src/app`](https://github.com/y-marui/alfred-paste-formatted-date/tree/a237369/src)).
- The packaged workflow now ships a single universal (amd64+arm64) compiled
  binary instead of a Python runtime + vendored packages, invoked directly by
  `workflow/info.plist` with no interpreter selection.
- Dropped the unused `Cache` SDK module and the `Log Level`/`Use uv`
  Configuration Builder settings (file logging via `alfred/logger.py`) — no
  Go sibling implements these, and neither was ever exercised by the date
  command behavior.

## [0.1.0] - 2024-01-01

### Added

- Initial release of the Alfred Workflow Template
- Alfred SDK: `response`, `cache`, `config`, `logger`, `router`, `safe_run`
- Command-based UX: `search`, `open`, `config`, `help`
- Vendor packaging via `scripts/vendor.sh`
- Build pipeline via `scripts/build.sh`
- GitHub Actions CI (lint, test, build)
- GitHub Actions Release (tag → `.alfredworkflow` → GitHub Release)
- Full pytest test suite

[Unreleased]: https://github.com/y-marui/alfred-paste-formatted-date/compare/a237369...HEAD
[0.1.0]: https://github.com/y-marui/alfred-paste-formatted-date/commit/a237369
