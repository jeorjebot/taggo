package commands

type GitRepoInfo struct {
	Path      string
	LastTag   string
	HasTag    bool
	HasOrigin bool
	NoPrefix  bool
}
