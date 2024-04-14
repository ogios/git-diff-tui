package left

import (
	"bytes"
	"log"
	"path"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/quick"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ogios/merge-repo/api"
	"github.com/ogios/merge-repo/config"
)

type ViewModel struct {
	cache api.StringCacher
	v     viewport.Model
}

var errorMsg = lipgloss.NewStyle().
	Bold(true).
	Inline(true).
	Background(lipgloss.Color("#ff2f5a")).
	Foreground(lipgloss.Color("#000000"))

func NewViewModel(block [2]int) tea.Model {
	view := &ViewModel{
		v: viewport.New(block[0], block[1]),
		cache: *api.NewStringCacher(func(k string) []byte {
			content, err := api.GetGitFile(config.GlobalConfig.Hash2, k)
			if err != nil {
				return []byte(errorMsg.Render(err.Error()))
			}
			buf := new(bytes.Buffer)
			lex := lexers.Match(path.Base(k)).Config().Name
			log.Println(lex, path.Base(k))
			err = quick.Highlight(buf, content, lex, "terminal16m", "catppuccin-mocha")
			if err != nil {
				return []byte(errorMsg.Render(err.Error()))
			}
			return buf.Bytes()
			// return []byte(content)
		}),
	}
	return view
}

func (v *ViewModel) ViewFile(p string) {
	content := v.cache.Get(p)
	log.Println(string(content))
	v.v.SetContent(string(content))
	v.v.GotoTop()
}

func (v *ViewModel) Init() tea.Cmd {
	return nil
}

func (v *ViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	vp, cmd := v.v.Update(msg)
	v.v = vp
	// cmds = append(cmds, viewport.Sync(v.v), cmd)
	cmds = append(cmds, cmd)
	return v, tea.Batch(cmds...)
}

func (v *ViewModel) View() string {
	return v.v.View()
}
