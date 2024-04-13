package api

import "path"

const (
	NODE_DIR  = 0
	NODE_FILE = 1
)

type Node struct {
	Parent   *Node
	Name     string
	Children []*Node
	Type     int
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
	for _, v := range parent.Children {
		if v.Name == p {
			next = v
			break
		}
	}
	if next == nil {
		if len(path) == 1 {
			next = newFile()
		} else {
			next = newDir()
		}
	}
	next.Name = p
	next.Parent = parent
	if len(path) > 1 {
		CreateNode(path[:len(path)-1], next)
	}
	parent.Children = append(parent.Children, next)
}

func newDir() *Node {
	return &Node{
		Type: NODE_DIR,
	}
}

func newFile() *Node {
	return &Node{
		Type: NODE_FILE,
	}
}
