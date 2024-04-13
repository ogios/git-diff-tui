package main

import (
	"fmt"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ogios/merge-repo/api"
	"github.com/ogios/merge-repo/config"
	"github.com/ogios/merge-repo/data"
	"github.com/ogios/merge-repo/template"
	"github.com/ogios/merge-repo/ui/left"
)

func main() {
	start := time.Now().UnixMilli()

	// args
	config.ParseArgs()

	// commits
	getCommits()

	// diffinfo html
	diffFileInfo()

	// ui
	if config.GlobalConfig.ShowUI {
		withUI()
	} else {
		noui()
	}

	fmt.Printf("cost: %dms\n", time.Now().UnixMilli()-start)
}

func getCommits() {
	hashes := &[2]string{
		config.GlobalConfig.Hash1, config.GlobalConfig.Hash2,
	}
	cs, err := api.GetCommits(hashes)
	if err != nil {
		panic(err)
	}
	data.COMMITS = cs
	// fmt.Println(cs)
	log.Println(cs)
}

// 生成diff表格
func diffFileInfo() {
	// data.DIFF_FILES := map[string][]api.Commit{}
	for i := 0; i < len(data.COMMITS)-1; i++ {
		old := data.COMMITS[i]
		change := data.COMMITS[i+1]
		fs, err := api.GetDiffFiles(old.Hash, change.Hash)
		log.Println(fs)
		fs = api.MatchRegex(fs, config.GlobalConfig.Regex)
		if err != nil {
			panic(err)
		}
		for _, f := range fs {
			data.DIFF_FILES[f] = append(data.DIFF_FILES[f], change)
		}
	}

	ls := make([][]string, len(data.DIFF_FILES))
	n := 0
	for fn, fcs := range data.DIFF_FILES {
		var comment string
		for _, c := range fcs {
			comment += c.Comment + "\n\r"
		}
		ls[n] = []string{
			fn,
			comment,
		}
		n++
	}
	template.GenTemplate(ls)
}

// 纯复制文件
func noui() {
	fs, err := api.GetDiffFiles(config.GlobalConfig.Hash1, config.GlobalConfig.Hash2)
	if err != nil {
		panic(err)
	}
	fs = api.MatchRegex(fs, config.GlobalConfig.Regex)
	api.CopyFiles(fs, "../..", "./copies")
}

// tui 未完成
func withUI() {
	// p := left.NewHomeModel()
	// fmt.Println(p.View())

	p := tea.NewProgram(left.NewHomeModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	// node := left.GetTreeNodes()
	// fmt.Println(node)
}
