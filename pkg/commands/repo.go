package commands

type GitRepoInfo struct {
	Path      string
	LastTag   string
	HasTag    bool
	HasOrigin bool
	NoPrefix  bool
	RemoteURL string // e.g. "https://github.com/user/repo"
}
