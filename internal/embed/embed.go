package embed

import "embed"

// Templates contains all embedded template files.
//go:embed templates/*
var Templates embed.FS
