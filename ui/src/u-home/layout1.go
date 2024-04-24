package uhome

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ogios/merge-repo/ui/comp"
	utree "github.com/ogios/merge-repo/ui/src/u-tree"
	uview "github.com/ogios/merge-repo/ui/src/u-view"
)

type Home struct {
	HomeCore
}

func newHome() *Home {
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

	modelCount := 2
	modelsHeight := h - 1
	modelsWidth := w - 2*modelCount
	getModelWidth := modelWidthCounter(modelCount, modelsWidth)
	ms := []tea.Model{
		utree.NewTreeModel(comp.TREE_NODE, [2]int{
			getModelWidth(0.2),
			modelsHeight,
		}),
		uview.NewViewModel([2]int{
			getModelWidth(0.8),
			modelsHeight,
		}),
	}

	home := &Home{
		HomeCore: HomeCore{
			Models: ms,
			Tree:   ms[0],
			Text:   ms[1].(*uview.ViewModel),
		},
	}

	return home
}

func (m *Home) Init() tea.Cmd {
	cmds := make([]tea.Cmd, 0)
	for _, m2 := range m.Models {
		cmds = append(cmds, m2.Init())
	}
	return tea.Batch(cmds...)
}

func (m *Home) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	update(msg, &m.HomeCore)
	return m, tea.Batch(cmds...)
}

func (m *Home) View() string {
	return view(&m.HomeCore)
}
