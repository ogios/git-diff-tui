package config

import (
	"errors"
	"os"

	"github.com/ogios/merge-repo/api"
)

type Config struct {
	Hash1, Hash2, Regex string
	ShowUI, AltScreen   bool
}

var GlobalConfig = Config{
	Hash1:     "",
	Hash2:     "",
	Regex:     "",
	ShowUI:    true,
	AltScreen: true,
}

type ArgFunc = func(index int) int

func ParseArgs() {
	onlyReg := false
	normalArgs := []string{}
	i := 1
	for i < len(os.Args) {
		switch os.Args[i] {
		case "-n":
			GlobalConfig.ShowUI = false
		case "-a":
			GlobalConfig.AltScreen = false
		case "-r":
			useRefLog()
			onlyReg = true
		default:
			normalArgs = append(normalArgs, os.Args[i])
		}
		i++
	}

	i = 0
	if !onlyReg {
		GlobalConfig.Hash1 = GetStrIfExist(normalArgs, i)
		i++
		GlobalConfig.Hash2 = GetStrIfExist(normalArgs, i)
		i++
	}
	GlobalConfig.Regex = GetStrIfExist(normalArgs, i)
}

func useRefLog() {
	cs, err := api.GetReflogCommit()
	if err != nil {
		panic(err)
	}
	if len(cs) < 2 {
		panic(errors.New("not enough commits"))
	}
	GlobalConfig.Hash2 = cs[0]
	GlobalConfig.Hash1 = cs[len(cs)-1]
}

func GetStrIfExist(l []string, i int) string {
	if i >= len(l) {
		return ""
	}
	return l[i]
}
