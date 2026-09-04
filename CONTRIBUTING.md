# Contributing

Thank you for contributing!

## Before You Start

- Check existing issues and PRs to avoid duplicate work.
- For large changes, open an issue first to discuss the approach.

For development setup, workflow, naming conventions, and code review guidelines,
see [DEVELOPING.md](DEVELOPING.md).

## Release Process

```bash
# 1. Update version in workflow/info.plist
# 2. Update CHANGELOG.md
git add workflow/info.plist CHANGELOG.md
git commit -m "chore: release v1.2.3"

# 3. Tag and push
git tag v1.2.3
git push origin main --tags
# GitHub Actions builds .alfredworkflow and creates a GitHub Release
```

## Security

### Supported Versions

Only the latest release is supported with security fixes.

### Reporting a Vulnerability

Please **do not** open a public GitHub issue for security vulnerabilities.

Instead, report them privately via
[GitHub Security Advisories](https://github.com/y-marui/alfred-paste-formatted-date/security/advisories/new)
or email the maintainer directly.

We aim to acknowledge reports within 48 hours and provide a fix within 7 days
for confirmed vulnerabilities.

### Scope

Common areas of concern:

- **Credential handling** — never store secrets in `workflow/info.plist` or
  committed files; use Alfred's built-in encrypted keychain instead.
- **Input sanitization** — Alfred query strings are passed to
  `cmd/paste-formatted-date-alfred`; they must not be interpolated into shell
  commands or SQL without sanitization.
- **Dependency security** — this project has no third-party Go dependencies
  by design; dependabot monitors `.github/workflows/` automatically.

### Automated Security Checks

| Hook | What it detects |
|---|---|
| `gitleaks` (`.gitleaks.toml`) | Hardcoded secrets, API keys, local absolute paths |
| `detect-private-key` | SSH/TLS private key headers |
| `no-commit-dotenv` | `.env` files accidentally staged |
| `check-added-large-files` | Files over 500 KB |

These hooks run on every commit (pre-commit) and in CI (`security` job).
