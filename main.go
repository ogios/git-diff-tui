package main

import (
	"fmt"
	"log"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/ogios/merge-repo/api"
	"github.com/ogios/merge-repo/config"
	"github.com/ogios/merge-repo/template"
	ui "github.com/ogios/merge-repo/ui/main"
)

var COMMITS []api.Commit

func main() {
	start := time.Now().UnixMilli()

	// args
	config.ParseArgs()
	hashes := &[2]string{
		config.GlobalConfig.Hash1, config.GlobalConfig.Hash2,
	}
	// h1 := os.Args[1]
	// h2 := os.Args[2]
	// reg := func() string {
	// 	if len(os.Args) >= 4 {
	// 		return os.Args[3]
	// 	}
	// 	return ""
	// }()

	// commits
	cs, err := api.GetCommits(hashes)
	if err != nil {
		panic(err)
	}
	COMMITS = cs
	// fmt.Println(cs)
	log.Println(cs)

	// operations
	diffFileInfo(config.GlobalConfig.Regex)
	if len(os.Args) >= 5 && os.Args[4] == "-n" {
		noui(hashes[0], hashes[1], config.GlobalConfig.Regex)
		return
	}
	fmt.Printf("cost: %dms\n", time.Now().UnixMilli()-start)
}

// 生成diff表格
func diffFileInfo(reg string) {
	m := map[string][]api.Commit{}
	for _, c := range COMMITS {
		fs, err := api.GetDiffFiles(c.Hash)
		fs = api.MatchRegex(fs, reg)
		if err != nil {
			panic(err)
		}
		for _, f := range fs {
			m[f] = append(m[f], c)
		}
	}

	ls := make([][]string, len(m))
	n := 0
	for fn, fcs := range m {
		var comment string
		for _, c := range fcs {
			comment += c.Comment + "\n\r"
		}
		ls[n] = []string{
			fn,
			comment,
			// c.Hash,
			// c.Comment,
			// c.Time,
			// c.Tag,
			// c.Author,
		}
		n++
	}
	template.GenTemplate(ls)
}

// 纯复制文件
func noui(h1, h2, reg string) {
	fs, err := api.GetDiffFiles(h1, h2)
	if err != nil {
		panic(err)
	}
	fs = api.MatchRegex(fs, reg)
	api.CopyFiles(fs, "../..", "./copies")
}

// tui 未完成
func withUI() {
	p := tea.NewProgram(ui.NewHomeModel())
	if _, err := p.Run(); err != nil {
		log.Fatal(err)
	}
}
