package left

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ogios/cropviewport"
	"github.com/ogios/cropviewport/process"
	"github.com/ogios/merge-repo/api"
	"github.com/ogios/merge-repo/data"
)

type CommentsModel struct {
	cache *api.ContentCacher[*api.ContentData]
	v     tea.Model
}

var commentIdentifierStyle = lipgloss.NewStyle().
	Foreground(lipgloss.Color("#5575e6"))

func NewCommentsModel(block [2]int) tea.Model {
	view := &CommentsModel{
		v: cropviewport.NewCropViewportModel(),
		cache: api.NewContentCacher(func(k string) *api.ContentData {
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
			at, sl := process.ProcessContent(s.String())
			return &api.ContentData{
				Table: at,
				Lines: sl,
			}
		}),
	}
	return view
}

func (c *CommentsModel) ViewComment(f string) {
	content := c.cache.Get(f)
	cv := c.v.(*cropviewport.CropViewportModel)
	cv.SetContentGivenData(content.Table, content.Lines)
	cv.BackToTop()
	cv.BackToLeft()
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
