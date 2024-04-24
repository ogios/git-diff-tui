package comp

import "github.com/charmbracelet/lipgloss"

var ErrorMsgStyle = lipgloss.NewStyle().
	Bold(true).
	Inline(true).
	Background(lipgloss.Color("#ff2f5a")).
	Foreground(lipgloss.Color("#000000"))
