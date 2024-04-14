package api

import (
	"path"
	"strings"
)

const (
	NODE_DIR  = 0
	NODE_FILE = 1
)

type Node struct {
	Parent     *Node
	CopiesNode map[string]struct{}
	Name       string
	Children   []*Node
	Type       int
}

func PathToNode(p string, root *Node) *Node {
	if root == nil {
		root = newDir()
	}
	bs := []string{}
	for {
		next, f := path.Split(p)
		bs = append(bs, f)
		if next == "" {
			break
		}
		p = path.Dir(next)
	}
	CreateNode(bs, root)
	return root
}

func CreateNode(path []string, parent *Node) {
	p := path[len(path)-1]
	var next *Node
	var currentType int
	if len(path) == 1 {
		currentType = NODE_FILE
	} else {
		currentType = NODE_DIR
	}
	for _, v := range parent.Children {
		if v.Name == p && v.Type == currentType {
			next = v
			break
		}
	}
	if next == nil {
		if currentType == NODE_FILE {
			next = newFile()
		} else {
			next = newDir()
		}
		parent.Children = append(parent.Children, next)
	}
	next.Name = p
	next.Parent = parent
	if len(path) > 1 {
		CreateNode(path[:len(path)-1], next)
	}
}

func newDir() *Node {
	return &Node{
		Type:       NODE_DIR,
		CopiesNode: map[string]struct{}{},
	}
}

func newFile() *Node {
	return &Node{
		Type: NODE_FILE,
	}
}

func NodeToPath(n *Node) string {
	names := []string{}
	count := 0
	for n.Parent != nil {
		names = append(names, n.Name)
		count += len(n.Name)
		n = n.Parent
	}
	var b strings.Builder
	b.Grow(count + len(names))
	for i := len(names) - 1; i >= 0; i-- {
		b.WriteString(names[i])
		if i > 0 {
			b.WriteString("/")
		}
	}
	return b.String()
}

func LoopFilesUnder(root *Node, handler func(n *Node)) {
	if root.Type == NODE_FILE {
		handler(root)
	} else {
		for _, v := range root.Children {
			LoopFilesUnder(v, handler)
		}
	}
}

func (t *Node) addCopy(n string) {
	t.CopiesNode[n] = struct{}{}
	if len(t.CopiesNode) == len(t.Children) {
		if t.Parent != nil {
			t.Parent.AddCopy(t.Name)
		}
	}
}

func (t *Node) AddCopy(n string) {
	if _, ok := t.CopiesNode[n]; !ok {
		t.addCopy(n)
	}
}

func (t *Node) RmCopy(n string) {
	delete(t.CopiesNode, n)
	if t.IsCopy() {
		if t.Parent != nil {
			t.Parent.RmCopy(t.Name)
		}
	}
}

func (t *Node) ToggleCopy() {
	if t.Type == NODE_FILE {
		if t.IsCopy() {
			t.Parent.RmCopy(t.Name)
		} else {
			t.Parent.AddCopy(t.Name)
		}
	} else {
		var handler func(n *Node)
		if len(t.Children) == len(t.CopiesNode) {
			handler = func(n *Node) {
				n.Parent.RmCopy(n.Name)
			}
		} else {
			handler = func(n *Node) {
				n.Parent.AddCopy(n.Name)
			}
		}
		LoopFilesUnder(t, handler)
	}
}

func (t *Node) IsCopy() bool {
	if t.Parent != nil {
		_, ok := t.Parent.CopiesNode[t.Name]
		return ok
	}
	return false
}
