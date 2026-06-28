package main

import (
	"bufio"
	"errors"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type DesktopEntry struct {
	Name        string
	DesktopFile string // e.g. "vlc.desktop"
	Exec        string // e.g. "vlc %U"
	MimeTypes   []string
}

func OpenWith(m *model) {
	selected := m.data[m.cursor]

	entries, err := FindDesktopEntries(selected.path)
	if err != nil {
		m.errorMsg = err.Error()
		return
	}

	if len(entries) == 0 {
		m.errorMsg = "No compatible applications found."
		return
	}

	options := make([]string, len(entries))
	for i, e := range entries {
		options[i] = e.Name
	}

	ShowPopup(
		m,
		"Open with",
		PopupVertical,
		options,
		func(idx int) (model, tea.Cmd) {
			return *m, LaunchDesktopEntry(entries[idx], selected.path)
		},
	)
}

func LaunchDesktopEntry(entry DesktopEntry, file string) tea.Cmd {
	args, err := BuildExec(entry.Exec, file)
	if err != nil {
		return func() tea.Msg { return err }
	}

	cmd := exec.Command(args[0], args[1:]...)

	return tea.ExecProcess(cmd, func(err error) tea.Msg {
		return err
	})
}

func FindDesktopEntries(filePath string) ([]DesktopEntry, error) {
	mimeType, err := detectMimeType(filePath)
	if err != nil {
		return nil, err
	}

	dirs := []string{
		"/usr/share/applications",
		"/usr/local/share/applications",
	}

	if home, err := os.UserHomeDir(); err == nil {
		dirs = append([]string{
			filepath.Join(home, ".local/share/applications"),
		}, dirs...)
	}

	var entries []DesktopEntry

	for _, dir := range dirs {
		filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
			if err != nil || d.IsDir() {
				return nil
			}

			if filepath.Ext(path) != ".desktop" {
				return nil
			}

			entry, err := parseDesktopEntry(path)
			if err != nil {
				return nil
			}

			if supportsMime(entry.MimeTypes, mimeType) {
				entries = append(entries, entry)
			}

			return nil
		})
	}

	return entries, nil
}

func parseDesktopEntry(path string) (DesktopEntry, error) {
	file, err := os.Open(path)
	if err != nil {
		return DesktopEntry{}, err
	}
	defer file.Close()

	var entry DesktopEntry

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		switch {
		case strings.HasPrefix(line, "Name="):
			entry.Name = strings.TrimPrefix(line, "Name=")

		case strings.HasPrefix(line, "Exec="):
			entry.Exec = strings.TrimPrefix(line, "Exec=")

		case strings.HasPrefix(line, "MimeType="):
			entry.MimeTypes = strings.Split(
				strings.TrimSuffix(strings.TrimPrefix(line, "MimeType="), ";"),
				";",
			)

		case line == "NoDisplay=true":
			return DesktopEntry{}, errors.New("hidden")

		case line == "Hidden=true":
			return DesktopEntry{}, errors.New("hidden")
		}
	}

	if err := scanner.Err(); err != nil {
		return DesktopEntry{}, err
	}

	entry.DesktopFile = filepath.Base(path)

	if entry.Name == "" || entry.Exec == "" {
		return DesktopEntry{}, errors.New("invalid")
	}

	return entry, nil
}

func detectMimeType(path string) (string, error) {
	out, err := exec.Command("xdg-mime", "query", "filetype", path).Output()
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(out)), nil
}

func supportsMime(list []string, mime string) bool {
	for _, m := range list {
		if m == mime {
			return true
		}
	}
	return false
}

func BuildExec(execLine, file string) ([]string, error) {
	args := strings.Fields(execLine)

	var result []string
	for _, arg := range args {
		switch arg {
		case "%f", "%u", "%F", "%U":
			result = append(result, file)

		case "%i", "%c", "%k":
			// Ignore these placeholders.

		default:
			result = append(result, arg)
		}
	}

	if len(result) == 0 {
		return nil, errors.New("invalid Exec line")
	}

	return result, nil
}
