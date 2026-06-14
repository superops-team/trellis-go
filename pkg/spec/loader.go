package spec

import (
	"fmt"
	"os"
	"path/filepath"
)

// Loader manages the spec directory.
type Loader struct {
	Root string
}

// NewLoader creates a spec loader.
func NewLoader(root string) *Loader {
	return &Loader{Root: root}
}

// Index builds an index of all accessible spec files.
func (l *Loader) Index() (*Index, error) {
	idx := &Index{
		Packages: make(map[string]PackageIndex),
		Guides:   []string{},
	}

	entries, err := os.ReadDir(l.Root)
	if err != nil {
		return nil, fmt.Errorf("read spec root: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			if entry.Name() == "guides" {
				guidesDir := filepath.Join(l.Root, entry.Name())
				guides, _ := os.ReadDir(guidesDir)
				for _, g := range guides {
					if !g.IsDir() {
						idx.Guides = append(idx.Guides, filepath.Join("guides", g.Name()))
					}
				}
			}
			continue
		}

		pkgName := entry.Name()
		pkgDir := filepath.Join(l.Root, pkgName)
		pkgIdx := PackageIndex{Layers: make(map[string]string)}

		layers, err := os.ReadDir(pkgDir)
		if err != nil {
			return nil, fmt.Errorf("read package %s: %w", pkgName, err)
		}
		for _, layer := range layers {
			if !layer.IsDir() {
				continue
			}
			indexPath := filepath.Join(pkgName, layer.Name(), "index.md")
			fullPath := filepath.Join(l.Root, indexPath)
			if _, err := os.Stat(fullPath); err == nil {
				pkgIdx.Layers[layer.Name()] = indexPath
			}
		}
		idx.Packages[pkgName] = pkgIdx
	}

	return idx, nil
}

// Load reads a single spec file.
func (l *Loader) Load(path string) (string, error) {
	fullPath := filepath.Join(l.Root, path)
	data, err := os.ReadFile(fullPath)
	if err != nil {
		return "", fmt.Errorf("read spec %s: %w", path, err)
	}
	return string(data), nil
}

// LoadPackage loads all specs for a given package.
func (l *Loader) LoadPackage(pkg string) (map[string]string, error) {
	result := make(map[string]string)
	pkgDir := filepath.Join(l.Root, pkg)
	entries, err := os.ReadDir(pkgDir)
	if err != nil {
		return nil, fmt.Errorf("read package %s: %w", pkg, err)
	}
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		indexPath := filepath.Join(pkgDir, entry.Name(), "index.md")
		data, err := os.ReadFile(indexPath)
		if err != nil {
			return nil, fmt.Errorf("read package %s layer %s: %w", pkg, entry.Name(), err)
		}
		result[entry.Name()] = string(data)
	}
	return result, nil
}

// LoadLayer loads a specific package/layer spec.
func (l *Loader) LoadLayer(pkg, layer string) (string, error) {
	path := filepath.Join(pkg, layer, "index.md")
	return l.Load(path)
}
