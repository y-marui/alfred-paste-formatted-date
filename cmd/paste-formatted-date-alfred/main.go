// Command paste-formatted-date-alfred is the binary the packaged Alfred
// Workflow invokes (see workflow/info.plist). Alfred's Script Filter node
// runs it with the query following the "date" keyword as $1, e.g. "-2d ISO".
package main

import (
	"fmt"
	"os"

	"github.com/y-marui/alfred-paste-formatted-date/internal/datecmd"
	"github.com/y-marui/alfred-paste-formatted-date/internal/scriptfilter"
)

func main() {
	query := ""
	if len(os.Args) > 1 {
		query = os.Args[1]
	}
	writeResponse(dispatch(query))
}

// dispatch recovers from any panic in datecmd, mirroring the Python
// workflow's safe_run: an unhandled failure must still produce a visible
// Script Filter error item rather than empty/invalid output.
func dispatch(query string) (resp scriptfilter.Response) {
	defer func() {
		if r := recover(); r != nil {
			resp = errorResponse(fmt.Sprintf("%v", r))
		}
	}()
	return datecmd.Dispatch(query)
}

func writeResponse(resp scriptfilter.Response) {
	if err := resp.Write(os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "paste-formatted-date-alfred: writing response:", err)
		os.Exit(1)
	}
}

func errorResponse(message string) scriptfilter.Response {
	return scriptfilter.Response{
		Items: []scriptfilter.Item{
			{
				Title:    "Workflow Error",
				Subtitle: message,
				Arg:      message,
				Valid:    scriptfilter.BoolPtr(false),
			},
		},
	}
}
