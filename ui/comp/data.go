package comp

import (
	"fmt"
	"os"

	"golang.org/x/term"
)

type UIData struct {
	MaxWidth  int
	MaxHeight int
}

var GlobalUIData UIData

func init() {
	w, h, err := term.GetSize(int(os.Stdout.Fd()))
	if err != nil {
		panic(err)
	}
	fmt.Println(w, h)
	GlobalUIData = UIData{MaxWidth: w, MaxHeight: h}
}
