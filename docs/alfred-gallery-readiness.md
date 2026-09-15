# Alfred Gallery Readiness

Tracks this Workflow's compliance with the
[Alfred Gallery submission requirements](https://alfred.app/submit/) and
[style guide](https://alfred.app/submit/styleguide/), per dev-charter's
[`topics/alfred/ALFRED_GALLERY.md`](docs/dev-charter/topics/alfred/ALFRED_GALLERY.md).
This is a checklist against an external, occasionally-changing policy —
re-read the linked pages before acting on stale entries here.

## Submission process

Alfred does not accept Gallery submissions directly. The documented path is:
share the workflow on the [Alfred Forum](https://www.alfredforum.com/) first;
once it is "generally stable and trusted by a number of users," the Alfred
team may invite an official Gallery submission. There is no self-service
form. Forum posting itself stays out of scope here unless explicitly
requested — this document only tracks technical/documentation readiness so
that submission is not blocked on our side whenever that step happens.

## Checklist

| Requirement | Status | Notes |
|---|---|---|
| Binaries signed and notarised | ✅ Done | Binary verified signed/notarised in the actual [v0.1.0 release](https://github.com/y-marui/alfred-paste-formatted-date/releases/tag/v0.1.0) (`codesign -dvvv`, `spctl -a -vvv -t install`, `gh attestation verify`) |
| No self-update | ✅ Done | Updates ship only as new `.alfredworkflow` releases; no self-update code path |
| No self-installed external software | ✅ Done | `go.mod` has no third-party dependencies; nothing is fetched at runtime |
| Icon ≥ 256×256 px | ✅ Done | `workflow/icon.png` is 376×376 |
| Keyword ≥ 3 characters | ✅ Done | `date` (4 characters) |
| Category set | ✅ Done | `Productivity` — see [`docs/alfred-workflow-notes/workflow-object-schema.md`](docs/alfred-workflow-notes/workflow-object-schema.md#category) for the valid value list |
| User Configuration over environment variables | ✅ N/A | No Config Builder variables are declared today |
| English instructions in About/README | ✅ Done | `README.md` is the reference (English) version; `README-jp.md` is canonical |
| README follows Gallery style guide | ✅ Done | `## Usage` opens with "via the `date` keyword" phrasing |
| Screenshots (full Alfred window, shadow, no background) | ✅ Done | `images/date.png`, real window capture (rounded corners + drop shadow, transparent background) |

## Out of scope here

- Posting to the Alfred Forum and the Gallery submission itself — a
  per-project decision, not mandated by this checklist
- Exporting a Developer ID certificate, generating a notarization API key,
  and registering the GitHub Actions secrets `release.yml` expects — manual
  steps only the repository owner can perform
