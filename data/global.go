package data

import (
	"errors"
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

func initGlobal() {
	p, err := api.ExecCmd("git", "rev-parse", "--show-toplevel")
	if err != nil {
		panic(err)
	}
	BASE_PATH = strings.TrimSpace(string(p))
	getDiffFileInfo()
}

func getCommits() {
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

func getDiffFileInfo() {
	if COMMITS == nil {
		getCommits()
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

func GetDiffFileComment(f string) ([]string, error) {
	if df, ok := DIFF_FILES[f]; ok {
		fcs := make([]string, len(df))
		for i, c := range df {
			fcs[i] = c.Comment
		}
		return fcs, nil
	}
	return nil, errors.New("diff file not found")
}
