package context

import "github.com/superops-team/trellis-go/pkg/manifest"

// Entry is a single reference in a context manifest.
type Entry = manifest.Entry

// Manifest corresponds to implement.jsonl / check.jsonl in memory.
type Manifest = manifest.Manifest

// LoadManifest reads a manifest from a JSONL file.
func LoadManifest(path string) (*Manifest, error) {
	return manifest.Load(path)
}
