package context

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"os"
	pathpkg "path"
	"path/filepath"
	"strings"

	"github.com/superops-team/trellis-go/pkg/spec"
)

const (
	injectMarker         = "<!-- trellis-hook-injected -->"
	maxContextEntryBytes = 256 * 1024
)

var ErrPRDRequired = errors.New("PRD is required")

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
	prd, err := LoadRequiredPRD(taskDir)
	if err != nil {
		return "", err
	}
	parts = append(parts, fmt.Sprintf("=== file: prd.md ===\n%s", prd))

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

	var skipped []string
	for _, entry := range manifest.Entries {
		content, err := b.loadEntry(entry)
		if err != nil {
			if entry.Required {
				return "", fmt.Errorf("required entry %s: %w", entry.Path, err)
			}
			skipped = append(skipped, fmt.Sprintf("- %s: %v", entry.Path, err))
			continue
		}
		parts = append(parts, fmt.Sprintf("=== file: %s ===\n%s", entry.Path, content))
	}
	appendSkippedOptional(&parts, skipped)

	return strings.Join(parts, "\n\n"), nil
}

// BuildCheckContext builds the check phase context.
func (b *Builder) BuildCheckContext(taskDir string) (string, error) {
	var parts []string
	parts = append(parts, injectMarker)

	// Load prd.md
	prd, err := LoadRequiredPRD(taskDir)
	if err != nil {
		return "", err
	}
	parts = append(parts, fmt.Sprintf("=== file: prd.md ===\n%s", prd))

	// Load check.jsonl entries
	manifestPath := filepath.Join(taskDir, "check.jsonl")
	manifest, err := LoadManifest(manifestPath)
	if err != nil {
		return "", fmt.Errorf("load check manifest: %w", err)
	}

	var skipped []string
	for _, entry := range manifest.Entries {
		content, err := b.loadEntry(entry)
		if err != nil {
			if entry.Required {
				return "", fmt.Errorf("required entry %s: %w", entry.Path, err)
			}
			skipped = append(skipped, fmt.Sprintf("- %s: %v", entry.Path, err))
			continue
		}
		parts = append(parts, fmt.Sprintf("=== file: %s ===\n%s", entry.Path, content))
	}
	appendSkippedOptional(&parts, skipped)

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

func (b *Builder) loadEntry(entry Entry) (string, error) {
	entryPath, err := NormalizeEntryPath(entry.Path)
	if err != nil {
		return "", err
	}
	path := filepath.Join(b.Root, filepath.FromSlash(entryPath))
	rootAbs, err := filepath.Abs(b.Root)
	if err != nil {
		return "", fmt.Errorf("resolve context root: %w", err)
	}
	pathAbs, err := filepath.Abs(path)
	if err != nil {
		return "", fmt.Errorf("resolve context entry %s: %w", entry.Path, err)
	}
	rel, err := filepath.Rel(rootAbs, pathAbs)
	if err != nil {
		return "", fmt.Errorf("validate context entry %s: %w", entry.Path, err)
	}
	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("invalid context path %q: escapes context root", entry.Path)
	}
	if isSensitiveContextPath(entryPath) {
		return "", fmt.Errorf("sensitive path")
	}
	f, err := os.Open(pathAbs)
	if err != nil {
		return "", err
	}
	defer f.Close()

	data, err := io.ReadAll(io.LimitReader(f, maxContextEntryBytes+1))
	if err != nil {
		return "", err
	}
	if len(data) > maxContextEntryBytes {
		return "", fmt.Errorf("too large: %d bytes exceeds %d", len(data), maxContextEntryBytes)
	}
	if isBinaryContent(data) {
		return "", fmt.Errorf("binary file")
	}
	return string(data), nil
}

func NormalizeEntryPath(rawPath string) (string, error) {
	normalized := strings.ReplaceAll(rawPath, "\\", "/")
	if filepath.IsAbs(rawPath) || strings.HasPrefix(normalized, "/") || strings.HasPrefix(normalized, "//") || hasWindowsVolume(normalized) {
		return "", fmt.Errorf("invalid context path %q: must be relative", rawPath)
	}
	cleaned := pathpkg.Clean(normalized)
	if cleaned == "." || cleaned == "" {
		return "", fmt.Errorf("invalid context path %q: path is required", rawPath)
	}
	for _, part := range strings.Split(cleaned, "/") {
		if part == ".." {
			return "", fmt.Errorf("invalid context path %q: cannot contain ..", rawPath)
		}
	}
	return cleaned, nil
}

func hasWindowsVolume(path string) bool {
	return len(path) >= 2 && path[1] == ':' && ((path[0] >= 'A' && path[0] <= 'Z') || (path[0] >= 'a' && path[0] <= 'z'))
}

func LoadRequiredPRD(taskDir string) (string, error) {
	data, err := os.ReadFile(filepath.Join(taskDir, "prd.md"))
	if err != nil {
		if os.IsNotExist(err) {
			return "", ErrPRDRequired
		}
		return "", fmt.Errorf("read PRD: %w", err)
	}
	if strings.TrimSpace(string(data)) == "" {
		return "", ErrPRDRequired
	}
	return string(data), nil
}

func isBinaryContent(data []byte) bool {
	limit := len(data)
	if limit > 8192 {
		limit = 8192
	}
	return bytes.IndexByte(data[:limit], 0) >= 0
}

func appendSkippedOptional(parts *[]string, skipped []string) {
	if len(skipped) == 0 {
		return
	}
	*parts = append(*parts, "=== skipped optional context ===\n"+strings.Join(skipped, "\n"))
}

func isSensitiveContextPath(entryPath string) bool {
	for _, segment := range strings.Split(entryPath, "/") {
		lower := strings.ToLower(segment)
		switch {
		case lower == ".env" || strings.HasPrefix(lower, ".env."):
			return true
		case lower == "id_rsa" || lower == "id_ed25519" || lower == "credentials.json":
			return true
		case lower == "secret" || lower == "secrets" || lower == "token" || lower == "tokens":
			return true
		case strings.HasSuffix(lower, ".pem") || strings.HasSuffix(lower, ".key"):
			return true
		}
	}
	return false
}
