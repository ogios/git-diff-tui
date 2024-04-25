package uhome

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ogios/merge-repo/data"
	"github.com/ogios/merge-repo/ui/comp"
	utree "github.com/ogios/merge-repo/ui/src/u-tree"
	uview "github.com/ogios/merge-repo/ui/src/u-view"
)

var homeStyle, focusStyle, unfocusStyle lipgloss.Style

func resetStyle() {
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
}

// type keyMap map[string]func(msg tea.Msg) tea.Cmd

type HomeCore struct {
	Tree            tea.Model
	Text            *uview.ViewModel
	CurrentFile     string
	Models          []tea.Model
	FocusModelIndex int
}

func update(msg tea.Msg, m *HomeCore) tea.Cmd {
	cmds := make([]tea.Cmd, 0)
	toFocusModel := false
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			fallthrough
		case "ctrl+c":
			return tea.Quit
		case "tab":
			m.FocusModelIndex = ((m.FocusModelIndex + 1) + len(m.Models)) % len(m.Models)
		case "shift+tab":
			m.FocusModelIndex = ((m.FocusModelIndex - 1) + len(m.Models)) % len(m.Models)
		default:
			toFocusModel = true
		}
	case utree.FileMsg:
		m.CurrentFile = msg.FileRelPath
		m.Text.ViewFile(m.CurrentFile)
		// m.Comment.ViewComment(m.CurrentFile)
	case utree.CopyFileMsg:
		data.PROGRAM.ReleaseTerminal()
		data.CopyFiles(msg.Files)
		data.PROGRAM.RestoreTerminal()
		return tea.Quit
	default:
		toFocusModel = true
	}
	if toFocusModel {
		_, cmd := m.Models[m.FocusModelIndex].Update(msg)
		cmds = append(cmds, cmd)
	}
	return tea.Batch(cmds...)
}

func view(m *HomeCore) string {
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

func modelWidthCounter(count int, total int) func(per float64) int {
	avaliable := total
	return func(per float64) int {
		count--
		if count == 0 {
			return avaliable
		} else {
			r := int(float64(total) * per)
			avaliable -= r
			return r
		}
	}
}
