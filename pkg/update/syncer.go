package update

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/superops-team/trellis-go/pkg/fsutil"
)

// Syncer syncs embedded templates to the project .trellis/ directory.
type Syncer struct {
	EmbedFS   fs.FS
	TargetDir string
	Config    SyncerConfig
}

// SyncerConfig holds update configuration.
type SyncerConfig struct {
	Skip    []string          // Paths to skip
	Migrate bool              // Force overwrite all
	DryRun  bool              // Preview only
	Sections []ConfigSection  // Config sections to append
}

// ConfigSection represents a sentinel-guarded config section to append.
type ConfigSection struct {
	Sentinel string // Unique text to check for existence
	Content  string // Content to append if sentinel is missing
}

// SyncResult holds the result of a sync operation.
type SyncResult struct {
	Added    []string
	Updated  []string
	Skipped  []string
	Sections []string
}

// Sync performs template synchronization.
func (s *Syncer) Sync() (*SyncResult, error) {
	result := &SyncResult{}

	if err := fsutil.EnsureDir(s.TargetDir); err != nil {
		return nil, fmt.Errorf("ensure target dir: %w", err)
	}

	// Walk embedded templates
	err := fs.WalkDir(s.EmbedFS, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}

		// Check skip list
		if s.shouldSkip(path) {
			result.Skipped = append(result.Skipped, path)
			return nil
		}

		targetPath := filepath.Join(s.TargetDir, path)

		// Read embedded content
		srcData, err := fs.ReadFile(s.EmbedFS, path)
		if err != nil {
			return fmt.Errorf("read embedded %s: %w", path, err)
		}

		// Check if target exists
		targetData, err := os.ReadFile(targetPath)
		if os.IsNotExist(err) {
			// New file
			if !s.Config.DryRun {
				if err := s.writeFile(targetPath, srcData); err != nil {
					return err
				}
			}
			result.Added = append(result.Added, path)
			return nil
		}
		if err != nil {
			return fmt.Errorf("read target %s: %w", targetPath, err)
		}

		// Compare hashes
		srcHash := hashBytes(srcData)
		targetHash := hashBytes(targetData)

		if srcHash == targetHash {
			// Same content, skip
			return nil
		}

		// Content differs
		if s.Config.Migrate {
			// Force overwrite
			if !s.Config.DryRun {
				if err := s.writeFile(targetPath, srcData); err != nil {
					return err
				}
			}
			result.Updated = append(result.Updated, path)
		} else {
			// User may have modified, skip
			result.Skipped = append(result.Skipped, path)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("walk templates: %w", err)
	}

	// Append config sections
	for _, sec := range s.Config.Sections {
		configPath := filepath.Join(s.TargetDir, "config.yaml")
		configData, err := os.ReadFile(configPath)
		if err != nil {
			continue
		}
		if strings.Contains(string(configData), sec.Sentinel) {
			continue // Already present
		}
		if !s.Config.DryRun {
			newData := string(configData) + "\n" + sec.Content + "\n"
			if err := os.WriteFile(configPath, []byte(newData), 0644); err != nil {
				return nil, fmt.Errorf("append config section: %w", err)
			}
		}
		result.Sections = append(result.Sections, sec.Sentinel)
	}

	return result, nil
}

// DryRun performs a dry-run sync and returns what would change.
func (s *Syncer) DryRun() (*SyncResult, error) {
	s.Config.DryRun = true
	return s.Sync()
}

// Migrate performs a forced migration, overwriting all files.
func (s *Syncer) Migrate() (*SyncResult, error) {
	s.Config.Migrate = true
	return s.Sync()
}

func (s *Syncer) shouldSkip(path string) bool {
	for _, skip := range s.Config.Skip {
		if path == skip || strings.HasPrefix(path, skip+"/") {
			return true
		}
	}
	return false
}

func (s *Syncer) writeFile(path string, data []byte) error {
	dir := filepath.Dir(path)
	if err := fsutil.EnsureDir(dir); err != nil {
		return fmt.Errorf("ensure dir %s: %w", dir, err)
	}
	return os.WriteFile(path, data, 0644)
}

func hashBytes(data []byte) string {
	h := sha256.Sum256(data)
	return fmt.Sprintf("%x", h)
}

// Ensure io is used (import check)
var _ = io.Discard
