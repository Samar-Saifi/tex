package main

import (
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
	popup
)

type model struct {
	currentDir     string
	data           []entry
	cursor         int
	currentMode    mode
	searchQuery    string
	startIndex     int
	terminalWidth  int
	terminalHeight int
	originalName   string
	errorMsg       string
}

func initialModel() model {
	dir, _ := os.Getwd()
	m := model{
		currentDir:  dir,
		cursor:      0,
		data:        nil,
		currentMode: normal,
		errorMsg:    "",
	}

	m.LoadData()

	return m
}

func (m *model) LoadData() {
	m.data = nil
	m.startIndex = 0

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

	case tea.WindowSizeMsg:
		m.terminalWidth = msg.Width
		m.terminalHeight = msg.Height
		return m, nil

	case error:
		m.errorMsg = msg.Error()

	case tea.KeyMsg:
		return handleKey(msg, m)
	}

	return m, nil
}

func (m model) View() string {
	return MainLayout(m)
}
