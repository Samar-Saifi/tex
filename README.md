# Tex

A keyboard-driven terminal based file manager. Navigate your filesystem, preview files, open applications, rename, delete, all without leaving the terminal.

Built as a project to explore TUI development in Go using the Bubble Tea framework.

---

## Preview

<img width="930" height="600" alt="image" src="https://github.com/user-attachments/assets/bda708dc-c86a-47ba-9f13-c7c490d23368" />

<img width="930" height="600" alt="image" src="https://github.com/user-attachments/assets/da796fd3-daff-4b2d-8dd6-9bf56384e200" />


---

## Features

- Live directory preview — hovering a folder shows its contents in the right pane without entering it
- Smart PATH handling — the installer detects your shell (bash/zsh/fish) and patches the right rc file
- Search without a shortcut — just start typing, no mode switch needed
- Rename and delete files
- xdg-mime integration — Open With doesn't just guess; it reads your installed .desktop files and matches them to the file's actual MIME type

---

## Installation

### Linux

**Requirements:** Go 1.21+

```bash
git clone https://github.com/yourusername/tex
cd tex
bash install.sh
```

Installs to `/usr/local/bin` (if run as root) or `~/.local/bin` (user install). The script adds `~/.local/bin` to your PATH automatically if needed.

---

### Windows

**Requirements:** Go 1.21+, run as Administrator

```bat
git clone https://github.com/yourusername/tex
cd tex
windows.bat
```

Installs to `C:\tex` and adds it to your user PATH via PowerShell.

---

> Open a **new** terminal after install, then type `tex`.

---

## Keybindings

| Key | Action |
|---|---|
| `↑` / `↓` | Move cursor |
| `→` / `Enter` | Open file or enter directory |
| `←` / `Backspace` | Go to parent directory |
| `Alt+S` | Search / filter current directory |
| `Alt+R` | Rename selected file |
| `Alt+T` | Open terminal here |
| `Alt+O` | Open with — choose application |
| `Delete` | Delete (asks for confirmation) |
| `Esc` | Cancel / close popup |
| `Alt+Q` | Quit |

> Tip: typing any character in normal mode jumps straight into search.

---

## Tech stack

- [Bubble Tea](https://github.com/charmbracelet/bubbletea) — TUI framework
- [Lip Gloss](https://github.com/charmbracelet/lipgloss) — terminal styling
- Go standard library for filesystem operations

---

## Building from source

```bash
cd src
go mod tidy
go build -o ../build/tex .
```
