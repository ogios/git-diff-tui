package utree

import (
	"log"
	"slices"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/mattn/go-runewidth"
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

type Tree struct {
	Root               *api.Node
	Lines              []*NodeLine
	ViewLineProcessors []LineProcessor
	Block              [2]int
	CurrentViewIndex   [2]int
	CurrentLine        int
}

type FileMsg struct {
	FileRelPath string
}

type CopyFileMsg struct {
	Files []string
}

var (
	copyColor                = lipgloss.Color("#00bd86")
	currentLineStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#000000"))
	currentLineSelectedStyle = currentLineStyle.Copy().Background(copyColor)
	currentLineNormalStyle   = currentLineStyle.Copy().Background(lipgloss.Color("#ffffff"))
	selectedStyle            = lipgloss.NewStyle().Foreground(copyColor)
)

func newTree(n *api.Node, block [2]int) *Tree {
	// block[1] -= 2
	lines := DrawNode(n, 0)
	log.Println("new tree", block)
	t := &Tree{
		Root:             n,
		Lines:            lines,
		Block:            block,
		CurrentViewIndex: [2]int{0, 0},
	}
	// t.updateSelection(0)
	t.ViewLineProcessors = []LineProcessor{
		func(input *strings.Builder, i int) {
			visibleLen := t.Block[0]
			input.Grow(visibleLen)
			start := t.CurrentViewIndex[0]
			end := start + visibleLen
			nl := t.Lines[i]
			inputLen := runewidth.StringWidth(nl.Text)
			if start < inputLen {
				s := runewidth.TruncateLeft(nl.Text, start, "")
				if end < inputLen {
					input.WriteString(runewidth.Truncate(s, end-start, ""))
				} else {
					input.WriteString(s)
					viewWidth := runewidth.StringWidth(s)
					input.WriteString(strings.Repeat(" ", visibleLen-viewWidth))
				}
			} else {
				input.WriteString(strings.Repeat(" ", visibleLen))
			}
		},
		func(input *strings.Builder, i int) {
			l := t.Lines[i]
			if i == t.CurrentLine {
				s := input.String()
				input.Reset()
				var rs string
				if l.Node.IsCopy() {
					rs = currentLineSelectedStyle.Render(s)
				} else {
					rs = currentLineNormalStyle.Render(s)
				}
				input.Grow(len(rs))
				input.WriteString(rs)
			} else if l.Node.IsCopy() {
				s := input.String()
				input.Reset()
				rs := selectedStyle.Render(s)
				input.Grow(len(rs))
				input.WriteString(rs)
			}
		},
	}
	return t
}

func (t *Tree) updateFileMsg() tea.Cmd {
	return func() tea.Msg {
		if len(t.Lines) == 0 {
			return nil
		}
		nl := t.Lines[t.CurrentLine]
		if nl.Node.Type == api.NODE_FILE {
			return FileMsg{FileRelPath: api.NodeToPath(nl.Node)}
		} else {
			return nil
		}
	}
}

func (t *Tree) copyFiles() tea.Cmd {
	fs := []string{}
	api.LoopFilesUnder(t.Root, func(n *api.Node) {
		if n.IsCopy() {
			fs = append(fs, api.NodeToPath(n))
		}
	})
	fs = slices.Clip(fs)

	return func() tea.Msg {
		if len(fs) == 0 {
			return nil
		}
		return CopyFileMsg{
			Files: fs,
		}
	}
}

func (t *Tree) Init() tea.Cmd {
	return t.updateFileMsg()
}

func (t *Tree) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	cmds := make([]tea.Cmd, 0)
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "j":
			cmds = append(cmds, t.nextLine(1))
		case "k":
			cmds = append(cmds, t.prevLine(1))
		case "h":
			t.prevCol(1)
		case "l":
			t.nextCol(1)
		case "H":
			t.prevCol(1)
		case "L":
			t.nextCol(1)
		case "ctrl+d":
			cmds = append(cmds, t.nextLine(10))
		case "ctrl+u":
			cmds = append(cmds, t.prevLine(10))

		case " ":
			n := t.Lines[t.CurrentLine].Node
			n.ToggleCopy()
		case "a":
			t.Root.ToggleCopy()
		case "c":
			cmds = append(cmds, t.copyFiles())
		}
	case tea.MouseMsg:
		mouse := tea.MouseEvent(msg)
		switch mouse.Button {
		case tea.MouseButtonWheelUp:
			if mouse.Ctrl {
				t.prevCol(1)
			} else {
				cmds = append(cmds, t.prevLine(1))
			}
		case tea.MouseButtonWheelDown:
			if mouse.Ctrl {
				t.nextCol(1)
			} else {
				cmds = append(cmds, t.nextLine(1))
			}
		case tea.MouseButtonWheelLeft:
			t.prevCol(1)
		case tea.MouseButtonWheelRight:
			t.nextCol(1)
		case tea.MouseButtonLeft:
			cmds = append(cmds, t.updateCurrentLine(t.CurrentViewIndex[1]+mouse.Y))

		}
	}
	return t, tea.Batch(cmds...)
}

func (t *Tree) nextLine(step int) tea.Cmd {
	return t.updateCurrentLine(min(t.CurrentLine+step, len(t.Lines)-1))
}

func (t *Tree) prevLine(step int) tea.Cmd {
	return t.updateCurrentLine(max(t.CurrentLine-step, 0))
}

func (t *Tree) nextCol(step int) {
	t.CurrentViewIndex[0] += step
}

func (t *Tree) prevCol(step int) {
	t.CurrentViewIndex[0] = max(t.CurrentViewIndex[0]-step, 0)
}

func (t *Tree) updateCurrentLine(i int) tea.Cmd {
	if t.CurrentLine != i {
		i = min(max(i, 0), len(t.Lines)-1)
		t.CurrentLine = i
		t.updateViewIndexY()
		return t.updateFileMsg()
	}
	return nil
}

func (t *Tree) updateViewIndexY() {
	moveViewThreshold := 3
	if t.CurrentViewIndex[1]+t.Block[1]-moveViewThreshold < t.CurrentLine {
		t.CurrentViewIndex[1] = max(t.CurrentLine-t.Block[1]+moveViewThreshold, 0)
	}
	if t.CurrentViewIndex[1]+moveViewThreshold > t.CurrentLine {
		t.CurrentViewIndex[1] = max(t.CurrentLine-moveViewThreshold, 0)
	}
}

func (t *Tree) View() string {
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
	return s.String()
}

const INDENT = "  "

// func DrawNode(n *api.Node, level int) string {
func DrawNode(n *api.Node, level int) []*NodeLine {
	ls := make([]*NodeLine, 0)
	var text strings.Builder
	if isRoot(n) {
		level = -1
	} else {
		text.Grow(len(INDENT)*level + len(n.Name) + 1)
		text.WriteString(strings.Repeat(INDENT, level))
		if n.Type == api.NODE_DIR {
			text.WriteString("▸")
		}
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
