package main

import (
	"fmt"
	"log"
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ogios/merge-repo/api"
	"github.com/ogios/merge-repo/config"
	"github.com/ogios/merge-repo/data"
	"github.com/ogios/merge-repo/template"
	_ "github.com/ogios/merge-repo/ui/comp"
	uhome "github.com/ogios/merge-repo/ui/src/u-home"
)

func main() {
	defer config.Exit()
	defer data.Exit()
	start := time.Now().UnixMilli()

	// html
	GenDiffHTML()

	fmt.Printf("cost: %dms\n", time.Now().UnixMilli()-start)

	// ui
	if config.GlobalConfig.ShowUI {
		withUI()
	} else {
		noUI()
	}
}

// 生成diff表格
func GenDiffHTML() {
	ls := make([][]string, len(data.DIFF_FILES))
	n := 0
	for fn := range data.DIFF_FILES {
		s, err := data.GetDiffFileComment(fn)
		if err != nil {
			panic(err)
		}
		ls[n] = []string{
			fn,
			strings.Join(s, "\n"),
		}
		n++
	}
	template.GenTemplate(ls)
}

// 纯复制文件
func noUI() {
	fs, err := api.GetDiffFiles(config.GlobalConfig.Hash1, config.GlobalConfig.Hash2)
	if err != nil {
		panic(err)
	}
	fs = api.MatchRegex(fs, config.GlobalConfig.Regex)
	data.CopyFiles(fs)
}

// tui 未完成
func withUI() {
	// v := uview.NewViewModel([2]int{154, 37}).(*uview.ViewModel)
	// v.ViewFile("README.md")
	// fmt.Println(v.View())
	// p := udiffview.NewDiffViewModel([2]int{154, 37}).(*udiffview.DiffViewModel)
	// p.ViewFile("README.md")
	// fmt.Println(p.View())

	data.PROGRAM = tea.NewProgram(uhome.NewHomeModel(), func() func(p *tea.Program) {
		if config.GlobalConfig.AltScreen {
			return tea.WithAltScreen()
		}
		return func(p *tea.Program) {}
	}())
	if _, err := data.PROGRAM.Run(); err != nil {
		log.Fatal(err)
	}
}
