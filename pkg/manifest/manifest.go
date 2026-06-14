package manifest

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
)

// Entry is a single reference in a context manifest.
type Entry struct {
	Path        string `json:"path"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required"`
}

// Manifest corresponds to implement.jsonl / check.jsonl in memory.
type Manifest struct {
	Version string  `json:"version"`
	Entries []Entry `json:"entries"`
}

// Load reads a manifest from a JSONL file.
func Load(path string) (*Manifest, error) {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return &Manifest{Version: "1.0", Entries: []Entry{}}, nil
		}
		return nil, fmt.Errorf("open manifest %s: %w", path, err)
	}
	defer f.Close()

	manifest := &Manifest{Version: "1.0", Entries: []Entry{}}
	scanner := bufio.NewScanner(f)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := scanner.Text()
		if line == "" {
			continue
		}
		var entry Entry
		if err := json.Unmarshal([]byte(line), &entry); err != nil {
			return nil, fmt.Errorf("parse manifest %s line %d: %w", path, lineNumber, err)
		}
		manifest.Entries = append(manifest.Entries, entry)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read manifest %s: %w", path, err)
	}
	return manifest, nil
}

// Save writes a manifest to a JSONL file.
func Save(path string, manifest *Manifest) error {
	f, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("create manifest %s: %w", path, err)
	}
	defer f.Close()

	for _, entry := range manifest.Entries {
		data, err := json.Marshal(entry)
		if err != nil {
			return fmt.Errorf("marshal entry: %w", err)
		}
		if _, err := f.WriteString(string(data) + "\n"); err != nil {
			return fmt.Errorf("write manifest %s: %w", path, err)
		}
	}
	return nil
}

// Save writes the manifest to a JSONL file.
func (m *Manifest) Save(path string) error {
	return Save(path, m)
}
