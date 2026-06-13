package context

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/mindfold/trellis/pkg/spec"
)

const injectMarker = "<!-- trellis-hook-injected -->"

// TaskInfo holds minimal task information for context building.
type TaskInfo struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// Builder assembles sub-agent context from task manifests.
type Builder struct {
	SpecLoader *spec.Loader
	Root       string
}

// BuildImplementContext builds the implementation phase context.
func (b *Builder) BuildImplementContext(taskDir string) (string, error) {
	var parts []string
	parts = append(parts, injectMarker)

	// Load prd.md
	prdPath := filepath.Join(taskDir, "prd.md")
	if data, err := os.ReadFile(prdPath); err == nil && len(data) > 0 {
		parts = append(parts, fmt.Sprintf("=== file: prd.md ===\n%s", data))
	}

	// Load info.md if exists
	infoPath := filepath.Join(taskDir, "info.md")
	if data, err := os.ReadFile(infoPath); err == nil && len(data) > 0 {
		parts = append(parts, fmt.Sprintf("=== file: info.md ===\n%s", data))
	}

	// Load implement.jsonl entries
	manifestPath := filepath.Join(taskDir, "implement.jsonl")
	manifest, err := LoadManifest(manifestPath)
	if err != nil {
		return "", fmt.Errorf("load implement manifest: %w", err)
	}

	for _, entry := range manifest.Entries {
		content, err := b.loadEntry(taskDir, entry)
		if err != nil {
			if entry.Required {
				return "", fmt.Errorf("required entry %s: %w", entry.Path, err)
			}
			continue
		}
		parts = append(parts, fmt.Sprintf("=== file: %s ===\n%s", entry.Path, content))
	}

	return strings.Join(parts, "\n\n"), nil
}

// BuildCheckContext builds the check phase context.
func (b *Builder) BuildCheckContext(taskDir string) (string, error) {
	var parts []string
	parts = append(parts, injectMarker)

	// Load prd.md
	prdPath := filepath.Join(taskDir, "prd.md")
	if data, err := os.ReadFile(prdPath); err == nil && len(data) > 0 {
		parts = append(parts, fmt.Sprintf("=== file: prd.md ===\n%s", data))
	}

	// Load check.jsonl entries
	manifestPath := filepath.Join(taskDir, "check.jsonl")
	manifest, err := LoadManifest(manifestPath)
	if err != nil {
		return "", fmt.Errorf("load check manifest: %w", err)
	}

	for _, entry := range manifest.Entries {
		content, err := b.loadEntry(taskDir, entry)
		if err != nil {
			if entry.Required {
				return "", fmt.Errorf("required entry %s: %w", entry.Path, err)
			}
			continue
		}
		parts = append(parts, fmt.Sprintf("=== file: %s ===\n%s", entry.Path, content))
	}

	return strings.Join(parts, "\n\n"), nil
}

// BuildResearchContext builds the research phase context.
func (b *Builder) BuildResearchContext() (string, error) {
	var parts []string
	parts = append(parts, injectMarker)

	if b.SpecLoader != nil {
		idx, err := b.SpecLoader.Index()
		if err == nil {
			parts = append(parts, idx.ToMarkdown())
		}
	}

	return strings.Join(parts, "\n\n"), nil
}

func (b *Builder) loadEntry(taskDir string, entry Entry) (string, error) {
	path := filepath.Join(taskDir, entry.Path)
	if !filepath.IsAbs(entry.Path) {
		path = filepath.Join(b.Root, entry.Path)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
