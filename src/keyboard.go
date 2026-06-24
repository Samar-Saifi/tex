package main

import (
	"os"
	"os/exec"
	"path/filepath"

	tea "github.com/charmbracelet/bubbletea"
)

func handleKey(msg tea.KeyMsg, m model) (model, tea.Cmd) {

	switch m.currentMode {
	case normal:
		return handleKeyNormal(msg, m)

	case rename:
		return handleKeyRename(msg, m)

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

func OpenSelected(m model) (model, tea.Cmd) {
	selected := m.data[m.cursor]

	if selected.isDir {
		m.currentDir = selected.path
		m.cursor = 0
		m.LoadData()
	} else {

	}

	return m, OpenFileWithDefaultApp(m)
}

func OpenFileWithDefaultApp(m model) tea.Cmd {

	if len(m.data) == 0 || m.cursor >= len(m.data) {
		return nil
	}

	selected := m.data[m.cursor]
	if selected.isDir {
		return nil
	}

	var cmdName string
	var args []string

	if os.PathSeparator == '\\' {
		cmdName = "cmd"
		args = []string{"/c", "start", "", selected.path}
	} else {
		cmdName = "xdg-open"
		args = []string{selected.path}
	}

	c := exec.Command(cmdName, args...)

	return tea.ExecProcess(c, func(err error) tea.Msg {
		return nil
	})
}

func OpenTerminal(m model) tea.Cmd {

	var shell string
	if os.PathSeparator == '\\' {
		shell = os.Getenv("COMSPEC")
	} else {
		shell = os.Getenv("SHELL")
	}

	cmd := exec.Command(shell)

	cmd.Dir = m.currentDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return tea.ExecProcess(cmd,
		func(err error) tea.Msg {
			return nil
		})
}
