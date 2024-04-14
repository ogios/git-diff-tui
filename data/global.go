package data

import (
	"log"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ogios/merge-repo/api"
	"github.com/ogios/merge-repo/config"
)

type DiffCommitMap map[string][]api.Commit

var (
	COMMITS    []api.Commit  = nil
	DIFF_FILES DiffCommitMap = DiffCommitMap{}
	BASE_PATH  string
	PROGRAM    *tea.Program
)

func GetCommits() {
	hashes := &[2]string{
		config.GlobalConfig.Hash1, config.GlobalConfig.Hash2,
	}
	cs, err := api.GetCommits(hashes)
	if err != nil {
		panic(err)
	}
	COMMITS = cs
	log.Println(cs)
}

func GetDiffFileInfo() {
	if COMMITS == nil {
		GetCommits()
	}
	for i := 0; i < len(COMMITS)-1; i++ {
		old := COMMITS[i]
		change := COMMITS[i+1]
		fs, err := api.GetDiffFiles(old.Hash, change.Hash)
		log.Println(fs)
		fs = api.MatchRegex(fs, config.GlobalConfig.Regex)
		if err != nil {
			panic(err)
		}
		for _, f := range fs {
			DIFF_FILES[f] = append(DIFF_FILES[f], change)
		}
	}
}

func init() {
	p, err := api.ExecCmd("git", "rev-parse", "--show-toplevel")
	if err != nil {
		panic(err)
	}
	BASE_PATH = strings.TrimSpace(p)
}

func CopyFiles(fs []string) {
	api.CopyFiles(fs, BASE_PATH, "./copies")
}
