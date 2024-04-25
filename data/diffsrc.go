package data

import (
	"errors"
	"os"
	"path"

	"github.com/ogios/merge-repo/config"
)

func initDiffsrc() {
	e := errors.New("diff src dir not exist")
	f, err := os.Stat(config.GlobalConfig.DiffSrc)
	if err != nil {
		panic(e)
	}
	if !f.IsDir() {
		panic(e)
	}
}

func GetDiffSrcFile(p string) ([]byte, error) {
	fp := path.Join(config.GlobalConfig.DiffSrc, p)
	return os.ReadFile(fp)
}
