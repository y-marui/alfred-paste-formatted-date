package datecmd

import (
	"strings"
	"testing"

	"github.com/y-marui/alfred-paste-formatted-date/internal/datefmt"
)

func TestDateEmptyQueryReturnsAllFormats(t *testing.T) {
	resp := Dispatch("")
	if len(resp.Items) != len(datefmt.All) {
		t.Fatalf("got %d items, want %d", len(resp.Items), len(datefmt.All))
	}
}

func TestDateAllItemsAreValid(t *testing.T) {
	resp := Dispatch("")
	for _, it := range resp.Items {
		if it.Valid == nil || !*it.Valid {
			t.Errorf("item %+v is not valid", it)
		}
	}
}

func TestDateFilterByFormatLabel(t *testing.T) {
	resp := Dispatch("ISO")
	if len(resp.Items) == 0 {
		t.Fatal("expected at least one item")
	}
	for _, it := range resp.Items {
		matches := strings.Contains(strings.ToLower(it.Subtitle), "iso") ||
			strings.Contains(strings.ToLower(it.Title), "iso") ||
			strings.Contains(strings.ToLower(it.UID), "iso")
		if !matches {
			t.Errorf("item %+v does not match filter %q", it, "iso")
		}
	}
}

func TestDateFilterByFormatUIDYYYYMMDD(t *testing.T) {
	resp := Dispatch("YYYYMMDD")
	if len(resp.Items) < 1 {
		t.Fatal("expected at least one item")
	}
	if resp.Items[0].Subtitle != "YYYYMMDD" {
		t.Errorf("Subtitle = %q, want %q", resp.Items[0].Subtitle, "YYYYMMDD")
	}
}

func TestDateNoMatchReturnsErrorItem(t *testing.T) {
	resp := Dispatch("xyzzy-nonexistent")
	if len(resp.Items) != 1 {
		t.Fatalf("got %d items, want 1", len(resp.Items))
	}
	if resp.Items[0].Valid == nil || *resp.Items[0].Valid {
		t.Errorf("item should be invalid, got %+v", resp.Items[0])
	}
}

func TestDateArgEqualsTitle(t *testing.T) {
	resp := Dispatch("")
	for _, it := range resp.Items {
		if it.Arg != it.Title {
			t.Errorf("Arg %q != Title %q", it.Arg, it.Title)
		}
	}
}

func TestDateUnixTimestampIsNumeric(t *testing.T) {
	resp := Dispatch("unix")
	if len(resp.Items) != 1 {
		t.Fatalf("got %d items, want 1", len(resp.Items))
	}
	for _, c := range resp.Items[0].Arg {
		if c < '0' || c > '9' {
			t.Errorf("Arg %q is not all digits", resp.Items[0].Arg)
			break
		}
	}
}

func TestDateRelativeOffsetShowsAllFormats(t *testing.T) {
	resp := Dispatch("-2d")
	if len(resp.Items) != len(datefmt.All) {
		t.Fatalf("got %d items, want %d", len(resp.Items), len(datefmt.All))
	}
}

func TestDateDirectDateISOValue(t *testing.T) {
	resp := Dispatch("2030/1/15")
	found := false
	for _, it := range resp.Items {
		if it.Subtitle == "YYYY-MM-DD" {
			found = true
			if it.Title != "2030-01-15" {
				t.Errorf("Title = %q, want %q", it.Title, "2030-01-15")
			}
		}
	}
	if !found {
		t.Fatal("no YYYY-MM-DD item found")
	}
}

func TestHelpShowsAllCommands(t *testing.T) {
	resp := Dispatch("help")
	if len(resp.Items) != len(helpCommands) {
		t.Fatalf("got %d items, want %d", len(resp.Items), len(helpCommands))
	}
}

func TestHelpAllItemsInvalid(t *testing.T) {
	resp := Dispatch("help")
	for _, it := range resp.Items {
		if it.Valid == nil || *it.Valid {
			t.Errorf("item %+v should be invalid", it)
		}
	}
}
