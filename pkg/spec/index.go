package spec

import (
	"fmt"
	"strings"
)

// Index holds the spec directory structure.
type Index struct {
	Packages map[string]PackageIndex `json:"packages"`
	Guides   []string                `json:"guides"`
}

// PackageIndex indexes layers within a package.
type PackageIndex struct {
	Layers map[string]string `json:"layers"` // layer -> index.md path
}

// ToMarkdown renders the index as an AI-readable spec directory.
func (idx *Index) ToMarkdown() string {
	var b strings.Builder
	b.WriteString("# Spec Index\n\n")

	if len(idx.Packages) > 0 {
		b.WriteString("## Packages\n\n")
		for pkgName, pkg := range idx.Packages {
			b.WriteString(fmt.Sprintf("### %s\n\n", pkgName))
			for layer, path := range pkg.Layers {
				b.WriteString(fmt.Sprintf("- **%s**: `%s`\n", layer, path))
			}
			b.WriteString("\n")
		}
	}

	if len(idx.Guides) > 0 {
		b.WriteString("## Guides\n\n")
		for _, g := range idx.Guides {
			b.WriteString(fmt.Sprintf("- `%s`\n", g))
		}
		b.WriteString("\n")
	}

	return b.String()
}
