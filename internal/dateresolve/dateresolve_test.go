package dateresolve

import (
	"testing"
	"time"
)

func sameDate(t time.Time, y int, m time.Month, d int) bool {
	yy, mm, dd := t.Date()
	return yy == y && mm == m && dd == d
}

func today() (int, time.Month, int) {
	return time.Now().Date()
}

func TestResolveEmptyReturnsToday(t *testing.T) {
	dt, remaining := Resolve("")
	y, m, d := today()
	if !sameDate(dt, y, m, d) {
		t.Errorf("got %v, want today", dt)
	}
	if remaining != "" {
		t.Errorf("remaining = %q, want empty", remaining)
	}
}

func TestResolveWhitespaceReturnsToday(t *testing.T) {
	dt, remaining := Resolve("   ")
	y, m, d := today()
	if !sameDate(dt, y, m, d) {
		t.Errorf("got %v, want today", dt)
	}
	if remaining != "" {
		t.Errorf("remaining = %q, want empty", remaining)
	}
}

func TestResolveRelativeDays(t *testing.T) {
	cases := []struct {
		name  string
		query string
		delta int
	}{
		{"minus_no_unit", "-2", -2},
		{"plus_no_unit", "+3", 3},
		{"minus_with_unit", "-2d", -2},
		{"plus_with_unit", "+1d", 1},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			dt, remaining := Resolve(c.query)
			want := time.Now().AddDate(0, 0, c.delta)
			wy, wm, wd := want.Date()
			if !sameDate(dt, wy, wm, wd) {
				t.Errorf("Resolve(%q) date = %v, want %04d-%02d-%02d", c.query, dt, wy, wm, wd)
			}
			if remaining != "" {
				t.Errorf("remaining = %q, want empty", remaining)
			}
		})
	}
}

func TestResolveRelativeWeeks(t *testing.T) {
	cases := []struct {
		query string
		weeks int
	}{
		{"+1w", 1},
		{"-2w", -2},
	}
	for _, c := range cases {
		dt, remaining := Resolve(c.query)
		want := time.Now().AddDate(0, 0, c.weeks*7)
		wy, wm, wd := want.Date()
		if !sameDate(dt, wy, wm, wd) {
			t.Errorf("Resolve(%q) date = %v, want %04d-%02d-%02d", c.query, dt, wy, wm, wd)
		}
		if remaining != "" {
			t.Errorf("remaining = %q, want empty", remaining)
		}
	}
}

func TestResolveRelativeMonths(t *testing.T) {
	dt, remaining := Resolve("-3m")
	if !dt.Before(time.Now()) {
		t.Errorf("-3m should be before now, got %v", dt)
	}
	if remaining != "" {
		t.Errorf("remaining = %q, want empty", remaining)
	}

	dt, remaining = Resolve("+1m")
	if !dt.After(time.Now()) {
		t.Errorf("+1m should be after now, got %v", dt)
	}
	if remaining != "" {
		t.Errorf("remaining = %q, want empty", remaining)
	}
}

func TestResolveRelativeYears(t *testing.T) {
	dt, remaining := Resolve("+1y")
	if !dt.After(time.Now()) {
		t.Errorf("+1y should be after now, got %v", dt)
	}
	if remaining != "" {
		t.Errorf("remaining = %q, want empty", remaining)
	}

	dt, remaining = Resolve("-1y")
	if !dt.Before(time.Now()) {
		t.Errorf("-1y should be before now, got %v", dt)
	}
	if remaining != "" {
		t.Errorf("remaining = %q, want empty", remaining)
	}
}

func TestResolveDirectDate(t *testing.T) {
	cases := []struct {
		name  string
		query string
	}{
		{"4digit_slash", "2026/7/9"},
		{"2digit_slash", "26/7/9"},
		{"4digit_hyphen", "2026-7-9"},
		{"zero_padded", "2026/07/09"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			dt, remaining := Resolve(c.query)
			if !sameDate(dt, 2026, time.July, 9) {
				t.Errorf("Resolve(%q) date = %v, want 2026-07-09", c.query, dt)
			}
			if remaining != "" {
				t.Errorf("remaining = %q, want empty", remaining)
			}
		})
	}
}

func TestResolveDirectDateMidnight(t *testing.T) {
	dt, _ := Resolve("2026/7/9")
	h, m, s := dt.Clock()
	if h != 0 || m != 0 || s != 0 {
		t.Errorf("time = %02d:%02d:%02d, want midnight", h, m, s)
	}
}

func TestResolveOverflowingRelativeOffsetPassthrough(t *testing.T) {
	q := "+999999999999999999999999999999d"
	dt, remaining := Resolve(q)
	y, m, d := today()
	if !sameDate(dt, y, m, d) {
		t.Errorf("got %v, want today (offset overflows int and must not silently resolve)", dt)
	}
	if remaining != q {
		t.Errorf("remaining = %q, want %q", remaining, q)
	}
}

func TestResolveInvalidDatePassthrough(t *testing.T) {
	dt, remaining := Resolve("2026/13/1")
	y, m, d := today()
	if !sameDate(dt, y, m, d) {
		t.Errorf("got %v, want today", dt)
	}
	if remaining != "2026/13/1" {
		t.Errorf("remaining = %q, want %q", remaining, "2026/13/1")
	}
}

func TestResolveFormatFilterPassthrough(t *testing.T) {
	dt, remaining := Resolve("ISO")
	y, m, d := today()
	if !sameDate(dt, y, m, d) {
		t.Errorf("got %v, want today", dt)
	}
	if remaining != "ISO" {
		t.Errorf("remaining = %q, want %q", remaining, "ISO")
	}
}

func TestResolveFormatFilterPreservesCase(t *testing.T) {
	_, remaining := Resolve("YYYY")
	if remaining != "YYYY" {
		t.Errorf("remaining = %q, want %q", remaining, "YYYY")
	}
}

func TestResolveCombinedRelativeWithFilter(t *testing.T) {
	dt, remaining := Resolve("-2d ISO")
	want := time.Now().AddDate(0, 0, -2)
	wy, wm, wd := want.Date()
	if !sameDate(dt, wy, wm, wd) {
		t.Errorf("date = %v, want %04d-%02d-%02d", dt, wy, wm, wd)
	}
	if remaining != "ISO" {
		t.Errorf("remaining = %q, want %q", remaining, "ISO")
	}
}

func TestResolveCombinedDirectDateWithFilter(t *testing.T) {
	dt, remaining := Resolve("2026/7/9 unix")
	if !sameDate(dt, 2026, time.July, 9) {
		t.Errorf("date = %v, want 2026-07-09", dt)
	}
	if remaining != "unix" {
		t.Errorf("remaining = %q, want %q", remaining, "unix")
	}
}

func TestResolveCombinedWeeksWithFilter(t *testing.T) {
	dt, remaining := Resolve("+1w YYYY")
	want := time.Now().AddDate(0, 0, 7)
	wy, wm, wd := want.Date()
	if !sameDate(dt, wy, wm, wd) {
		t.Errorf("date = %v, want %04d-%02d-%02d", dt, wy, wm, wd)
	}
	if remaining != "YYYY" {
		t.Errorf("remaining = %q, want %q", remaining, "YYYY")
	}
}

func TestResolveFilterOnlyNotConfusedWithDateSpec(t *testing.T) {
	dt, remaining := Resolve("YYYYMMDD")
	y, m, d := today()
	if !sameDate(dt, y, m, d) {
		t.Errorf("got %v, want today", dt)
	}
	if remaining != "YYYYMMDD" {
		t.Errorf("remaining = %q, want %q", remaining, "YYYYMMDD")
	}
}

func TestAddMonths(t *testing.T) {
	cases := []struct {
		base   time.Time
		months int
		want   time.Time
	}{
		{date(2026, 4, 14), 3, date(2026, 7, 14)},
		{date(2026, 4, 14), -3, date(2026, 1, 14)},
		{date(2026, 1, 31), 1, date(2026, 2, 28)}, // clamp to Feb 28
		{date(2024, 1, 31), 1, date(2024, 2, 29)}, // leap year
		{date(2026, 12, 31), 2, date(2027, 2, 28)},
	}
	for _, c := range cases {
		got := addMonths(c.base, c.months)
		if !got.Equal(c.want) {
			t.Errorf("addMonths(%v, %d) = %v, want %v", c.base, c.months, got, c.want)
		}
	}
}

func date(y int, m time.Month, d int) time.Time {
	return time.Date(y, m, d, 0, 0, 0, 0, time.UTC)
}
