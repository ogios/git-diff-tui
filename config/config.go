package config

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/ogios/merge-repo/api"
)

type Config struct {
	Hash1, Hash2, Regex, DiffSrc string
	ShowUI, AltScreen            bool
}

var GlobalConfig = Config{
	Hash1:     "",
	Hash2:     "",
	Regex:     "",
	ShowUI:    true,
	AltScreen: true,
	DiffSrc:   "",
}

type grgFunc = func(index int) int

func initConfig() {
	start := time.Now().UnixMicro()
	parseArgs()
	log.Println("init config cost:", time.Now().UnixMicro()-start)
}

func parseArgs() {
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
		case "-d":
			i++
			GlobalConfig.DiffSrc = os.Args[i]
		default:
			normalArgs = append(normalArgs, os.Args[i])
		}
		i++
	}

	i = 0
	if !onlyReg {
		GlobalConfig.Hash1 = getStrIfExist(normalArgs, i)
		i++
		GlobalConfig.Hash2 = getStrIfExist(normalArgs, i)
		i++
	}
	GlobalConfig.Regex = getStrIfExist(normalArgs, i)
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

func getStrIfExist(l []string, i int) string {
	if i >= len(l) {
		return ""
	}
	return l[i]
}
