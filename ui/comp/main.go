package comp

import tea "github.com/charmbracelet/bubbletea"

type Ctx struct {
	Focusd tea.Model
}

type NamedModel struct {
	Parent tea.Model
	Ctx    Ctx
	Name   string
	Models []tea.Model
}

type NamedModelInterface interface {
	Call(caller tea.Model)
}

func (m *NamedModel) Init() tea.Cmd {
	return nil
}

func (m *NamedModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	return m, nil
}

func (m *NamedModel) View() string {
	return ""
}

func (m *NamedModel) Call(caller any) {
}

func NewNamedModel(name string, parent tea.Model) *NamedModel {
	n := &NamedModel{
		Name:   name,
		Ctx:    Ctx{Focusd: nil},
		Models: make([]tea.Model, 0),
	}
	if parent != nil {
		n.Parent = parent
	}
	return n
}
