package left

import (
	"strings"

	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ogios/merge-repo/api"
	"github.com/ogios/merge-repo/data"
)

type CommentsModel struct {
	cache api.StringCacher
	v     viewport.Model
}

var commentIdentifierStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#5575e6"))

func NewCommentsModel(block [2]int) tea.Model {
	view := &CommentsModel{
		v: viewport.New(block[0], block[1]),
		cache: *api.NewStringCacher(func(k string) []byte {
			s := strings.Builder{}
			content, err := data.GetDiffFileComment(k)
			if err != nil {
				s.WriteString(errorMsg.Render(err.Error()))
			}
			for _, v := range content {
				s.WriteString(commentIdentifierStyle.Render("◌ "))
				s.WriteString(v)
				s.WriteString("\n")
			}
			return []byte(s.String())
		}),
	}
	return view
}

func (c *CommentsModel) ViewComment(f string) {
	content := c.cache.Get(f)
	c.v.SetContent(string(content))
	c.v.GotoTop()
}

func (c *CommentsModel) Init() tea.Cmd {
	return c.v.Init()
}

func (c *CommentsModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	c.v, cmd = c.v.Update(msg)
	return c, cmd
}

func (c *CommentsModel) View() string {
	return c.v.View()
}
