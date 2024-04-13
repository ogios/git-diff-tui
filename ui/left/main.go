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
	Tree   *TreeModel
	Models []tea.Model
}

var homeStyle = lipgloss.NewStyle().
	Width(comp.GlobalUIData.MaxWidth - 2).
	Height(comp.GlobalUIData.MaxHeight - 2).
	// Align(lipgloss.Center, lipgloss.Center).
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("69"))

func getTreeNodes() *api.Node {
	var node *api.Node = nil
	for k := range data.DIFF_FILES {
		node = api.PathToNode(k, node)
	}
	fmt.Println(node)
	return node
}

func NewHomeModel() *HomeModel {
	ms := []tea.Model{
		NewTreeModel(getTreeNodes()),
	}
	h := &HomeModel{
		Models: ms,
		Tree:   ms[0].(*TreeModel),
	}
	return h
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
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		default:
			for _, m2 := range m.Models {
				_, cmd := m2.Update(msg)
				cmds = append(cmds, cmd)
			}
		}
	}
	return m, tea.Batch(cmds...)
}

func (m *HomeModel) View() string {
	var v string
	// v = lipgloss.JoinHorizontal(m.Tree.View())
	v = m.Tree.View()
	return homeStyle.Render(v)
}
