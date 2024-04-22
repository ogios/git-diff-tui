package left

import (
	"bytes"
	"path"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/quick"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ogios/cropviewport"
	"github.com/ogios/cropviewport/process"
	"github.com/ogios/merge-repo/api"
	"github.com/ogios/merge-repo/config"
)

type ViewModel struct {
	cache *api.ContentCacher[*api.ContentData]
	v     tea.Model
}

var errorMsg = lipgloss.NewStyle().
	Bold(true).
	Inline(true).
	Background(lipgloss.Color("#ff2f5a")).
	Foreground(lipgloss.Color("#000000"))

func NewViewModel(block [2]int) tea.Model {
	view := &ViewModel{
		v: cropviewport.NewCropViewportModel(),
		cache: api.NewContentCacher(func(p string) *api.ContentData {
			var finalContent string
			content, err := api.GetGitFile(config.GlobalConfig.Hash2, p)
			if err != nil {
				finalContent = errorMsg.Render(err.Error())
			} else {
				buf := new(bytes.Buffer)
				lex := lexers.Match(path.Base(p))
				lang := "plaintext"
				if lex != nil {
					lang = lex.Config().Name
				}
				err = quick.Highlight(buf, content, lang, "terminal16m", "catppuccin-mocha")
				if err != nil {
					finalContent = errorMsg.Render(err.Error())
				} else {
					finalContent = buf.String()
				}
			}
			at, sl := process.ProcessContent(finalContent)
			return &api.ContentData{
				Table: at,
				Lines: sl,
			}
		}),
	}
	return view
}

func (v *ViewModel) ViewFile(p string) {
	content := v.cache.Get(p)
	cv := v.v.(*cropviewport.CropViewportModel)
	cv.SetContentGivenData(content.Table, content.Lines)
	cv.BackToTop()
	cv.BackToLeft()
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
