package data

import (
	"fmt"
	"os"
	"path"
	"regexp"

	"github.com/ogios/merge-repo/api"
	"github.com/ogios/merge-repo/config"
)

var (
	TEMP_PATH      string
	TEMP_DATA_PATH string
)

func initTemp() {
	p, err := os.MkdirTemp(os.TempDir(), "git-diff-tui-*")
	if err != nil {
		panic(err)
	}
	TEMP_PATH = p
	TEMP_DATA_PATH = path.Join(TEMP_PATH, "data")
	pattern, err := regexp.Compile("^fatal: path '.*?' does not exist in '.*?'.*")
	if err != nil {
		panic(err)
	}

	for p := range DIFF_FILES {
		c, err := api.GetGitFile(config.GlobalConfig.Hash2, p)
		if err != nil {
			if !pattern.Match(c) {
				panic(fmt.Errorf("error: init temp diff files: %s\n %v", string(c), err))
			}
			continue
		}
		saveTempDiffFile(p, c)
	}
}

func exitTemp() {
	fmt.Println("removeing temp")
	err := os.RemoveAll(TEMP_PATH)
	if err != nil {
		panic(err)
	}
}

func saveTempDiffFile(p string, data []byte) error {
	fp := path.Join(TEMP_DATA_PATH, p)
	err := api.EnsureExist(path.Dir(fp))
	if err != nil {
		return err
	}
	err = os.WriteFile(fp, data, 0766)
	if err != nil {
		return err
	}
	return nil
}

func GetTempDiffFile(p string) ([]byte, error) {
	fp := path.Join(TEMP_DATA_PATH, p)
	return os.ReadFile(fp)
}

func CopyFiles(fs []string) {
	api.CopyFiles(fs, TEMP_DATA_PATH, "./copies")
}
