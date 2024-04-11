package ui

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ogios/merge-repo/ui/comp"
)

type HomeModelChildren struct {
	Choose *ChooseModel
}

type HomeModel struct {
	Text1 tea.Model
	Text2 tea.Model
	*comp.NamedModel
	Focusd int
}

var homeFocusedModelStyle = lipgloss.NewStyle().
	// Width().
	// Height(1).
	// Align(lipgloss.Bottom, lipgloss.Bottom).
	BorderStyle(lipgloss.NormalBorder()).
	BorderForeground(lipgloss.Color("69"))

func NewHomeModel() *HomeModel {
	h := &HomeModel{
		NamedModel: comp.NewNamedModel("home", nil),
	}
	h.Text1 = NewInputModel(h)
	h.Text2 = NewInputModel(h)
	ms := []tea.Model{
		NewChooseModel(h),
		h.Text1,
		h.Text2,
	}
	h.Models = ms
	h.Focusd = 0
	return h
}

func (m *HomeModel) Init() tea.Cmd {
	return nil
}

func (m *HomeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "tab":
			m.Focusd = (m.Focusd + 1) % len(m.Models)
		}
	}
	cmds := make([]tea.Cmd, 0)
	_, cmd := m.Models[m.Focusd].Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m *HomeModel) View() string {
	var v string
	for i, model := range m.Models {
		mv := model.View()
		if i == m.Focusd {
			mv = homeFocusedModelStyle.Render(mv)
		}
		v += mv + "\n"
	}
	return v
}

func (m *HomeModel) Call(caller tea.Model) {
}
