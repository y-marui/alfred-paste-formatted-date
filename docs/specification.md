# Specification

> Functional specification, behavior definition, and data flow for
> alfred-paste-formatted-date.

## Overview

This workflow is an Alfred 5 Script Filter that accepts a keyword + query,
resolves a target date (today, a relative offset, or a direct date),
renders it in every registered format (optionally filtered), and returns a
JSON result list. Selecting a result copies/pastes its formatted value.

## Commands

### `date` (default)

**Trigger:** `date`, `date <offset>`, `date <YYYY/M/D>`, `date <filter>`

**Behavior:**
1. `internal/dateresolve.Resolve` reads the query's leading token:
   - Empty → today.
   - A relative offset matching `^[+-]\d+[dwmy]?$` (e.g. `-2`, `-2d`, `+1w`,
     `-3m`, `+1y`; no suffix defaults to days) → that many days/weeks/months/
     years from today.
   - A direct date matching `^\d{2,4}[/-]\d{1,2}[/-]\d{1,2}$` (e.g.
     `2026/7/9`, 2- or 4-digit year) → that calendar date.
   - Anything else → today, and the whole query is treated as a format
     filter instead.
2. The remaining text (if any) is used as a case-insensitive substring
   filter against each format's label, rendered value, or UID (e.g. `ISO`,
   `YYYY`, `unix`).
3. `internal/datefmt.All` is rendered for the resolved date via
   `datefmt.Value`, keeping only formats the filter matches (all of them,
   if the filter is empty).
4. If nothing matches → a single invalid item, `title: 'No format matches
   "<args>"'`, `subtitle` suggesting example queries.
5. Otherwise → one result item per matching format.

**Registered formats** (`internal/datefmt.All`, in display order): YYYYMMDD,
YYMMDD, YYYY-MM-DD, YYYY/MM/DD, MM/DD/YYYY, DD/MM/YYYY, MMM DD YYYY (abbreviated
month), MMMM DD YYYY (full month), YYYY-MM-DDThh:mm:ss, and a Unix timestamp.

**Result item fields:**

| Field | Source | Notes |
|---|---|---|
| `title` | the rendered date value | Primary display text |
| `subtitle` | the format's label (e.g. `"YYYY-MM-DD"`) | Secondary display text |
| `arg` | the rendered date value | Passed to Alfred's Copy to Clipboard node on Enter |
| `uid` | the format's UID (e.g. `"iso-date"`) | Used by Alfred for learned ordering |

### `date help`

**Trigger:** `date help`

**Behavior:** Display the two registered commands (`date`, `date help`)
with descriptions and autocomplete strings (`valid: false` for all items).

## Data Flow

```
Alfred input (keyword "date" + query string)
  │
  ▼
cmd/paste-formatted-date-alfred/main.go   reads os.Args[1]
  │
  ▼
dispatch(query)                            recover()-wraps the call below → error item on panic
  │
  ▼
internal/datecmd.Dispatch(query)           "help" vs. everything else
  │
  ├─ internal/dateresolve.Resolve(args)     target date + remaining filter text
  │
  └─ internal/datefmt.All / .Value()        renders the target date in each matching format
  │
  ▼
scriptfilter.Response.Write()              prints JSON to stdout → Alfred renders result list
```

## Error Handling

- Any panic during `internal/datecmd.Dispatch` is recovered by `dispatch` in
  `main.go`, which emits a single error result item containing the panic
  value.
- An unresolvable relative offset or direct date silently falls back to
  today rather than erroring — the user sees today's date in every format,
  which reads as "nothing matched" rather than a crash.

## Configuration Variables

None currently — `workflow/info.plist`'s `userconfigurationconfig` is an
empty array. See `docs/configuration-builder.md` for the mechanism to use
when a future command needs one.

## Constraints

- Script Filter response time target: **< 100 ms** (a compiled binary,
  usually well within budget).
- All output must go through `scriptfilter.Response.Write()` — never
  `fmt.Print*` directly.
- `cmd/paste-formatted-date-alfred/main.go` contains no business logic; it
  only dispatches argv and writes the response.
