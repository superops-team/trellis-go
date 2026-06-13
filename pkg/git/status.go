package git

// Status represents the repository status.
type Status struct {
	Branch      string
	IsDirty     bool
	Ahead       int
	Behind      int
	StagedFiles []string
}
