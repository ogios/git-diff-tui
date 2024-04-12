package config

import (
	"fmt"
	"os"

	"github.com/ogios/merge-repo/api"
)

type Config struct {
	Hash1, Hash2, Regex string
	ShowUI              bool
}

var GlobalConfig = Config{
	Hash1:  "",
	Hash2:  "",
	Regex:  "",
	ShowUI: true,
}

type ArgFunc = func(index int) int

func ParseArgs() {
	argFuncs := map[string]ArgFunc{
		"-n": func(i int) int {
			GlobalConfig.ShowUI = false
			return i
		},
	}
	normalFunc := getNormalArgs()
	i := 1
	for i < len(os.Args) {
		v := os.Args[i]
		f := argFuncs[v]
		if f != nil {
			i = f(i)
		} else {
			i = normalFunc(i)
		}
		i++
	}

	if (GlobalConfig.Hash1 == "") || (GlobalConfig.Hash2 == "") {
		useRefLog()
	}
}

func getNormalArgs() func(i int) int {
	parsed := false
	return func(i int) int {
		if !parsed {
			GlobalConfig.Hash1 = os.Args[i]
			i++
			GlobalConfig.Hash2 = os.Args[i]
			parsed = true
		} else {
			GlobalConfig.Regex = os.Args[i]
		}
		return i
	}
}

func useRefLog() {
	cs, err := api.GetReflogCommit()
	if err != nil {
		panic(err)
	}
	if len(cs) < 2 {
		panic(fmt.Errorf("not enough commits"))
	}
	GlobalConfig.Hash2 = cs[0]
	GlobalConfig.Hash1 = cs[len(cs)-1]
}
