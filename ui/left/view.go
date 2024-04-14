package left

import (
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
				content = errorMsg.Render(err.Error())
			}
			return []byte(content)
		}),
	}
	return view
}

func (v *ViewModel) ViewFile(p string) {
	content := v.cache.Get(p)
	v.v.SetContent(string(content))
	v.v.GotoTop()
}

func (v *ViewModel) Init() tea.Cmd {
	return v.v.Init()
}

func (v *ViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	v.v, cmd = v.v.Update(msg)
	return v, cmd
}

func (v *ViewModel) View() string {
	return v.v.View()
}
