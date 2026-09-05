# UI Design

Alfred Script Filter workflows present results as a list of items in the Alfred
launcher. This document defines the UI conventions for result items in this
workflow.

## Result Item Structure

Alfred result items are JSON objects with the following fields used in this workflow:

| Field | Type | Required | Description |
|---|---|---|---|
| `title` | string | yes | Primary text (large, always visible) — the rendered date value |
| `subtitle` | string | no | Secondary text (small, below title) — the format's label |
| `arg` | string | no | The rendered date value; copied/pasted by Alfred's native output node on Enter |
| `uid` | string | no | The format's UID; used by Alfred for learned ordering |
| `valid` | bool | yes | If false, Enter does not trigger an action |
| `autocomplete` | string | no | Text inserted into Alfred's input on Tab — used by the `date help` rows |

## Text Guidelines

### No Unicode Emoji in `title` / `subtitle`

- **Prohibited:** `🔍 Search`, `✅ Done`, `📄 Document`
- **Allowed:** ASCII symbols — `>`, `*`, `[x]`, `(!)`, `--`
- **Reason:** Emoji rendering is inconsistent across Alfred versions and macOS
  updates. ASCII symbols are universally stable.

### Empty / Error States

- No format matches the filter → `title: 'No format matches "<args>"'`,
  `subtitle` suggesting example queries, `valid: false`.
- Error → panic recovery automatically shows a `"Workflow Error"` item; do
  not hide errors silently.

## Icon

- Workflow icon: `workflow/icon.png` (PNG, any size — Alfred scales it).
- Alfred controls light/dark mode; do not ship separate light/dark icons.
- No per-item icons are used in this workflow.

## Keyboard Shortcuts

These are standard Alfred behaviors — do not override them in the workflow:

| Key | Action |
|---|---|
| ↩ Enter | Copy/paste `arg` (native output node) |
| ⇥ Tab | Insert `autocomplete` text into Alfred's input (`date help` rows) |
| ⌘C | Copy `arg` to clipboard |
| ⌘L | Show `title` in Large Type |

## Layout Conventions by Command

### `date` results

```
title:    <rendered date value, e.g. "2026-09-05">
subtitle: <format label, e.g. "YYYY-MM-DD">
arg:      <rendered date value>
uid:      <format UID, e.g. "iso-date">
valid:    true
```

### `date help` items

```
title:    date <command>
subtitle: <command description>
valid:    false
autocomplete: <command trigger string>
```

### No-match / error rows

```
title:    'No format matches "<args>"' | "Workflow Error"
subtitle: <suggested example queries, or the error message>
valid:    false
```
