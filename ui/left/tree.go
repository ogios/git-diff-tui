package left

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/ogios/merge-repo/api"
	"github.com/ogios/merge-repo/ui/comp"
)

type NodeLine struct {
	Node *api.Node
	Text string
}
type LineProcessor func(input *strings.Builder, index int)

type TreeModel struct {
	Root               *api.Node
	Lines              []NodeLine
	CurrentLine        int
	CurrentViewLine    int
	Block              [2]int
	CurrentViewIndex   [2]int
	ViewLineProcessors []LineProcessor
}

var (
	selectedNodeLineStyle = lipgloss.NewStyle().Background(lipgloss.Color("#ffffff")).Foreground(lipgloss.Color("#000000"))
	treeStyle             lipgloss.Style
)

func NewTreeModel(n *api.Node) tea.Model {
	width := int(float64(comp.GlobalUIData.MaxWidth) * 0.1)
	height := comp.GlobalUIData.MaxHeight - 4
	treeStyle = lipgloss.NewStyle().BorderStyle(lipgloss.NormalBorder()).BorderForeground(lipgloss.Color("#ffad00")).Width(width).Height(height)
	t := &TreeModel{
		Root:             n,
		Lines:            DrawNode(n, 0),
		CurrentLine:      0,
		CurrentViewLine:  0,
		Block:            [2]int{width, height},
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
			t.CurrentLine = ((t.CurrentLine + 1) + len(t.Lines)) % len(t.Lines)
		case "k":
			t.CurrentLine = ((t.CurrentLine - 1) + len(t.Lines)) % len(t.Lines)
		case "h":
			if t.CurrentViewIndex[0] > 0 {
				t.CurrentViewIndex[0]--
			}
		case "l":
			t.CurrentViewIndex[0]++
		}
	}
	return t, tea.Batch(cmds...)
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
		s.WriteString("\n")
	}
	return treeStyle.Render(s.String())
}

const INDENT = "  "

// func DrawNode(n *api.Node, level int) string {
func DrawNode(n *api.Node, level int) []NodeLine {
	ls := make([]NodeLine, 0)
	var text strings.Builder
	if isRoot(n) {
		level = -1
	} else {
		text.Grow(len(INDENT)*level + len(n.Name))
		text.WriteString(strings.Repeat(INDENT, level))
		text.WriteString(n.Name)
		ls = append(ls, NodeLine{
			Node: n,
			Text: text.String(),
		})
	}
	for _, child := range n.Children {
		ls = append(ls, DrawNode(child, level+1)...)
	}
	return ls
}

func isRoot(n *api.Node) bool {
	return n.Parent == nil
}
