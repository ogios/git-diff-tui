package api

import (
	"fmt"
	"os/exec"
	"strings"
)

type Commit struct {
	Hash, Author, Time, Comment, Tag string
}

func parseCommits(raw string) Commit {
	l := strings.Split(raw, "\n")
	return Commit{
		Hash:    l[0],
		Author:  l[1],
		Time:    strings.Replace(l[2], " +0800", "", 1),
		Comment: l[3],
	}
}

func GetCommitLog(hashes *[2]string) ([]string, error) {
	args := []string{
		"log",
		"--pretty=format:%h%n%an%n%ci%n%s%n",
	}
	if hashes != nil {
		args = append(args, hashes[0]+".."+hashes[1])
	}
	cmd := exec.Command("git", args...)
	o, err := cmd.Output()
	if err != nil {
		return nil, fmt.Errorf("git log error: %v", err)
	}
	rcs := strings.Split(strings.TrimSpace(string(o)), "\n\n")
	return rcs, nil
}

func GetTags() map[string]string {
	cmd := exec.Command("git", "show-ref", "--tags")
	o, err := cmd.Output()
	m := make(map[string]string)
	if err != nil {
		return m
	}
	rt := strings.Split(strings.TrimSpace(string(o)), "\n")
	for _, v := range rt {
		ct := strings.Split(v, " ")
		m[ct[0][:7]] = strings.Replace(ct[1], "refs/tags/", "", 1)
	}
	return m
}

func GetCommits(hashes *[2]string) ([]Commit, error) {
	rcs, err := GetCommitLog(hashes)
	if err != nil {
		return nil, fmt.Errorf("git commit error: %v", err)
	}
	rts := GetTags()
	cs := make([]Commit, len(rcs))
	for i, v := range rcs {
		v = strings.TrimSpace(v)
		if v != "" {
			c := parseCommits(v)
			t := rts[c.Hash]
			c.Tag = t
			cs[i] = c
		}
	}
	return cs, nil
}

// func GetCommitsBetween(hash1, hash2 string) ([]string, error) {
// 	cmd := exec.Command("git", "log", "--pretty=format:%h", hash1+".."+hash2)
// 	o, err := cmd.Output()
// 	if err != nil {
// 		return nil, err
// 	}
// 	return strings.Split(string(o), "\n"), nil
// }
