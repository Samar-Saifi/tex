package main

import (
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

func OpenSelected(m model) model {
	selected := m.data[m.cursor]

	if selected.isDir {
		m.currentDir = selected.path
		m.cursor = 0
		m.LoadData()
	}

	return m
}
