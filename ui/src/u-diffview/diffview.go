package udiffview

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ogios/cropviewport"
	"github.com/ogios/merge-repo/api"
	"github.com/ogios/merge-repo/ui/comp"
)

type DiffViewModel struct {
	cache *api.ContentCacher[*DiffContentData]
	v     tea.Model
}
type DiffContentData comp.ContentData

var (
	errorMsg = lipgloss.NewStyle().
			Bold(true).
			Inline(true).
			Background(lipgloss.Color("#ff2f5a")).
			Foreground(lipgloss.Color("#000000"))
	diffInsertBg = lipgloss.NewStyle().
			Background(lipgloss.Color("#1aff66"))
	diffDeleteBg = lipgloss.NewStyle().
			Background(lipgloss.Color("#ff1a4d"))
	viewBlockStyle lipgloss.Style
)

func NewDiffViewModel(block [2]int) tea.Model {
	viewBlockStyle = lipgloss.NewStyle().Width(block[0]).Height(block[1])
	view := &DiffViewModel{
		v: cropviewport.NewCropViewportModel(),
		cache: api.NewContentCacher(func(p string) *DiffContentData {
			atl, sl, err := diffContent(p, p)
			if err != nil {
				atl, sl = cropviewport.ProcessContent(errorMsg.Render(err.Error()))
			} else if atl == nil || sl == nil {
				atl, sl = cropviewport.ProcessContent(errorMsg.Render("diff view diffcontent ansi tablelist and sublines are nil"))
			}
			return &DiffContentData{
				Table: atl,
				Lines: sl,
			}
		}),
	}
	view.v.(*cropviewport.CropViewportModel).SetBlock(0, 0, block[0], block[1])
	return view
}

func (v *DiffViewModel) ViewFile(p string) {
	log.Println("diffing file", p)
	content := v.cache.Get(p)
	cv := v.v.(*cropviewport.CropViewportModel)
	cv.SetContentGivenData(content.Table, content.Lines)
	cv.BackToTop()
	cv.BackToLeft()
}

func (v *DiffViewModel) Init() tea.Cmd {
	return nil
}

func (v *DiffViewModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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

func (v *DiffViewModel) View() string {
	return viewBlockStyle.Render(v.v.View())
}
