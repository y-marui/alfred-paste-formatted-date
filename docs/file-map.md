# File Map

> File-level dependency map for alfred-paste-formatted-date.
> Add to this as you explore the codebase during development.

## Entry Points

| File | Role |
|---|---|
| `cmd/paste-formatted-date-alfred/main.go` | Alfred executes this binary — the sole entry point |

## Call Flow

```
cmd/paste-formatted-date-alfred/main.go
  └─ dispatch(query)                          ← recover()-wrapped
       └─ internal/datecmd.Dispatch(query)
            ├─ internal/dateresolve.Resolve()  ← parses the target date / relative offset
            └─ internal/datefmt.All / .Value() ← renders that date in every registered format
       └─ scriptfilter.Response.Write(os.Stdout)
```

## Module Dependency Table

### `internal/`

| File | Imports from | Notes |
|---|---|---|
| `dateresolve/dateresolve.go` | stdlib only (`time`, `regexp`) | Resolves a query's leading token into a target date: today, a relative offset (`-2d`, `+1w`, `-3m`, `+1y`), or a direct date (`2026/7/9`); returns the remaining text as a format filter |
| `datefmt/datefmt.go` | stdlib only (`time`) | The registered `Format` list (label + `time.Time` → string) and `Value()` to render one |
| `datecmd/datecmd.go` | `internal/dateresolve`, `internal/datefmt`, `internal/scriptfilter` | Dispatches the query (`help` vs. everything else) and builds the Script Filter response |
| `scriptfilter/scriptfilter.go` | stdlib only | Script Filter JSON types + writer |

### `cmd/`

| File | Imports from | Notes |
|---|---|---|
| `paste-formatted-date-alfred/main.go` | `internal/datecmd`, `internal/scriptfilter` | Alfred boundary; argv dispatch + `recover()` only |

## Key Files for Customization

| File | What to change |
|---|---|
| `workflow/info.plist` | `bundleid`, keyword, UIDs, category, description |
| `internal/datefmt/datefmt.go` | The registered date formats (`All`) |
| `internal/dateresolve/dateresolve.go` | Relative/direct date parsing rules |
| `go.mod` | Module path |
| `workflow/icon.png` | Workflow icon |
