package data

import "github.com/ogios/merge-repo/api"

type DiffCommitMap map[string][]api.Commit

var (
	COMMITS    []api.Commit
	DIFF_FILES = DiffCommitMap{}
)
