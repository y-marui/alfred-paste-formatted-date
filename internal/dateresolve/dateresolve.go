// Package dateresolve parses a query string into a target date/time.
//
// # Supported query forms
//
// Relative offsets (sign + number + optional unit):
//
//	-2   -2d    two days ago
//	+1w         one week later
//	-3m         three months ago
//	+1y         one year later
//	Unit defaults to 'd' when omitted.
//
// Direct dates (year/month/day or year-month-day):
//
//	2026/7/9    26/7/9    2026-7-9
//	Two-digit years are interpreted as 2000+YY.
//
// Combined form (date spec + format filter, space-separated):
//
//	-2d ISO        two days ago, ISO formats only
//	2026/7/9 unix  specific date, unix timestamp only
//	+1w YYYY       one week later, formats containing "YYYY"
//
// Anything else is treated as a format filter; today is used as the target.
package dateresolve

import (
	"regexp"
	"strconv"
	"strings"
	"time"
)

// Resolve returns (targetTime, remainingFilterQuery) for the given query.
//
// The first whitespace-separated token is tested as a date spec. If it
// matches, the remainder of the query is returned as the filter. Otherwise
// the whole query is used as a format filter against today.
func Resolve(query string) (time.Time, string) {
	q := strings.TrimSpace(query)
	if q == "" {
		return time.Now(), ""
	}

	first, rest := splitFirst(q)

	if dt, ok := tryRelative(first); ok {
		return dt, strings.TrimSpace(rest)
	}

	if dt, ok := tryDirectDate(first); ok {
		return dt, strings.TrimSpace(rest)
	}

	return time.Now(), q
}

// splitFirst splits q at the first space character only (mirroring Python's
// str.partition(" ")), so any further internal whitespace in rest is left
// untouched until the caller trims it.
func splitFirst(q string) (first, rest string) {
	if idx := strings.IndexByte(q, ' '); idx >= 0 {
		return q[:idx], q[idx+1:]
	}
	return q, ""
}

var relativeRE = regexp.MustCompile(`(?i)^([+-])(\d+)([dwmy]?)$`)

func tryRelative(q string) (time.Time, bool) {
	m := relativeRE.FindStringSubmatch(q)
	if m == nil {
		return time.Time{}, false
	}

	sign := 1
	if m[1] == "-" {
		sign = -1
	}
	amount, _ := strconv.Atoi(m[2])
	amount *= sign
	unit := strings.ToLower(m[3])
	if unit == "" {
		unit = "d"
	}

	now := time.Now()
	today := truncateToDate(now)

	var target time.Time
	switch unit {
	case "d":
		target = today.AddDate(0, 0, amount)
	case "w":
		target = today.AddDate(0, 0, amount*7)
	case "m":
		target = addMonths(today, amount)
	default: // "y"
		target = addMonths(today, amount*12)
	}

	// Combine the resolved date with today's current time-of-day (seconds
	// precision), matching Python's
	// datetime.combine(target, datetime.now().time().replace(microsecond=0)).
	h, mnt, s := now.Clock()
	return time.Date(target.Year(), target.Month(), target.Day(), h, mnt, s, 0, now.Location()), true
}

var directDateRE = regexp.MustCompile(`^(\d{2,4})[/-](\d{1,2})[/-](\d{1,2})$`)

func tryDirectDate(q string) (time.Time, bool) {
	m := directDateRE.FindStringSubmatch(q)
	if m == nil {
		return time.Time{}, false
	}

	year, _ := strconv.Atoi(m[1])
	month, _ := strconv.Atoi(m[2])
	day, _ := strconv.Atoi(m[3])
	if year < 100 {
		year += 2000
	}

	dt := time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.Local)
	// time.Date normalizes out-of-range components (e.g. month 13) instead
	// of failing, so reject anything that didn't round-trip — mirroring
	// Python's datetime.date(...) raising ValueError for impossible dates.
	if dt.Year() != year || int(dt.Month()) != month || dt.Day() != day {
		return time.Time{}, false
	}
	return dt, true
}

func truncateToDate(t time.Time) time.Time {
	y, m, d := t.Date()
	return time.Date(y, m, d, 0, 0, 0, 0, t.Location())
}

// addMonths adds (or subtracts) a number of months, clamping the day to the
// last valid day of the resulting month. time.Time.AddDate can't be used
// here: it normalizes day overflow instead (e.g. Jan 31 + 1 month becomes
// Mar 3, not the Feb 28 clamp this function (and the ported Python
// _add_months) requires.
func addMonths(base time.Time, months int) time.Time {
	y, m, d := base.Date()
	totalMonths := int(m) - 1 + months

	yearOffset, monthIdx := floorDivMod(totalMonths, 12)
	year := y + yearOffset
	month := time.Month(monthIdx + 1)

	day := d
	if last := daysInMonth(year, month); day > last {
		day = last
	}
	return time.Date(year, month, day, 0, 0, 0, 0, base.Location())
}

// floorDivMod returns the floor-division quotient and non-negative
// remainder of a/b, matching Python's // and % operators (Go's native /
// and % truncate toward zero and can yield a negative remainder).
func floorDivMod(a, b int) (q, r int) {
	q = a / b
	r = a % b
	if r != 0 && (r < 0) != (b < 0) {
		q--
		r += b
	}
	return q, r
}

func daysInMonth(year int, month time.Month) int {
	// Day 0 of the following month is the last day of this one.
	return time.Date(year, month+1, 0, 0, 0, 0, 0, time.UTC).Day()
}
