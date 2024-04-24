package utree

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ogios/merge-repo/api"
	"github.com/ogios/merge-repo/config"
)

func NewTreeModel(n *api.Node, block [2]int) tea.Model {
	if config.GlobalConfig.DiffSrc == "" {
		return newTree(n, block)
	} else {
		return newTreeWithComment(n, block)
	}
}
