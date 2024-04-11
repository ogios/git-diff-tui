package api

import (
	"fmt"
	"os/exec"
	"strings"
)

/*
`git diff` 传1-2个hash

传一个则为获取某个commit里修改的内容

传两个则为获取从第一个commit至第二个commit之间的所有修改文件
*/
func GetDiffFiles(hashes ...string) ([]string, error) {
	if len(hashes) < 1 || len(hashes) > 2 {
		return nil, fmt.Errorf("wrong hashes length: %v", hashes)
	}
	args := []string{
		"git",
		"diff",
		"--name-only",
	}
	args = append(args, hashes...)
	res := exec.Command(args[0], args[1:]...)
	out, err := res.Output()
	if err != nil {
		panic(fmt.Errorf("git diff error: %v", err))
	}
	return strings.Split(string(out), "\n"), nil
}
