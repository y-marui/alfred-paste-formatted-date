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
//	date help             — show available commands
package datecmd

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"

	"github.com/y-marui/alfred-paste-formatted-date/internal/datefmt"
	"github.com/y-marui/alfred-paste-formatted-date/internal/dateresolve"
	"github.com/y-marui/alfred-paste-formatted-date/internal/scriptfilter"
)

// Dispatch parses the raw Alfred query and routes it to the matching
// command, falling back to the date command when the leading token isn't a
// registered command name (e.g. "ISO" or "-2d").
func Dispatch(query string) scriptfilter.Response {
	trimmed := strings.TrimSpace(query)
	command, rest := splitCommand(trimmed)

	switch strings.ToLower(command) {
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

type helpEntry struct {
	cmd, desc, autocomplete string
}

var helpCommands = []helpEntry{
	{"date", "List all date formats (default command)", "date "},
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
