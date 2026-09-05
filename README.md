# Paste Formatted Date

> **This is the English (reference) version.**
> For the Japanese canonical version, see [README-jp.md](README-jp.md).

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![CI](https://github.com/y-marui/alfred-paste-formatted-date/actions/workflows/ci.yml/badge.svg)](https://github.com/y-marui/alfred-paste-formatted-date/actions/workflows/ci.yml)
[![Charter Check](https://github.com/y-marui/alfred-paste-formatted-date/actions/workflows/dev-charter-check.yml/badge.svg)](https://github.com/y-marui/alfred-paste-formatted-date/actions/workflows/dev-charter-check.yml)
[![GitHub Sponsors](https://img.shields.io/github/sponsors/y-marui?style=social)](https://github.com/sponsors/y-marui)
[![Buy Me a Coffee](https://img.shields.io/badge/Buy%20Me%20a%20Coffee-donate-yellow.svg)](https://www.buymeacoffee.com/y.marui)

Generate and paste today's date in multiple formats via Alfred 5.

## Usage

List today's date in multiple formats via the `date` keyword, then select
one to copy and auto-paste it.

```
date             — list all formats
date <filter>    — filter by format name or value (e.g. "ISO", "YYYY", "unix")
date help        — show available commands
```

### Available formats

| Format | Example |
|---|---|
| YYYYMMDD | 20260414 |
| YYMMDD | 260414 |
| YYYY-MM-DD | 2026-04-14 |
| YYYY/MM/DD | 2026/04/14 |
| MM/DD/YYYY | 04/14/2026 |
| DD/MM/YYYY | 14/04/2026 |
| MMM DD, YYYY | Apr 14, 2026 |
| MMMM DD, YYYY | April 14, 2026 |
| YYYY-MM-DDThh:mm:ss | 2026-04-14T12:00:00 |
| Unix timestamp | 1744588800 |

## Requirements

- Alfred 5 (Powerpack required for Script Filter)

## Installation

Download the latest `.alfredworkflow` from [Releases](https://github.com/y-marui/alfred-paste-formatted-date/releases) and double-click to install.

## Development

Requires Go (see `go.mod` for the toolchain version) — see [DEVELOPING.md](DEVELOPING.md) for the full workflow.

```bash
# Simulate Alfred locally
go run ./cmd/paste-formatted-date-alfred ""
go run ./cmd/paste-formatted-date-alfred "ISO"

# Run tests
make test

# Build workflow package
make build-workflow
# → dist/paste-formatted-date-0.1.0.alfredworkflow
```

## Project Structure

```
alfred-paste-formatted-date/
├── cmd/
│   └── paste-formatted-date-alfred/  # The binary Alfred invokes
├── internal/
│   ├── datecmd/         # Command dispatch (date / help)
│   ├── dateresolve/     # Query → target date parsing
│   ├── datefmt/         # Date format table + rendering
│   └── scriptfilter/    # Alfred Script Filter JSON types
├── workflow/            # Alfred package (info.plist, icon.png)
└── docs/                # Architecture documentation
```

## License

MIT — see [LICENSE](LICENSE)

---

*This is the reference (English) version. The canonical Japanese version is [README-jp.md](README-jp.md). Update both files in the same commit.*
