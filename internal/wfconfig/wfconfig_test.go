package wfconfig

import "testing"

func TestGetMissingReturnsFalse(t *testing.T) {
	s := New(t.TempDir())
	if v, ok := s.Get("missing"); ok || v != nil {
		t.Errorf("Get(missing) = (%v, %v), want (nil, false)", v, ok)
	}
}

func TestSetAndGet(t *testing.T) {
	s := New(t.TempDir())
	if err := s.Set("key", "value"); err != nil {
		t.Fatal(err)
	}
	v, ok := s.Get("key")
	if !ok || v != "value" {
		t.Errorf("Get(key) = (%v, %v), want (value, true)", v, ok)
	}
}

func TestDelete(t *testing.T) {
	s := New(t.TempDir())
	_ = s.Set("key", "value")
	if err := s.Delete("key"); err != nil {
		t.Fatal(err)
	}
	if _, ok := s.Get("key"); ok {
		t.Errorf("Get(key) after Delete: ok = true, want false")
	}
}

func TestReset(t *testing.T) {
	s := New(t.TempDir())
	_ = s.Set("a", float64(1))
	_ = s.Set("b", float64(2))
	if err := s.Reset(); err != nil {
		t.Fatal(err)
	}
	if all := s.All(); len(all) != 0 {
		t.Errorf("All() after Reset = %v, want empty", all)
	}
}

func TestAll(t *testing.T) {
	s := New(t.TempDir())
	_ = s.Set("x", float64(10))
	_ = s.Set("y", float64(20))
	all := s.All()
	if len(all) != 2 || all["x"] != float64(10) || all["y"] != float64(20) {
		t.Errorf("All() = %v, want map[x:10 y:20]", all)
	}
}
