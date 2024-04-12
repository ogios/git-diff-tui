package api

import (
	"fmt"
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
	baseArgs := []string{
		"git",
		"log",
		"--pretty=format:%h%n%an%n%ci%n%s%n",
	}
	var o string
	if hashes != nil {
		args := make([]string, len(baseArgs)+1)
		copy(args, baseArgs)
		args[len(args)-1] = hashes[0]
		r, err := ExecCmd(args...)
		if err != nil {
			return nil, fmt.Errorf("git log error: %v", err)
		}
		o += r + "\n"
		args[len(args)-1] = hashes[0] + ".." + hashes[1]
		r, err = ExecCmd(args...)
		if err != nil {
			return nil, fmt.Errorf("git log error: %v", err)
		}
		o += r
	} else {
		r, err := ExecCmd(baseArgs...)
		if err != nil {
			return nil, fmt.Errorf("git log error: %v", err)
		}
		o = r
	}
	rcs := strings.Split(strings.TrimSpace(o), "\n\n")
	return rcs, nil
}

func GetTags() map[string]string {
	o, err := ExecCmd("git", "show-ref", "--tags")
	m := make(map[string]string)
	if err != nil {
		return m
	}
	rt := strings.Split(strings.TrimSpace(o), "\n")
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

func GetCurrentBranch() (string, error) {
	o, err := ExecCmd("git", "branch", "--show-current")
	if err != nil {
		return "", fmt.Errorf("get current branch error: %v", err)
	}
	return strings.TrimSpace(o), nil
}

func GetReflogCommit() ([]string, error) {
	bn, err := GetCurrentBranch()
	if err != nil {
		return nil, err
	}
	o, err := ExecCmd("git", "reflog", bn)
	if err != nil {
		return nil, fmt.Errorf("reflog error: %v", err)
	}

	rows := strings.Split(strings.TrimSpace(string(o)), "\n")
	cs := make([]string, len(rows))
	for i, v := range rows {
		cs[i] = strings.Split(v, " ")[0]
	}
	return cs, nil
}
