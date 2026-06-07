package main

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func SearchView(m model) string {
	s := "Search Mode\n\n"

	s += ViewFilesAndFolders(m)

	s += fmt.Sprintf("\n\nSearch: %s", m.searchQuery)

	s += "\n\n[enter] open  [esc] cancel"

	return s
}

func handleKeySearch(msg tea.KeyMsg, m model) (model, tea.Cmd) {

	key := msg.String()
	key = strings.ToLower(key)

	switch key {

	case keymap.cancel:
		m.currentMode = normal
		m.searchQuery = ""
		return m, nil

	case keymap.back:
		if len(m.searchQuery) > 0 {
			m.searchQuery = m.searchQuery[:len(m.searchQuery)-1]
		}

	case keymap.confirm:
		m.currentMode = normal
		m = OpenSelected(m)

	default:
		if len(msg.String()) == 1 {
			m.searchQuery += msg.String()
		}
	}

	query := strings.ToLower(m.searchQuery)

	for i, e := range m.data {
		if strings.Contains(strings.ToLower(e.name), query) {
			m.cursor = i
			break
		}
	}

	return m, nil
}
