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

	// commits & diffinfo
	data.GetDiffFileInfo()
	// html
	GenDiffHTML()

	// ui
	if config.GlobalConfig.ShowUI {
		withUI()
	} else {
		noui()
	}

	fmt.Printf("cost: %dms\n", time.Now().UnixMilli()-start)
}

// 生成diff表格
func GenDiffHTML() {
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
	api.CopyFiles(fs, data.BASE_PATH, "./copies")
}

// tui 未完成
func withUI() {
	// p := left.NewHomeModel()
	// fmt.Println(p.View())
	logger := config.CraeteLogger()
	defer logger.Close()

	p := tea.NewProgram(left.NewHomeModel(), tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}

	// node := left.GetTreeNodes()
	// fmt.Println(node)
}
