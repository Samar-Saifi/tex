package main

import (
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

const cols = 5

func main() {
	program := tea.NewProgram(initialModel())

	if _, err := program.Run(); err != nil {
		print("Error: " + err.Error())
		os.Exit(1)
	}
}
