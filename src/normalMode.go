package main

import (
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func handleKeyNormal(msg tea.KeyMsg, m model) (model, tea.Cmd) {

	key := msg.String()
	key = strings.ToLower(key)

	switch key {
	case keymap.left, keymap.up, keymap.down, keymap.left, keymap.right, keymap.confirm:
		return handleNavigation(m, key)

	case keymap.search:
		m.currentMode = search
		m.searchQuery = ""

	case keymap.terminal:
		return m, OpenTerminal(m)

	case keymap.rename:
		m.currentMode = rename
		m.originalName = m.data[m.cursor].name
		m.data[m.cursor].name = ""

	case keymap.dlt:
		fileToBeDeleted := filepath.Join(m.currentDir, m.data[m.cursor].name)
		callback := func(selectedIndex int) (model, tea.Cmd) {
			if selectedIndex == 0 {
				m = deleteFile(m, fileToBeDeleted)
				m.LoadData()
			}
			return m, nil
		}

		ShowPopup(&m, "Are you sure you want to delete?", []string{"Yes, Delete", "No, Cancel"}, callback)
	default:
		if len(key) == 1 {
			m.currentMode = search
			m.searchQuery = key
		}
	}

	return m, nil
}

func handleNavigation(m model, key string) (model, tea.Cmd) {
	switch key {
	case keymap.up:
		if m.cursor > 0 {
			m.cursor--

			if m.cursor < m.startIndex {
				m.startIndex = m.cursor
			}
		}

	case keymap.down:
		if m.cursor < len(m.data)-1 {
			m.cursor++

			maxVisibleLines := m.terminalHeight - 5
			if m.currentMode == search {
				maxVisibleLines = m.terminalHeight - 8
			}

			if m.cursor >= m.startIndex+maxVisibleLines {
				m.startIndex = m.cursor - maxVisibleLines + 1
			}
		}

	case keymap.left:
		m = OpenParent(m)

	case keymap.right:
		if m.data[m.cursor].name != ".." {
			return OpenSelected(m)
		}

	case keymap.confirm:
		return OpenSelected(m)

	case keymap.back:
		m = OpenParent(m)

	}

	return m, nil
}
