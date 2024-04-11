package ui

import (
	"math"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ogios/merge-repo/ui/comp"
)

const (
	STATE_COMMIT = 0
	STATE_BRANCH = 1
)

type ChooseModel struct {
	*comp.NamedModel
	State int
}

var (
	modelStyle = lipgloss.NewStyle().
			Width(15).
			Height(1).
			Align(lipgloss.Center, lipgloss.Center).
			BorderStyle(lipgloss.HiddenBorder())
	focusedModelStyle = lipgloss.NewStyle().
				Width(15).
				Height(1).
				Align(lipgloss.Center, lipgloss.Center).
				BorderStyle(lipgloss.NormalBorder()).
				BorderForeground(lipgloss.Color("69"))
)

func NewChooseModel(parent tea.Model) *ChooseModel {
	return &ChooseModel{
		State:      STATE_COMMIT,
		NamedModel: comp.NewNamedModel("choose", parent),
	}
}

func (m *ChooseModel) Init() tea.Cmd {
	return func() tea.Cmd {
		return func() tea.Msg {
			m.Parent.(*comp.NamedModel).Call(m)
			return nil
		}
	}()
}

func (m *ChooseModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyLeft:
			fallthrough
		case tea.KeyRight:
			m.State = int(math.Abs(float64(m.State) - 1))
		}
	}
	return m, nil
}

func (m *ChooseModel) View() string {
	sort := []int{STATE_COMMIT, STATE_BRANCH}
	cs := map[int]string{
		STATE_COMMIT: "commit",
		STATE_BRANCH: "branch",
	}
	rs := make([]string, len(sort))
	for i, v := range sort {
		var r string
		raw := cs[v]
		if v == m.State {
			r = focusedModelStyle.Render(raw)
		} else {
			r = modelStyle.Render(raw)
		}
		rs[i] = r
	}

	s := lipgloss.JoinHorizontal(
		lipgloss.Top,
		rs...,
	)
	return s
}
