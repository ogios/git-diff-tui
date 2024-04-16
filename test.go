package main

import (
	"fmt"
	"os"

	"github.com/charmbracelet/lipgloss"
	"github.com/ogios/merge-repo/ui/comp"
)

var (
	copyColor                = lipgloss.Color("#00bd86")
	currentLineStyle         = lipgloss.NewStyle().Foreground(lipgloss.Color("#000000"))
	currentLineSelectedStyle = currentLineStyle.Copy().Background(copyColor)
	currentLineNormalStyle   = currentLineStyle.Copy().Background(lipgloss.Color("#ffffff"))
	selectedStyle            = lipgloss.NewStyle().Foreground(copyColor)
)

func main() {
	s := "一二三"
	a := s + selectedStyle.Render(s) + s
	a = s + currentLineNormalStyle.Render(a)
	a = currentLineSelectedStyle.Render(a)
	// fmt.Println(a, len(a))
	os.WriteFile("./test.log", []byte(a), 0766)

	ansi, shit := comp.GetANSIs(a)
	fmt.Println(ansi)
	fmt.Println(shit)
}
