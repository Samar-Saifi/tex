package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func handleKeyRename(msg tea.KeyMsg, m model) (model, tea.Cmd) {

	key := msg.String()
	lowerKey := strings.ToLower(key)

	switch lowerKey {

	case keymap.cancel:
		m.currentMode = normal
		m.data[m.cursor].name = m.originalName
		return m, nil

	case keymap.back:
		r := []rune(m.data[m.cursor].name)
		if len(r) > 0 {
			m.data[m.cursor].name = string(r[:len(r)-1])
		}

	case keymap.confirm:
		oldName := filepath.Join(m.currentDir, m.originalName)
		newName := filepath.Join(m.currentDir, m.data[m.cursor].name)

		if err := os.Rename(oldName, newName); err != nil {
			m.errorMsg = fmt.Sprintf(
				"Failed to rename %q to %q: %v",
				m.originalName,
				m.data[m.cursor].name,
				err,
			)
			return m, nil
		}

		m.errorMsg = ""
		m.currentMode = normal

	default:
		if len(msg.String()) == 1 {
			m.data[m.cursor].name += key
		}
	}

	return m, nil
}
