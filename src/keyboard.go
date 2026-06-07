package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func handleKey(msg tea.KeyMsg, m model) (model, tea.Cmd) {

	switch m.currentMode {
	case normal:
		return handleKeyNormal(msg, m)
	}

	return m, nil
}

func handleKeyNormal(msg tea.KeyMsg, m model) (model, tea.Cmd) {

	key := msg.String()
	key = strings.ToLower(key)

	switch key {

	case keymap.quit:
		return m, tea.Quit

	case keymap.left:
		if m.cursor > 0 {
			m.cursor--
		}

	case keymap.up:
		if m.cursor >= cols {
			m.cursor -= cols
		}

	case keymap.down:
		if m.cursor+cols < len(m.data) {
			m.cursor += cols
		}

	case keymap.right:
		if m.cursor < len(m.data)-1 {
			m.cursor++
		}

	case keymap.confirm:
		if len(m.data) == 0 {
			return m, nil
		}

		selected := m.data[m.cursor]

		if selected.isDir {
			m.currentDir = selected.path
			m.cursor = 0
			m.LoadData()
		}

	}

	return m, nil
}
