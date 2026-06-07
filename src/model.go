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
	searchQuery string
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
			name:     "..",
			path:     filepath.Dir(m.currentDir),
			isDir:    true,
			isParent: true,
		})
	}

	var dirs []entry
	var regularFiles []entry

	for _, f := range files {
		e := entry{
			name:  f.Name(),
			path:  filepath.Join(m.currentDir, f.Name()),
			isDir: f.IsDir(),
		}

		if e.isDir == true {
			dirs = append(dirs, e)
		} else {
			regularFiles = append(regularFiles, e)
		}
	}

	m.data = append(m.data, dirs...)
	m.data = append(m.data, regularFiles...)
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
	return ViewSwitch(m)
}

func ViewSwitch(m model) string {
	s := ""

	switch m.currentMode {

	case normal:
		s = NormalView(m)

	case search:
		s = SearchView(m)

	}

	return s
}

func ViewFilesAndFolders(m model) string {

	s := ""

	if len(m.data) == 0 {
		return s + "No files found. \n\n [q] quit"
	}

	for i, e := range m.data {

		icon := "🗀"

		if i%cols == 0 && i != 0 {
			s += "\n\n"
		}

		cursorCharacter := " "
		if i == m.cursor {
			cursorCharacter = ">"
		}

		name := e.name

		if len(name) > 15 {
			name = name[:12] + "..."
		}

		if e.isDir && !e.isParent {
			name += "/"
		}

		if !e.isDir {
			icon = "🗎"
		}

		s += fmt.Sprintf("%s %s %-20s", cursorCharacter, icon, name)
	}

	return s
}
