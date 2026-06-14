package task

import "github.com/superops-team/trellis-go/pkg/manifest"

// ContextEntry is a single reference in a context manifest.
type ContextEntry = manifest.Entry

// Manifest corresponds to implement.jsonl / check.jsonl in memory.
type Manifest = manifest.Manifest

// loadManifest reads a manifest from a JSONL file.
func loadManifest(path string) (*Manifest, error) {
	return manifest.Load(path)
}
