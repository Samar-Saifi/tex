package main

import (
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func handleKeyNormal(msg tea.KeyMsg, m model) (model, tea.Cmd) {

	key := msg.String()
	key = strings.ToLower(key)

	switch key {

	case keymap.quit:
		return m, tea.Quit

	case keymap.left:
		m = OpenParent(m)

	case keymap.search:
		m.currentMode = search
		m.searchQuery = ""

	case keymap.up:
		if m.cursor > 0 {
			m.cursor--

			if m.cursor < m.startIndex {
				m.startIndex = m.cursor
			}
		}

	case keymap.terminal:
		return m, OpenTerminal(m)

	case keymap.rename:
		m.currentMode = rename
		m.originalName = m.data[m.cursor].name
		m.data[m.cursor].name = ""

	case keymap.dlt:
		fileToBeDeleted := filepath.Join(m.currentDir, m.data[m.cursor].name)

		var err error

		if m.data[m.cursor].isDir {
			err = os.RemoveAll(fileToBeDeleted)
		} else {
			err = os.Remove(fileToBeDeleted)
		}

		if err != nil {
			m.errorMsg = err.Error()
			return m, nil
		}

		m.errorMsg = ""

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

	case keymap.right:
		return OpenSelected(m)

	case keymap.confirm:
		return OpenSelected(m)

	case keymap.back:
		m = OpenParent(m)
	}

	return m, nil
}
