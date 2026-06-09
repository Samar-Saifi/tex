package main

import (
	"os/exec"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

func handleKey(msg tea.KeyMsg, m model) (model, tea.Cmd) {

	switch m.currentMode {
	case normal:
		return handleKeyNormal(msg, m)

	case search:
		return handleKeySearch(msg, m)
	}

	return m, nil
}

func OpenParent(m model) model {

	if m.currentDir == "/" {
		return m
	}

	m.currentDir = filepath.Dir(m.currentDir)
	m.cursor = 0
	m.LoadData()

	return m
}

func OpenSelected(m model) (model, tea.Cmd) {
	selected := m.data[m.cursor]

	if selected.isDir {
		m.currentDir = selected.path
		m.cursor = 0
		m.LoadData()
	} else {

	}

	return m, OpenFileWithDefaultApp(m)
}

func OpenFileWithDefaultApp(m model) tea.Cmd {

	if len(m.data) == 0 || m.cursor >= len(m.data) {
		return nil
	}

	selected := m.data[m.cursor]
	if selected.isDir {
		return nil
	}

	c := exec.Command("xdg-open", selected.path)

	return tea.ExecProcess(c, func(err error) tea.Msg {
		return nil
	})
}
