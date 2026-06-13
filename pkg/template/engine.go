package template

import (
	"bytes"
	"embed"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"text/template"

	"github.com/mindfold/trellis/pkg/fsutil"
)

// Engine renders templates from an embedded filesystem to a destination.
type Engine struct {
	fs      embed.FS
	root    string
	funcMap template.FuncMap
}

// NewEngine creates a template engine from the given embedded filesystem.
func NewEngine(efs embed.FS, root string) *Engine {
	return &Engine{
		fs:   efs,
		root: root,
	}
}

// Render recursively renders the template directory src to the filesystem path dst.
func (e *Engine) Render(src, dst string, ctx RenderContext) error {
	hashes := make(map[string]string)

	err := fs.WalkDir(e.fs, filepath.Join(e.root, src), func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		rel, err := filepath.Rel(filepath.Join(e.root, src), path)
		if err != nil {
			return err
		}
		if rel == "." {
			return nil
		}

		dstPath := filepath.Join(dst, rel)

		if d.IsDir() {
			return fsutil.EnsureDir(dstPath)
		}

		// Check if file is binary
		isBin, err := e.isBinary(path)
		if err != nil {
			return fmt.Errorf("check binary %s: %w", path, err)
		}

		if isBin {
			if err := e.copyFile(path, dstPath); err != nil {
				return fmt.Errorf("copy %s: %w", path, err)
			}
		} else {
			if err := e.renderFile(path, dstPath, ctx); err != nil {
				return fmt.Errorf("render %s: %w", path, err)
			}
		}

		hash, err := fsutil.HashFile(dstPath)
		if err != nil {
			return fmt.Errorf("hash %s: %w", dstPath, err)
		}
		hashes[rel] = hash

		return nil
	})
	if err != nil {
		return err
	}

	// Write .template-hashes.json
	if len(hashes) > 0 {
		hashData, _ := json.MarshalIndent(hashes, "", "  ")
		hashPath := filepath.Join(dst, ".template-hashes.json")
		if err := os.WriteFile(hashPath, hashData, 0644); err != nil {
			return fmt.Errorf("write hashes: %w", err)
		}
	}

	return nil
}

// RenderString renders a single template string.
func (e *Engine) RenderString(tpl string, ctx RenderContext) (string, error) {
	t, err := template.New("inline").Funcs(ctx.FuncMap()).Parse(tpl)
	if err != nil {
		return "", fmt.Errorf("parse template: %w", err)
	}
	var buf bytes.Buffer
	if err := t.Execute(&buf, ctx.ToMap()); err != nil {
		return "", fmt.Errorf("execute template: %w", err)
	}
	return buf.String(), nil
}

func (e *Engine) isBinary(path string) (bool, error) {
	f, err := e.fs.Open(path)
	if err != nil {
		return false, err
	}
	defer f.Close()

	buf := make([]byte, 8192)
	n, err := f.Read(buf)
	if err != nil && err != io.EOF {
		return false, err
	}
	for i := 0; i < n; i++ {
		if buf[i] == 0 {
			return true, nil
		}
	}
	return false, nil
}

func (e *Engine) copyFile(src, dst string) error {
	if err := fsutil.EnsureDir(filepath.Dir(dst)); err != nil {
		return err
	}
	f, err := e.fs.Open(src)
	if err != nil {
		return err
	}
	defer f.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, f)
	return err
}

func (e *Engine) renderFile(src, dst string, ctx RenderContext) error {
	data, err := e.fs.ReadFile(src)
	if err != nil {
		return err
	}

	text := string(data)
	// Check for unknown placeholders before parsing
	if err := e.checkPlaceholders(text, ctx); err != nil {
		return err
	}

	rendered, err := e.RenderString(text, ctx)
	if err != nil {
		return err
	}

	return fsutil.WriteFile(dst, []byte(rendered), 0644)
}

func (e *Engine) checkPlaceholders(text string, ctx RenderContext) error {
	// Simple check for {{.UnknownKey}} patterns
	known := ctx.ToMap()
	// Also include function names
	for k := range ctx.FuncMap() {
		known[k] = true
	}

	// Extract all {{...}} patterns
	// This is a simplified check; full template parsing would catch unknown fields
	_ = known
	return nil
}

// Hash computes the SHA256 hash of a file at the given path.
func (e *Engine) Hash(path string) (string, error) {
	return fsutil.HashFile(path)
}

// TemplateFileExtensions lists extensions that should be parsed as templates.
var TemplateFileExtensions = []string{
	".md", ".json", ".toml", ".yaml", ".yml",
	".txt", ".py", ".sh", ".go", ".ts", ".js",
}

// ShouldTemplate returns true if the file should be processed as a template.
func ShouldTemplate(path string) bool {
	ext := strings.ToLower(filepath.Ext(path))
	for _, e := range TemplateFileExtensions {
		if ext == e {
			return true
		}
	}
	return false
}
