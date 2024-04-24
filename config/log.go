package config

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

var LoggerFile *os.File

func initLog() {
	f, err := tea.LogToFile("debug.log", "")
	if err != nil {
		panic(err)
	}
	LoggerFile = f
}

func exitLog() {
	LoggerFile.Close()
}
