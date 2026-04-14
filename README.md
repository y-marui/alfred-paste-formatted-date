# Paste Formatted Date

> **This is the English (reference) version.**
> For the Japanese canonical version, see [README-jp.md](README-jp.md).

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](LICENSE)
[![CI](https://github.com/y-marui/alfred-paste-formatted-date/actions/workflows/ci.yml/badge.svg)](https://github.com/y-marui/alfred-paste-formatted-date/actions/workflows/ci.yml)
[![Charter Check](https://github.com/y-marui/alfred-paste-formatted-date/actions/workflows/dev-charter-check.yml/badge.svg)](https://github.com/y-marui/alfred-paste-formatted-date/actions/workflows/dev-charter-check.yml)

Generate and paste today's date in multiple formats via Alfred 5.

## Usage

Type `date` in Alfred to see all available formats. Select one to copy and auto-paste it.

```
date             — list all formats
date <filter>    — filter by format name or value (e.g. "ISO", "YYYY", "unix")
date config      — view or reset configuration
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
- Python 3.9+

## Installation

Download the latest `.alfredworkflow` from [Releases](https://github.com/y-marui/alfred-paste-formatted-date/releases) and double-click to install.

## Development

```bash
# Install dev dependencies
make install

# Simulate Alfred locally
make run Q=""
make run Q="ISO"

# Run tests
make test

# Build workflow package
make build
# → dist/alfred-paste-formatted-date-0.1.0.alfredworkflow
```

## Project Structure

```
alfred-paste-formatted-date/
├── src/
│   ├── alfred/         # Alfred SDK (response, router, cache, config, logger, safe_run)
│   └── app/            # Application layer (commands)
├── workflow/           # Alfred package (info.plist, scripts/entry.py, vendor/)
├── tests/              # pytest test suite
├── scripts/            # build.sh, dev.sh, release.sh, vendor.sh
└── docs/               # Architecture and development documentation
```

## Support

If this workflow saves you time, support is appreciated.

- [Buy Me a Coffee](https://www.buymeacoffee.com/y.marui)
- [GitHub Sponsors](https://github.com/sponsors/y-marui)

## License

MIT — see [LICENSE](LICENSE)

---

*This is the reference (English) version. The canonical Japanese version is [README-jp.md](README-jp.md). Update both files in the same commit.*
