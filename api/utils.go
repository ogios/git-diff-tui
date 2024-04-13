package api

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
)

func MatchRegex(ss []string, regex string) []string {
	var res []string
	for _, v := range ss {
		y, _ := regexp.MatchString(regex, v)
		if y {
			res = append(res, v)
		}
	}
	return res
}

func EnsureExist(path string) error {
	return os.MkdirAll(path, 0755)
}

func MoveFile(srcPath, destPath string) error {
	srcFile, err := os.Open(srcPath)
	if err != nil {
		return err
	}
	defer srcFile.Close()

	destFile, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer destFile.Close()

	_, err = io.Copy(destFile, srcFile)
	if err != nil {
		return err
	}

	return nil
}

func CopyFiles(fs []string, from, to string) {
	from, err := filepath.Abs(from)
	if err != nil {
		panic(err)
	}
	to, err = filepath.Abs(to)
	if err != nil {
		panic(err)
	}
	fmt.Println(from, "->", to)

	removed := make([]string, 0)
	for _, v := range fs {
		from := filepath.Join(from, v)
		to := filepath.Join(to, v)
		err := EnsureExist(filepath.Dir(to))
		if err != nil {
			panic(err)
		}
		err = MoveFile(from, to)
		if err != nil {
			if os.IsNotExist(err) {
				removed = append(removed, v)
			} else {
				panic(err)
			}
		}
	}

	res, err := ExecCmd("tree", to)
	if err != nil {
		return
	}
	fmt.Println(string(res))
	rlog := "removed(%d): "
	rlogf := make([]any, len(removed)+1)
	rlogf[0] = len(removed)
	for i, v := range removed {
		rlog += "\n%s"
		rlogf[i+1] = v
	}
	fmt.Printf(rlog, rlogf...)
}

func SliceFrom[S ~[]E, E comparable](src S, start, end int) S {
	inputLen := len(src)
	if start < inputLen {
		if end < inputLen {
			s := src[start:end]
			to := make(S, len(s))
			copy(to, s)
			return to
		} else {
			s := src[start:]
			to := make(S, len(s))
			copy(to, s)
			return to
		}
	}
	return make(S, 0)
}
