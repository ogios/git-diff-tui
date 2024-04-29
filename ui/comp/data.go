package comp

import (
	"fmt"
	"os"

	process "github.com/ogios/ansisgr-process"
	"github.com/ogios/cropviewport"
	"github.com/ogios/merge-repo/api"
	"github.com/ogios/merge-repo/data"
	"golang.org/x/term"
)

type UIData struct {
	MaxWidth  int
	MaxHeight int
}

var (
	GlobalUIData UIData
	TREE_NODE    *api.Node
)

func init() {
	// ui data
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		panic(err)
	}
	fmt.Println(w, h)
	GlobalUIData = UIData{MaxWidth: w, MaxHeight: h}

	// node
	TREE_NODE = getTreeNodes()
}

func getTreeNodes() *api.Node {
	var node *api.Node = nil
	for k := range data.DIFF_FILES {
		node = api.PathToNode(k, node)
	}
	fmt.Println(node)
	return node
}

type ContentData struct {
	Table *process.ANSITableList
	Lines []*cropviewport.SubLine
}
