package uhome

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type childModel struct {
	m     tea.Model
	style lipgloss.Style
	block [2]int
}

func newChild(block [2]int) *childModel {
	return &childModel{
		block: block,
		style: lipgloss.NewStyle().Width(block[0]).Height(block[1]),
	}
}
