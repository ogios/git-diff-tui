package api

import (
	"os"
)

func ReadFile(p string) ([]byte, error) {
	return os.ReadFile(p)
}

func GetGitFile(hash, p string) (string, error) {
	return ExecCmd("git", "show", hash+":"+p)
}

type StringCacher struct {
	pool map[string][]byte
	new  func(p string) []byte
}

func (c *StringCacher) Get(key string) []byte {
	if s, ok := c.pool[key]; ok {
		return s
	}
	s := c.new(key)
	c.pool[key] = s
	return s
}

func NewStringCacher(new func(p string) []byte) *StringCacher {
	return &StringCacher{pool: map[string][]byte{}, new: new}
}
