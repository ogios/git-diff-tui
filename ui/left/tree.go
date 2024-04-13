package left

import (
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ogios/merge-repo/api"
)

type NodeLine struct {
	Node *api.Node
	Text string
}
type SelectionNodeLine struct {
	Line      *NodeLine
	LineIndex int
}
type LineProcessor func(input *strings.Builder, index int)

type TreeModel struct {
	Root               *api.Node
	Lines              []*NodeLine
	selections         []*SelectionNodeLine
	ViewLineProcessors []LineProcessor
	Block              [2]int
	CurrentViewIndex   [2]int
	selectIndex        int
	CurrentLine        int
	CurrentViewLine    int
}

var (
	selectedNodeLineStyle = lipgloss.NewStyle().Background(lipgloss.Color("#ffffff")).Foreground(lipgloss.Color("#000000"))
	treeStyle             lipgloss.Style
)

func NewTreeModel(n *api.Node, block [2]int) tea.Model {
	block[1] -= 2
	treeStyle = lipgloss.NewStyle().
		BorderStyle(lipgloss.NormalBorder()).
		BorderForeground(lipgloss.Color("#ffad00")).
		Width(block[0]).
		Height(block[1])
	lines := DrawNode(n, 0)
	selections := make([]*SelectionNodeLine, 0)
	for i, nl := range lines {
		if nl.Node.Type == api.NODE_FILE {
			selections = append(selections, &SelectionNodeLine{
				LineIndex: i,
				Line:      nl,
			})
		}
	}
	t := &TreeModel{
		Root:             n,
		Lines:            lines,
		selections:       slices.Clip(selections),
		selectIndex:      0,
		CurrentLine:      0,
		CurrentViewLine:  0,
		Block:            block,
		CurrentViewIndex: [2]int{0, 0},
	}
	t.ViewLineProcessors = []LineProcessor{
		func(input *strings.Builder, i int) {
			visibleLen := t.Block[0]
			input.Grow(visibleLen)
			start := t.CurrentViewIndex[0]
			end := start + visibleLen
			nl := t.Lines[i]
			inputLen := len(nl.Text)
			if start < inputLen {
				if end < inputLen {
					input.WriteString(nl.Text[start:end])
					// try remove big chars
				} else {
					input.WriteString(nl.Text[start:])
					input.WriteString(strings.Repeat(" ", end-inputLen))
				}
			} else {
				input.WriteString(strings.Repeat(" ", visibleLen))
			}
		},
		func(input *strings.Builder, i int) {
			if i == t.CurrentLine {
				s := input.String()
				input.Reset()
				rs := selectedNodeLineStyle.Render(s)
				input.Grow(len(rs))
				input.WriteString(rs)
			}
		},
	}
	return t
}

func (t *TreeModel) Init() tea.Cmd {
	return nil
}

func (t *TreeModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j":
			// t.CurrentLine = ((t.CurrentLine + 1) + len(t.Lines)) % len(t.Lines)
			t.nextSelection(1)
		case "k":
			// t.CurrentLine = ((t.CurrentLine - 1) + len(t.Lines)) % len(t.Lines)
			t.prevSelection(1)
		case "h":
			t.CurrentViewIndex[0] = max(t.CurrentViewIndex[0]-1, 0)
		case "l":
			t.CurrentViewIndex[0]++
		case "ctrl+d":
			t.nextSelection(5)
		case "ctrl+u":
			t.prevSelection(5)
		}
	}
	return t, tea.Batch(cmds...)
}

func (t *TreeModel) nextSelection(step int) {
	t.updateSelection(min(t.selectIndex+step, len(t.selections)-1))
}

func (t *TreeModel) prevSelection(step int) {
	t.updateSelection(max(t.selectIndex-step, 0))
}

func (t *TreeModel) updateSelection(i int) {
	t.selectIndex = i
	t.updateCurrentLine(t.selections[i].LineIndex)
}

func (t *TreeModel) nextLine(step int) {
	t.updateCurrentLine(min(t.CurrentLine+step, len(t.Lines)-1))
}

func (t *TreeModel) prevLine(step int) {
	t.updateCurrentLine(max(t.CurrentLine-step, 0))
}

func (t *TreeModel) updateCurrentLine(i int) {
	t.CurrentLine = i
	t.updateViewIndexY()
}

func (t *TreeModel) updateViewIndexY() {
	moveViewThreshold := 3
	if t.CurrentViewIndex[1]+t.Block[1]-moveViewThreshold < t.CurrentLine {
		t.CurrentViewIndex[1] = max(t.CurrentLine-t.Block[1]+moveViewThreshold, 0)
	}
	if t.CurrentViewIndex[1]+moveViewThreshold > t.CurrentLine {
		t.CurrentViewIndex[1] = max(t.CurrentLine-moveViewThreshold, 0)
	}
}

func (t *TreeModel) View() string {
	// var s string
	var s strings.Builder
	startLine := t.CurrentViewIndex[1]
	endLine := startLine + t.Block[1]
	visibleLines := api.SliceFrom(t.Lines, startLine, endLine)
	for i := range visibleLines {
		input := func() string {
			inputBuilder := new(strings.Builder)
			for _, processor := range t.ViewLineProcessors {
				processor(inputBuilder, i+startLine)
			}
			return inputBuilder.String()
		}()
		s.Grow(len(input) + 1)
		s.WriteString(input)
		if i < len(visibleLines)-1 {
			s.WriteString("\n")
		}
	}
	return treeStyle.Render(s.String())
}

const INDENT = "  "

// func DrawNode(n *api.Node, level int) string {
func DrawNode(n *api.Node, level int) []*NodeLine {
	ls := make([]*NodeLine, 0)
	var text strings.Builder
	if isRoot(n) {
		level = -1
	} else {
		text.Grow(len(INDENT)*level + len(n.Name))
		text.WriteString(strings.Repeat(INDENT, level))
		text.WriteString(n.Name)
		ls = append(ls, &NodeLine{
			Node: n,
			Text: text.String(),
		})
	}
	for _, child := range n.Children {
		ls = append(ls, DrawNode(child, level+1)...)
	}
	return slices.Clip(ls)
}

func isRoot(n *api.Node) bool {
	return n.Parent == nil
}
