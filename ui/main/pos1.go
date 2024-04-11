package ui

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ogios/merge-repo/ui/comp"
)

type InputModel struct {
	*comp.NamedModel
	M textinput.Model
}

func NewInputModel(parent tea.Model) *InputModel {
	return &InputModel{
		M:          textinput.New(),
		NamedModel: comp.NewNamedModel("choose", parent),
	}
}

func (m *InputModel) Init() tea.Cmd {
	return func() tea.Cmd {
		return func() tea.Msg {
			m.Parent.(*comp.NamedModel).Call(m)
			return nil
		}
	}()
}

func (m *InputModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := []tea.Cmd{}
	if m.M.Focused() {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.Type {
			case tea.KeyEsc:
				m.M.Blur()
				return m, nil
			default:
				tm, cmd := m.M.Update(msg)
				m.M = tm
				cmds = append(cmds, cmd)
			}
		}
	} else {
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "i":
				cmds = append(cmds, m.M.Focus())
				return m, tea.Batch(cmds...)
			}
		}
	}
	return m, tea.Batch(cmds...)
}

func (m *InputModel) View() string {
	if m.M.Focused() {
		return homeFocusedModelStyle.Render(m.M.View())
	}
	return m.M.View()
}
