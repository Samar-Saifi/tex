package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

func NormalView(m model) string {
	s := "Normal Mode\n\n"

	s += ViewFilesAndFolders(m)

	s += "\n\n[enter] open  [q] quit [s] search [backspace] go to parent"

	return s
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

	case keymap.search:
		m.currentMode = search

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

		m = OpenSelected(m)

	case keymap.back:
		m = OpenParent(m)
	}

	return m, nil
}
