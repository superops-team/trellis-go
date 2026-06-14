package main

import (
	"fmt"
	"os"
	"path/filepath"
)

type resolvedPaths struct {
	RepoRoot   string
	TrellisDir string
	TasksDir   string
	SpecDir    string
}

func resolveCommandPaths() (resolvedPaths, error) {
	return resolvePaths(root, os.Getwd)
}

func resolvePaths(rootFlag string, getwd func() (string, error)) (resolvedPaths, error) {
	repoRoot := rootFlag
	if repoRoot == "" {
		cwd, err := getwd()
		if err != nil {
			return resolvedPaths{}, fmt.Errorf("get working directory: %w", err)
		}
		repoRoot = cwd
	}

	repoRoot = filepath.Clean(repoRoot)
	trellisDir := filepath.Join(repoRoot, ".trellis")
	if filepath.Base(repoRoot) == ".trellis" {
		trellisDir = repoRoot
		repoRoot = filepath.Dir(repoRoot)
	}

	return resolvedPaths{
		RepoRoot:   repoRoot,
		TrellisDir: trellisDir,
		TasksDir:   filepath.Join(trellisDir, "tasks"),
		SpecDir:    filepath.Join(trellisDir, "spec"),
	}, nil
}
