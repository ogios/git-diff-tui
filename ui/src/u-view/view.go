package uview

import (
	"bytes"
	"path"

	"github.com/alecthomas/chroma/v2/lexers"
	"github.com/alecthomas/chroma/v2/quick"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ogios/cropviewport"
	"github.com/ogios/merge-repo/api"
	"github.com/ogios/merge-repo/data"
	"github.com/ogios/merge-repo/ui/comp"
)

type ViewModel struct {
	cache *api.ContentCacher[*comp.ContentData]
	v     tea.Model
	path  string
}

var viewBlockStyle lipgloss.Style

func NewViewModel(block [2]int) tea.Model {
	viewBlockStyle = lipgloss.NewStyle().Width(block[0]).Height(block[1])
	view := &ViewModel{
		v: cropviewport.NewCropViewportModel(),
		cache: api.NewContentCacher(func(p string) *comp.ContentData {
			var finalContent string
			content, err := data.GetTempDiffFile(p)
			if err != nil {
				finalContent = comp.ErrorMsgStyle.Render(err.Error())
			} else {
				buf := new(bytes.Buffer)
				lex := lexers.Match(path.Base(p))
				lang := "plaintext"
				if lex != nil {
					lang = lex.Config().Name
				}
				err = quick.Highlight(buf, string(content), lang, "terminal16m", "catppuccin-mocha")
				if err != nil {
					finalContent = comp.ErrorMsgStyle.Render(err.Error())
				} else {
					finalContent = buf.String()
				}
			}
			at, sl := cropviewport.ProcessContent(finalContent)
			return &comp.ContentData{
				Table: at,
				Lines: sl,
			}
		}),
	}
	view.v.(*cropviewport.CropViewportModel).SetBlock(0, 0, block[0], block[1])
	return view
}

func (v *ViewModel) ViewFile(p string) {
	content := v.cache.Get(p)
	cv := v.v.(*cropviewport.CropViewportModel)
	cv.SetContentGivenData(content.Table, content.Lines)
	cv.BackToTop()
	cv.BackToLeft()
}

func (v *ViewModel) SetFile(p string) tea.Cmd {
	v.path = p
	return func() tea.Msg {
		content := v.cache.Get(p)
		cv := v.v.(*cropviewport.CropViewportModel)
		if v.path != p {
			return nil
		}
		cv.SetContentGivenData(content.Table, content.Lines)
		cv.BackToTop()
		cv.BackToLeft()
		return 1
	}
}

func (v *ViewModel) Init() tea.Cmd {
	return nil
}

func (v *ViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.MouseMsg:
		mouse := tea.MouseEvent(msg)
		cv := v.v.(*cropviewport.CropViewportModel)
		switch mouse.Button {
		case tea.MouseButtonWheelUp:
			if mouse.Ctrl {
				cv.PrevCol(1)
			} else {
				cv.PrevLine(1)
			}
		case tea.MouseButtonWheelDown:
			if mouse.Ctrl {
				cv.NextCol(1)
			} else {
				cv.NextLine(1)
			}
		case tea.MouseButtonWheelLeft:
			cv.PrevCol(1)
		case tea.MouseButtonWheelRight:
			cv.NextCol(1)
		}
		return v, nil
	}
	vp, cmd := v.v.Update(msg)
	v.v = vp
	// cmds = append(cmds, viewport.Sync(v.v), cmd)
	cmds = append(cmds, cmd)
	return v, tea.Batch(cmds...)
}

func (v *ViewModel) View() string {
	return viewBlockStyle.Render(v.v.View())
}
