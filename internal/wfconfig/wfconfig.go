// Package wfconfig is a small persistent key/value store for the "config"
// command, backed by a JSON file in Alfred's workflow data directory.
//
// Alfred exposes the data directory via the alfred_workflow_data
// environment variable. Values are stored as a flat JSON file, which Alfred
// can also read/write via its built-in "Set Variable" / "Universal Action"
// objects.
//
// Outside Alfred, DataDir falls back to
// ~/.config/alfred-workflow/<bundle-id>/.
package wfconfig

import (
	"encoding/json"
	"os"
	"path/filepath"
)

const fileName = "config.json"

// DataDir returns the workflow data directory, creating it if necessary.
func DataDir() string {
	dir := os.Getenv("alfred_workflow_data")
	if dir == "" {
		bundleID := os.Getenv("alfred_workflow_bundleid")
		if bundleID == "" {
			bundleID = "alfred-workflow-template"
		}
		home, err := os.UserHomeDir()
		if err != nil {
			home = "."
		}
		dir = filepath.Join(home, ".config", "alfred-workflow", bundleID)
	}
	_ = os.MkdirAll(dir, 0o755)
	return dir
}

// Store is a persistent key-value store for workflow configuration.
type Store struct {
	dir string
}

// New returns a Store backed by a config.json file inside dir.
func New(dir string) *Store {
	return &Store{dir: dir}
}

func (s *Store) path() string {
	return filepath.Join(s.dir, fileName)
}

func (s *Store) load() map[string]any {
	data, err := os.ReadFile(s.path())
	if err != nil {
		return map[string]any{}
	}
	var out map[string]any
	if err := json.Unmarshal(data, &out); err != nil {
		return map[string]any{}
	}
	return out
}

func (s *Store) save(data map[string]any) error {
	encoded, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(s.path(), encoded, 0o644)
}

// Get returns the stored value for key and whether it was present.
func (s *Store) Get(key string) (any, bool) {
	v, ok := s.load()[key]
	return v, ok
}

// Set stores value under key.
func (s *Store) Set(key string, value any) error {
	data := s.load()
	data[key] = value
	return s.save(data)
}

// Delete removes key, if present.
func (s *Store) Delete(key string) error {
	data := s.load()
	delete(data, key)
	return s.save(data)
}

// All returns every stored key/value pair.
func (s *Store) All() map[string]any {
	return s.load()
}

// Reset clears every stored key/value pair.
func (s *Store) Reset() error {
	return s.save(map[string]any{})
}
