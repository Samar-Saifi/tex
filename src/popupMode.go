package main

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type PopupLayout int

const (
	PopupHorizontal PopupLayout = iota
	PopupVertical
)

type Popup struct {
	Message     string
	Options     []string
	ActiveIndex int
	Layout      PopupLayout
	OnSelect    func(selectedIndex int) (model, tea.Cmd)
}

var ActivePopup *Popup

func ShowPopup(m *model, message string, layout PopupLayout, options []string, callback func(int) (model, tea.Cmd)) {
	ActivePopup = &Popup{
		Message:     message,
		Options:     options,
		ActiveIndex: 0,
		Layout:      layout,
		OnSelect:    callback,
	}
	m.currentMode = popup
}

func handleKeyPopup(msg tea.KeyMsg, m model) (model, tea.Cmd) {
	if ActivePopup == nil {
		m.currentMode = normal
		return m, nil
	}

	key := strings.ToLower(msg.String())

	switch key {
	case keymap.left, keymap.up:
		if ActivePopup.ActiveIndex > 0 {
			ActivePopup.ActiveIndex--
		}

	case keymap.right, keymap.down:
		if ActivePopup.ActiveIndex < len(ActivePopup.Options)-1 {
			ActivePopup.ActiveIndex++
		}

	case keymap.confirm:
		updatedModel, cmd := ActivePopup.OnSelect(ActivePopup.ActiveIndex)
		ActivePopup = nil
		m.currentMode = normal
		return updatedModel, cmd

	case keymap.cancel:
		ActivePopup = nil
		m.currentMode = normal
	}

	return m, nil
}
