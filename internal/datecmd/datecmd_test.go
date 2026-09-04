package datecmd

import (
	"strings"
	"testing"

	"github.com/y-marui/alfred-paste-formatted-date/internal/datefmt"
	"github.com/y-marui/alfred-paste-formatted-date/internal/wfconfig"
)

func setupConfigEnv(t *testing.T) {
	t.Helper()
	t.Setenv("alfred_workflow_data", t.TempDir())
	t.Setenv("alfred_workflow_bundleid", "com.example.test")
}

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

func TestConfigEmptyShowsNoSettings(t *testing.T) {
	setupConfigEnv(t)
	resp := Dispatch("config")
	found := false
	for _, it := range resp.Items {
		if strings.Contains(it.Title, "No settings") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected a 'No settings' item, got %+v", resp.Items)
	}
}

func TestConfigResetClearsConfig(t *testing.T) {
	setupConfigEnv(t)
	store := wfconfig.New(wfconfig.DataDir())
	if err := store.Set("key", "value"); err != nil {
		t.Fatal(err)
	}

	resp := Dispatch("config reset")
	if len(resp.Items) == 0 || !strings.Contains(strings.ToLower(resp.Items[0].Title), "reset") {
		t.Errorf("expected a reset confirmation item, got %+v", resp.Items)
	}
	if all := store.All(); len(all) != 0 {
		t.Errorf("config not cleared: %v", all)
	}
}

func TestConfigShowsExistingSettings(t *testing.T) {
	setupConfigEnv(t)
	store := wfconfig.New(wfconfig.DataDir())
	if err := store.Set("some_key", "some_value"); err != nil {
		t.Fatal(err)
	}

	resp := Dispatch("config")
	found := false
	for _, it := range resp.Items {
		if strings.Contains(it.Title, "some_key") {
			found = true
		}
	}
	if !found {
		t.Errorf("expected an item mentioning some_key, got %+v", resp.Items)
	}
}

func TestConfigUnknownSubcommandShowsCurrentConfig(t *testing.T) {
	setupConfigEnv(t)
	resp := Dispatch("config unknown-subcommand")
	if len(resp.Items) == 0 {
		t.Error("expected at least one item")
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
