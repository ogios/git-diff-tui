package left

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ogios/merge-repo/api"
	"github.com/ogios/merge-repo/data"
	"github.com/ogios/merge-repo/ui/comp"
)

type HomeModel struct {
	Tree            *TreeModel
	Text            *ViewModel
	Comment         *CommentsModel
	CurrentFile     string
	Models          []tea.Model
	FocusModelIndex int
}

var homeStyle, focusStyle, unfocusStyle lipgloss.Style

func GetTreeNodes() *api.Node {
	var node *api.Node = nil
	for k := range data.DIFF_FILES {
		node = api.PathToNode(k, node)
	}
	fmt.Println(node)
	return node
}

func NewHomeModel() *HomeModel {
	w := comp.GlobalUIData.MaxWidth - 2
	h := comp.GlobalUIData.MaxHeight - 2

	homeStyle = lipgloss.NewStyle().
		Width(w).
		Height(h).
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("69"))

	focusStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#ffad00"))
	unfocusStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#ff5b00"))

	modelCount := 3
	modelsHeight := h - 1
	modelsWidth := w - 2*modelCount
	getModelWidth := func() func(per float64) int {
		count := modelCount
		avaliable := modelsWidth
		return func(per float64) int {
			count--
			if count == 0 {
				return avaliable
			} else {
				r := int(float64(modelsWidth) * per)
				avaliable -= r
				return r
			}
		}
	}()
	ms := []tea.Model{
		NewTreeModel(GetTreeNodes(), [2]int{
			getModelWidth(0.2),
			modelsHeight,
		}),
		NewViewModel([2]int{
			getModelWidth(0.4),
			modelsHeight,
		}),
		NewCommentsModel([2]int{
			getModelWidth(0.4),
			modelsHeight,
		}),
	}

	home := &HomeModel{
		Models:          ms,
		Tree:            ms[0].(*TreeModel),
		Text:            ms[1].(*ViewModel),
		Comment:         ms[2].(*CommentsModel),
		FocusModelIndex: 0,
	}

	return home
}

func (m *HomeModel) Init() tea.Cmd {
	cmds := make([]tea.Cmd, 0)
	for _, m2 := range m.Models {
		cmds = append(cmds, m2.Init())
	}
	return tea.Batch(cmds...)
}

func (m *HomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	toFocusModel := false
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			fallthrough
		case "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.FocusModelIndex = ((m.FocusModelIndex + 1) + len(m.Models)) % len(m.Models)
		case "shift+tab":
			m.FocusModelIndex = ((m.FocusModelIndex - 1) + len(m.Models)) % len(m.Models)
		default:
			toFocusModel = true
		}
	case FileMsg:
		m.CurrentFile = msg.FileRelPath
		m.Text.ViewFile(m.CurrentFile)
		m.Comment.ViewComment(m.CurrentFile)
	case CopyFileMsg:
		data.PROGRAM.ReleaseTerminal()
		data.CopyFiles(msg.Files)
		data.PROGRAM.RestoreTerminal()
		return m, tea.Quit
	default:
		toFocusModel = true
	}
	if toFocusModel {
		_, cmd := m.Models[m.FocusModelIndex].Update(msg)
		cmds = append(cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m *HomeModel) View() string {
	var v string

	ms := make([]string, len(m.Models))
	for i, m2 := range m.Models {
		if m.FocusModelIndex == i {
			ms[i] = focusStyle.Render(m2.View())
		} else {
			ms[i] = unfocusStyle.Render(m2.View())
		}
	}
	v = lipgloss.JoinVertical(lipgloss.Left,
		m.CurrentFile,
		lipgloss.JoinHorizontal(lipgloss.Top,
			ms...,
		),
	)
	// return homeStyle.Render(v)
	return v
}
