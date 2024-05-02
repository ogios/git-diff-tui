package app

import (
	"fmt"
	"log"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ogios/merge-repo/api"
	"github.com/ogios/merge-repo/config"
	"github.com/ogios/merge-repo/data"
	_ "github.com/ogios/merge-repo/ui/comp"
	uhome "github.com/ogios/merge-repo/ui/src/u-home"
)

func RunApp() {
	defer config.Exit()
	defer data.Exit()
	start := time.Now().UnixMilli()

	fmt.Printf("cost: %dms\n", time.Now().UnixMilli()-start)

	// ui
	if config.GlobalConfig.ShowUI {
		withUI()
	} else {
		noUI()
	}
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
	// p.ViewFile("ui/src/u-home/layout2.go")
	// fmt.Println(p.View())

	data.PROGRAM = tea.NewProgram(uhome.NewHomeModel(), func() func(p *tea.Program) {
		if config.GlobalConfig.AltScreen {
			return tea.WithAltScreen()
		}
		return func(p *tea.Program) {}
	}(), tea.WithMouseAllMotion())
	if _, err := data.PROGRAM.Run(); err != nil {
		log.Fatal(err)
	}
}
