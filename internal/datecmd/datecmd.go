// Package datecmd dispatches an Alfred "date" query to the matching
// command handler and builds the Script Filter response.
//
// Commands:
//
//	date                  — today in all formats (default)
//	date -2  / -2d        — 2 days ago
//	date +1w              — 1 week later
//	date -3m              — 3 months ago
//	date +1y              — 1 year later
//	date 2026/7/9         — specific date (4- or 2-digit year, / or - separator)
//	date <filter>         — filter by format name or value (e.g. "ISO", "YYYY", "unix")
//	date config           — view current configuration
//	date config reset     — clear all stored configuration
//	date help             — show available commands
package datecmd

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"unicode"

	"github.com/y-marui/alfred-paste-formatted-date/internal/datefmt"
	"github.com/y-marui/alfred-paste-formatted-date/internal/dateresolve"
	"github.com/y-marui/alfred-paste-formatted-date/internal/scriptfilter"
	"github.com/y-marui/alfred-paste-formatted-date/internal/wfconfig"
)

// Dispatch parses the raw Alfred query and routes it to the matching
// command, falling back to the date command when the leading token isn't a
// registered command name (e.g. "ISO" or "-2d").
func Dispatch(query string) scriptfilter.Response {
	trimmed := strings.TrimSpace(query)
	command, rest := splitCommand(trimmed)

	switch strings.ToLower(command) {
	case "config":
		return handleConfig(rest)
	case "help":
		return handleHelp()
	case "date":
		return handleDate(rest)
	default:
		return handleDate(trimmed)
	}
}

// splitCommand splits s into its first whitespace-delimited token and the
// remainder, with the remainder's leading whitespace stripped — mirroring
// Python's "s.strip().split(None, 1)".
func splitCommand(s string) (command, rest string) {
	i := strings.IndexFunc(s, unicode.IsSpace)
	if i < 0 {
		return s, ""
	}
	return s[:i], strings.TrimLeftFunc(s[i:], unicode.IsSpace)
}

// handleDate resolves the target date from args, then shows matching date formats.
func handleDate(args string) scriptfilter.Response {
	targetDT, filterQuery := dateresolve.Resolve(args)
	query := strings.ToLower(strings.TrimSpace(filterQuery))

	var items []scriptfilter.Item
	for _, f := range datefmt.All {
		value := datefmt.Value(f, targetDT)
		if query != "" &&
			!strings.Contains(strings.ToLower(f.Label), query) &&
			!strings.Contains(strings.ToLower(value), query) &&
			!strings.Contains(f.UID, query) {
			continue
		}
		items = append(items, scriptfilter.Item{
			Title:    value,
			Subtitle: f.Label,
			Arg:      value,
			UID:      f.UID,
			Valid:    scriptfilter.BoolPtr(true),
		})
	}

	if len(items) == 0 {
		return scriptfilter.Response{
			Items: []scriptfilter.Item{{
				Title:    fmt.Sprintf(`No format matches "%s"`, args),
				Subtitle: "Try: YYYY, MM, DD, unix, ISO, -2d, +1w, 2026/7/9 ...",
				Valid:    scriptfilter.BoolPtr(false),
			}},
		}
	}
	return scriptfilter.Response{Items: items}
}

// handleConfig shows config items or performs a config action ("reset").
func handleConfig(args string) scriptfilter.Response {
	store := wfconfig.New(wfconfig.DataDir())
	sub := strings.ToLower(strings.TrimSpace(args))

	if sub == "reset" {
		if err := store.Reset(); err != nil {
			panic(err)
		}
		return scriptfilter.Response{
			Items: []scriptfilter.Item{{
				Title:    "Configuration reset",
				Subtitle: "All settings have been cleared",
				Valid:    scriptfilter.BoolPtr(false),
			}},
		}
	}

	resetItem := scriptfilter.Item{
		Title:        "Reset all settings",
		Subtitle:     "date config reset  — clear all stored configuration",
		Arg:          "reset",
		UID:          "config-reset",
		Autocomplete: "config reset",
		Valid:        scriptfilter.BoolPtr(true),
	}

	current := store.All()
	if len(current) == 0 {
		return scriptfilter.Response{
			Items: []scriptfilter.Item{
				{
					Title:    "No settings configured",
					Subtitle: "Settings will appear here once set",
					Valid:    scriptfilter.BoolPtr(false),
				},
				resetItem,
			},
		}
	}

	keys := make([]string, 0, len(current))
	for k := range current {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	items := make([]scriptfilter.Item, 0, len(keys)+1)
	for _, k := range keys {
		v := current[k]
		items = append(items, scriptfilter.Item{
			Title:    fmt.Sprintf("%s: %v", k, v),
			Subtitle: "Current setting",
			Arg:      fmt.Sprintf("%v", v),
			UID:      "config-" + k,
			Valid:    scriptfilter.BoolPtr(false),
		})
	}
	items = append(items, resetItem)

	return scriptfilter.Response{Items: items}
}

type helpEntry struct {
	cmd, desc, autocomplete string
}

var helpCommands = []helpEntry{
	{"date", "List all date formats (default command)", "date "},
	{"date config", "View or reset configuration", "date config"},
	{"date help", "Show this help", "date help"},
}

// handleHelp displays all available commands.
func handleHelp() scriptfilter.Response {
	items := make([]scriptfilter.Item, len(helpCommands))
	for i, c := range helpCommands {
		items[i] = scriptfilter.Item{
			Title:        c.cmd,
			Subtitle:     c.desc,
			UID:          "help-" + strconv.Itoa(i),
			Valid:        scriptfilter.BoolPtr(false),
			Autocomplete: c.autocomplete,
		}
	}
	return scriptfilter.Response{Items: items}
}
