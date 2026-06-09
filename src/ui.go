package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	NormalHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#7D56F4")).
				Padding(0, 1)

	SearchHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#FF5F00")).
				Padding(0, 1)

	ListStyle = lipgloss.NewStyle().
			PaddingRight(1)

	CursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00FF87")).
			Bold(true)

	DirStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#00ADFF")).
			Bold(true)

	FileStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#EEEEEE"))

	PreviewStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(lipgloss.Color("#3A3A3A")).
			PaddingLeft(2).
			Foreground(lipgloss.Color("#A0A0A0"))

	SearchInputBoxStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color("#FF5F00")).
				Padding(0, 1).
				MarginTop(1)

	SearchPromptStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FF5F00")).
				Bold(true)

	SearchValueStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color("#3A3A3A"))

	MetadataStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#8A8A8A")).
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(lipgloss.Color("#262626")).
			PaddingTop(1)

	FooterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#585858")).
			PaddingTop(0)
)

func MainLayout(m model) string {
	width := m.terminalWidth
	if width == 0 {
		width = 80
	}
	height := m.terminalHeight
	if height == 0 {
		height = 24
	}

	neededChromeHeight := 5
	if m.currentMode == search {
		neededChromeHeight = 8
	}

	contentHeight := height - neededChromeHeight
	if contentHeight < 3 {
		contentHeight = 3
	}

	leftWidth := int(float64(width) * 0.4)
	rightWidth := width - leftWidth - 1

	var header string
	if m.currentMode == search {
		header = SearchHeaderStyle.Width(width).Render(" MODE: SEARCHING目录")
	} else {
		header = NormalHeaderStyle.Width(width).Render(" TEX")
	}

	leftContent := RenderVerticalFileList(m, leftWidth, contentHeight)
	rightContent := RenderFilePreview(m, rightWidth, contentHeight)

	mainContent := lipgloss.JoinHorizontal(
		lipgloss.Top,
		ListStyle.Width(leftWidth).Height(contentHeight).Render(leftContent),
		PreviewStyle.Width(rightWidth).Height(contentHeight).Render(rightContent),
	)

	var searchField string
	if m.currentMode == search {
		displayText := m.searchQuery
		displayText += "█"

		inputContent := fmt.Sprintf("%s %s", SearchPromptStyle.Render("🔍 Search target:"), SearchValueStyle.Render(displayText))
		searchField = SearchInputBoxStyle.Width(width - 2).Render(inputContent)
	}

	metadataStr := " No File Selected"
	if len(m.data) > 0 && m.cursor < len(m.data) {
		selected := m.data[m.cursor]
		info, err := os.Stat(selected.path)
		if err == nil {
			sizeStr := fmt.Sprintf("%d B", info.Size())
			if info.IsDir() {
				sizeStr = "DIR"
			} else if info.Size() > 1024*1024 {
				sizeStr = fmt.Sprintf("%.2f MB", float64(info.Size())/(1024*1024))
			} else if info.Size() > 1024 {
				sizeStr = fmt.Sprintf("%.2f KB", float64(info.Size())/1024)
			}
			metadataStr = fmt.Sprintf(" 📄 %s  •  💾 %s  •  🕒 %s  •  🔒 %s",
				info.Name(), sizeStr, info.ModTime().Format("2006-01-02 15:04"), info.Mode().String())
		}
	}
	metadata := MetadataStyle.Width(width).Render(metadataStr)

	var controlKeys string
	if m.currentMode == search {
		controlKeys = " [Type characters to search]  •  [enter] confirm/open focus  •  [esc] cancel search"
	} else {
		controlKeys = " [↑/↓] navigate  •  [enter/→] open  •  [←/backspace] parent  •  [s] search  •  [q] quit"
	}
	footer := FooterStyle.Width(width).Render(controlKeys)

	if m.currentMode == search {
		return lipgloss.JoinVertical(lipgloss.Left, header, mainContent, searchField, metadata, footer)
	}
	return lipgloss.JoinVertical(lipgloss.Left, header, mainContent, metadata, footer)
}

func RenderVerticalFileList(m model, width, height int) string {
	if len(m.data) == 0 {
		return " No files found."
	}

	var lines []string
	for i := m.startIndex; i < len(m.data); i++ {
		cursorChar := "  "
		if i == m.cursor {
			cursorChar = "> "
		}

		e := m.data[i]

		icon := "🗎 "
		name := e.name
		if e.isDir {
			icon = "📂 "
			if !e.isParent {
				name += "/"
			}
		}

		maxNameLen := width - 6
		if maxNameLen > 0 && len(name) > maxNameLen {
			name = name[:maxNameLen-3] + "..."
		}

		lineStr := fmt.Sprintf("%s%s%s", cursorChar, icon, name)

		if i == m.cursor {
			lines = append(lines, CursorStyle.Render(lineStr))
		} else if e.isDir {
			lines = append(lines, DirStyle.Render(lineStr))
		} else {
			lines = append(lines, FileStyle.Render(lineStr))
		}

		if len(lines) >= height {
			break
		}
	}

	return lipgloss.JoinVertical(lipgloss.Left, lines...)
}

func RenderFilePreview(m model, width, height int) string {
	if len(m.data) == 0 || m.cursor >= len(m.data) {
		return "No preview available"
	}

	selected := m.data[m.cursor]

	if selected.isDir {
		files, err := os.ReadDir(selected.path)
		if err != nil {
			return fmt.Sprintf("Error reading directory:\n%v", err)
		}

		var lines []string
		lines = append(lines, fmt.Sprintf("Directory contents (%d items):", len(files)))
		maxFileLines := height - len(lines) - 1

		for idx, f := range files {
			if idx >= maxFileLines-3 {
				lines = append(lines, " ...and more items")
				break
			}
			if f.IsDir() {
				lines = append(lines, fmt.Sprintf(" 📂 %s/", f.Name()))
			} else {
				lines = append(lines, fmt.Sprintf(" 📄 %s", f.Name()))
			}
		}
		return strings.Join(lines, "\n")
	}

	content, err := os.ReadFile(selected.path)
	if err != nil {
		return "Binary or unreadable file content"
	}

	strContent := string(content)
	for _, b := range strContent {
		if b == 0 {
			return "Binary file content preview not available"
		}
	}

	lines := strings.Split(strContent, "\n")
	var boundedLines []string

	for idx, line := range lines {
		if idx >= height {
			break
		}
		if len(line) > width-3 && width > 3 {
			line = line[:width-6] + "..."
		}
		boundedLines = append(boundedLines, line)
	}

	return strings.Join(boundedLines, "\n")
}
