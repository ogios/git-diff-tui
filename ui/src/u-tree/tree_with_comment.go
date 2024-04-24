package utree

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/ogios/merge-repo/api"
	ucomment "github.com/ogios/merge-repo/ui/src/u-comment"
)

type TreeWithComment struct {
	tree    *Tree
	comment tea.Model
	current tea.Model
}

func newTreeWithComment(n *api.Node, block [2]int) *TreeWithComment {
	t := &TreeWithComment{
		tree:    newTree(n, block),
		comment: ucomment.NewCommentsModel(block),
	}
	t.current = t.tree
	return t
}

func (t *TreeWithComment) Init() tea.Cmd {
	return tea.Batch(t.tree.Init(), t.comment.Init())
}

func (t *TreeWithComment) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "i":
			if t.current != t.comment {
				t.current = t.comment
				return t, nil
			}
		case "esc":
			if t.current == t.comment {
				t.current = t.tree
				return t, nil
			}
		}
	}
	_, cmd := t.current.Update(msg)
	return t, cmd
}

func (t *TreeWithComment) View() string {
	return t.current.View()
}
