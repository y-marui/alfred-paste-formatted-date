// Package datefmt holds the list of date formats offered by the "date"
// command and formats a time.Time into each of them.
package datefmt

import (
	"strconv"
	"time"
)

// Format is one selectable date representation.
type Format struct {
	UID    string
	Label  string
	Layout string // Go reference-time layout; empty for the "unix" special case.
}

// All is the ordered list of formats shown for an unfiltered "date" query.
var All = []Format{
	{"yyyymmdd", "YYYYMMDD", "20060102"},
	{"yymmdd", "YYMMDD", "060102"},
	{"iso-date", "YYYY-MM-DD", "2006-01-02"},
	{"slash-ymd", "YYYY/MM/DD", "2006/01/02"},
	{"slash-mdy", "MM/DD/YYYY", "01/02/2006"},
	{"slash-dmy", "DD/MM/YYYY", "02/01/2006"},
	{"abbr-month", "MMM DD, YYYY", "Jan 02, 2006"},
	{"full-month", "MMMM DD, YYYY", "January 02, 2006"},
	{"iso-datetime", "YYYY-MM-DDThh:mm:ss", "2006-01-02T15:04:05"},
	{"unix", "Unix timestamp", ""},
}

// Value renders dt using f's layout, or as a Unix timestamp for the "unix" format.
func Value(f Format, dt time.Time) string {
	if f.UID == "unix" {
		return strconv.FormatInt(dt.Unix(), 10)
	}
	return dt.Format(f.Layout)
}
