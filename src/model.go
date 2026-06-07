package main

import (
	"fmt"
	"os"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

type entry struct {
	name     string
	path     string
	isDir    bool
	isParent bool
}

type mode int

const (
	normal mode = iota
	search
	rename
)

type model struct {
	currentDir  string
	data        []entry
	cursor      int
	currentMode mode
}

func initialModel() model {
	dir, _ := os.Getwd()
	m := model{
		currentDir:  dir,
		cursor:      0,
		data:        nil,
		currentMode: normal,
	}

	m.LoadData()

	return m
}

func (m *model) LoadData() {
	m.data = nil

	files, err := os.ReadDir(m.currentDir)
	if err != nil {
		return
	}

	if m.currentDir != "/" {
		m.data = append(m.data, entry{
			name:     "...",
			path:     filepath.Dir(m.currentDir),
			isDir:    true,
			isParent: true,
		})
	}

	for _, f := range files {
		m.data = append(m.data, entry{
			name:  f.Name(),
			path:  filepath.Join(m.currentDir, f.Name()),
			isDir: f.IsDir(),
		})
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {

	case tea.KeyMsg:
		return handleKey(msg, m)
	}

	return m, nil
}

func (m model) View() string {
	s := fmt.Sprintf("Current Mode: %d\n", m.currentMode)

	if len(m.data) == 0 {
		return s + "No files found. \n\n [q] quit"
	}

	for i, e := range m.data {
		if i%cols == 0 && i != 0 {
			s += "\n"
		}

		cursorCharacter := " "
		if i == m.cursor {
			cursorCharacter = ">"
		}

		name := e.name

		if e.isDir && !e.isParent {
			name += "/"
		}

		s += fmt.Sprintf("%s %-20s", cursorCharacter, name)
	}

	s += "\n\n[enter] open  [q] quit"

	return s
}
