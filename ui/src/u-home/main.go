package uhome

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ogios/merge-repo/config"
)

func NewHomeModel() tea.Model {
	resetStyle()
	if config.GlobalConfig.DiffSrc == "" {
		return newHome()
	} else {
		return newHomeDiff()
	}
}
