package manifest

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"

	"github.com/superops-team/trellis-go/pkg/fsutil"
)

const maxManifestLineBytes = 1024 * 1024

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
	scanner.Buffer(make([]byte, 0, 64*1024), maxManifestLineBytes)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		line := scanner.Bytes()
		if len(line) == 0 {
			continue
		}
		var entry Entry
		if err := json.Unmarshal(line, &entry); err != nil {
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
	var data []byte
	for _, entry := range manifest.Entries {
		line, err := json.Marshal(entry)
		if err != nil {
			return fmt.Errorf("marshal entry: %w", err)
		}
		data = append(data, line...)
		data = append(data, '\n')
	}
	if err := fsutil.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("write manifest %s: %w", path, err)
	}
	return nil
}

// Save writes the manifest to a JSONL file.
func (m *Manifest) Save(path string) error {
	return Save(path, m)
}
