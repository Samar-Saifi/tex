package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

const (
	BackgroundDark = "#1E222A"
	SurfaceMedium  = "#282C34"
	TextMain       = "#ABB2BF"
	TextMuted      = "#5C6370"
	AccentTeal     = "#4C566A"
	FocusTeal      = "#00E676"
	FocusBlue      = "#61AFEF"
	AlertCoral     = "#E06C75"
	AlertWarning   = "#D19A66"
)

var (
	NormalHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#1E222A")).
				Background(lipgloss.Color(FocusBlue)).
				Padding(0, 1)

	SearchHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#1E222A")).
				Background(lipgloss.Color(AlertWarning)).
				Padding(0, 1)

	RenameHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#1E222A")).
				Background(lipgloss.Color(FocusTeal)).
				Padding(0, 1)

	PopupHeaderStyle = lipgloss.NewStyle().
				Bold(true).
				Foreground(lipgloss.Color("#1E222A")).
				Background(lipgloss.Color(AlertCoral)).
				Padding(0, 1)

	ListStyle = lipgloss.NewStyle().
			PaddingRight(1)

	PreviewStyle = lipgloss.NewStyle().
			Border(lipgloss.NormalBorder(), false, false, false, true).
			BorderForeground(lipgloss.Color(TextMuted)).
			PaddingLeft(2).
			Foreground(lipgloss.Color(TextMain)).
			Background(lipgloss.Color(SurfaceMedium))

	CursorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(FocusTeal)).
			Bold(true)

	DirStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(FocusBlue)).
			Bold(true)

	FileStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(TextMain))

	SearchInputBoxStyle = lipgloss.NewStyle().
				Border(lipgloss.RoundedBorder()).
				BorderForeground(lipgloss.Color(AlertWarning)).
				Padding(0, 1).
				MarginTop(1)

	SearchPromptStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color(AlertWarning)).
				Bold(true)

	SearchValueStyle = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color(SurfaceMedium))

	PopupStyle = lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color(AlertCoral)).
			Padding(1, 2).
			Background(lipgloss.Color(SurfaceMedium))

	ActiveButtonSpec = lipgloss.NewStyle().
				Foreground(lipgloss.Color("#FAFAFA")).
				Background(lipgloss.Color(AlertCoral)).
				Padding(0, 2).
				Bold(true)

	InactiveButtonSpec = lipgloss.NewStyle().
				Foreground(lipgloss.Color(TextMain)).
				Background(lipgloss.Color(BackgroundDark)).
				Padding(0, 2)

	MetadataStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(TextMain)).
			Border(lipgloss.NormalBorder(), true, false, false, false).
			BorderForeground(lipgloss.Color(TextMuted)).
			PaddingTop(1)

	FooterStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color(TextMuted)).
			PaddingTop(0)

	ErrorStyle = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#FFFFFF")).
			Background(lipgloss.Color(AlertCoral)).
			Bold(true).
			Padding(0, 1)
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

	header := renderHeader(m, width)

	mainContent := renderWorkspace(m, width, contentHeight)

	if m.currentMode == search {
		mainContent = renderSearchOverlay(m, width, mainContent)
	} else if m.currentMode == popup && ActivePopup != nil {
		mainContent = renderPopupOverlay(m, width, contentHeight)
	}

	metadata := MetadataStyle.Width(width).Render(renderMetadataString(m))
	footer := FooterStyle.Width(width).Render(renderFooterString(m))

	parts := []string{header, mainContent}

	if m.errorMsg != "" {
		parts = append(parts, ErrorStyle.Width(width).Render("⚠ "+m.errorMsg))
	}

	parts = append(parts, metadata, footer)
	return lipgloss.JoinVertical(lipgloss.Left, parts...)
}

func renderHeader(m model, width int) string {
	switch m.currentMode {
	case search:
		return SearchHeaderStyle.Width(width).Render(" MODE: SEARCHING DIRECTORY")
	case rename:
		return RenameHeaderStyle.Width(width).Render(" MODE: RENAMING...")
	case popup:
		return PopupHeaderStyle.Width(width).Render(" MODE: ACTION REQUIRED")
	default:
		return NormalHeaderStyle.Width(width).Render(" TEX")
	}
}

func renderWorkspace(m model, width, contentHeight int) string {
	leftWidth := int(float64(width) * 0.4)
	rightWidth := width - leftWidth - 1

	leftContent := RenderVerticalFileList(m, leftWidth, contentHeight)
	rightContent := RenderFilePreview(m, rightWidth, contentHeight)

	return lipgloss.JoinHorizontal(
		lipgloss.Top,
		ListStyle.Width(leftWidth).Height(contentHeight).Render(leftContent),
		PreviewStyle.Width(rightWidth).Height(contentHeight).Render(rightContent),
	)
}

func renderSearchOverlay(m model, width int, currentWorkspace string) string {
	displayText := m.searchQuery + "█"
	inputContent := fmt.Sprintf("%s %s", SearchPromptStyle.Render("🔍 Search for:"), SearchValueStyle.Render(displayText))
	searchField := SearchInputBoxStyle.Width(width - 2).Render(inputContent)

	return lipgloss.JoinVertical(lipgloss.Left, currentWorkspace, searchField)
}

func renderPopupOverlay(m model, width, contentHeight int) string {
	boxWidth := width * 7 / 10
	if boxWidth > 80 {
		boxWidth = 80
	}
	if boxWidth < 40 {
		boxWidth = 40
	}

	innerWidth := boxWidth - 4

	title := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("#FFFFFF")).
		Background(lipgloss.Color(SurfaceMedium)). // ✅
		Width(innerWidth).
		Render(ActivePopup.Message)

	var body string

	switch ActivePopup.Layout {
	case PopupHorizontal:
		var buttons []string
		for i, option := range ActivePopup.Options {
			if i == ActivePopup.ActiveIndex {
				buttons = append(buttons, ActiveButtonSpec.Render(option))
			} else {
				buttons = append(buttons, InactiveButtonSpec.Render(option))
			}
		}
		body = lipgloss.NewStyle().
			Width(innerWidth).
			Background(lipgloss.Color(SurfaceMedium)). // ✅
			Render(lipgloss.JoinHorizontal(lipgloss.Center, buttons...))

	case PopupVertical:
		selectedStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(FocusTeal)).
			Background(lipgloss.Color(SurfaceMedium)).
			Width(innerWidth). // ✅ pad to full width
			Bold(true)

		normalStyle := lipgloss.NewStyle().
			Foreground(lipgloss.Color(TextMain)).
			Background(lipgloss.Color(SurfaceMedium)).
			Width(innerWidth) // ✅ pad to full width

		var items []string
		for i, option := range ActivePopup.Options {
			if i == ActivePopup.ActiveIndex {
				items = append(items, selectedStyle.Render("▶ "+option))
			} else {
				items = append(items, normalStyle.Render("  "+option))
			}
		}
		body = lipgloss.JoinVertical(lipgloss.Left, items...)
		// body = lipgloss.NewStyle().
		// 	Width(innerWidth).
		// 	Background(lipgloss.Color(SurfaceMedium)). // ✅
		// 	Render(lipgloss.JoinVertical(lipgloss.Left, items...))
	}

	help := lipgloss.NewStyle().
		Foreground(lipgloss.Color(TextMuted)).
		Background(lipgloss.Color(SurfaceMedium)). // ✅
		Width(innerWidth).
		Render(func() string {
			if ActivePopup.Layout == PopupHorizontal {
				return "←/→ Select • Enter Confirm • Esc Cancel"
			}
			return "↑/↓ Select • Enter Open • Esc Cancel"
		}())

	content := lipgloss.JoinVertical(
		lipgloss.Left,
		title,
		"",
		body,
		"",
		help,
	)

	popupBox := PopupStyle.
		Width(boxWidth).
		Render(content)

	return lipgloss.Place(
		width,
		contentHeight,
		lipgloss.Center,
		lipgloss.Center,
		popupBox,
		lipgloss.WithWhitespaceBackground(lipgloss.Color(BackgroundDark)),
		lipgloss.WithWhitespaceForeground(lipgloss.Color(BackgroundDark)),
	)
}

func renderMetadataString(m model) string {
	if len(m.data) == 0 || m.cursor >= len(m.data) {
		return " No File Selected"
	}

	selected := m.data[m.cursor]
	info, err := os.Stat(selected.path)
	if err != nil {
		return " No File Selected"
	}

	sizeStr := fmt.Sprintf("%d B", info.Size())
	if info.IsDir() {
		sizeStr = "DIR"
	} else if info.Size() > 1024*1024 {
		sizeStr = fmt.Sprintf("%.2f MB", float64(info.Size())/(1024*1024))
	} else if info.Size() > 1024 {
		sizeStr = fmt.Sprintf("%.2f KB", float64(info.Size())/1024)
	}

	return fmt.Sprintf(" 📄 %s  •  💾 %s  •  🕒 %s  •  🔒 %s",
		info.Name(), sizeStr, info.ModTime().Format("2006-01-02 15:04"), info.Mode().String())
}

func renderFooterString(m model) string {
	switch m.currentMode {
	case search:
		return " [Type characters to filter]  •  [enter] focus match  •  [esc] cancel search"
	case popup:
		return " [←/→] switch choice  •  [enter] execute action  •  [esc] dismiss modal"
	default:
		return " [↑/↓] navigate  •  [enter/→] open  •  [←/backspace] back  •  [alt+s] search  •  [alt+t] terminal  •  [alt+r] rename"
	}
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
