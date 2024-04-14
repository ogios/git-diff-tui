package config

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

func CraeteLogger() *os.File {
	loggerFile, err := tea.LogToFile("debug.log", "")
	if err != nil {
		panic(err)
	}
	return loggerFile
}
