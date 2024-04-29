package utree

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ogios/merge-repo/api"
)

func NewTreeModel(n *api.Node, block [2]int) tea.Model {
	return newTreeWithComment(n, block)
}
