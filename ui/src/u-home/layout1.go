package uhome

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ogios/merge-repo/ui/comp"
	utree "github.com/ogios/merge-repo/ui/src/u-tree"
	uview "github.com/ogios/merge-repo/ui/src/u-view"
)

type Home struct {
	comment tea.Model
	HomeCore
}

func newHome() *Home {
	w := comp.GlobalUIData.MaxWidth - 2
	h := comp.GlobalUIData.MaxHeight - 2

	modelCount := 2
	modelsHeight := h - 1
	modelsWidth := w - 2*modelCount
	getModelWidth := modelWidthCounter(modelCount, modelsWidth)
	ms := []*childModel{
		newChild([2]int{getModelWidth(0.2), modelsHeight}),
		newChild([2]int{getModelWidth(0.8), modelsHeight}),
	}
	ms[0].m = utree.NewTreeModel(comp.TREE_NODE, ms[0].block)
	ms[1].m = uview.NewViewModel(ms[1].block)

	home := &Home{
		HomeCore: HomeCore{
			Models: ms,
			Text:   ms[1].m.(*uview.ViewModel),
		},
	}

	return home
}

func (m *Home) Init() tea.Cmd {
	cmds := make([]tea.Cmd, 0)
	for _, m2 := range m.Models {
		cmds = append(cmds, m2.m.Init())
	}
	return tea.Batch(cmds...)
}

func (m *Home) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, update(msg, &m.HomeCore)
}

func (m *Home) View() string {
	return view(&m.HomeCore)
}
