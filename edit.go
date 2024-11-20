package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type edit_model struct {
	project string
}

func (e *edit_model) UpdateEdit(msg tea.Msg) (tea.Cmd, state) {
	model_state := state(edit)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return tea.Quit, 0
		}
	}
	return nil, model_state
}

func (e *edit_model) ViewEdit() string {
	s := ""
	s += lipgloss.NewStyle().
		SetString("Edit Flashcards for", e.project).
		Foreground(lipgloss.Color("3")).
		Bold(true).
		Italic(true).String()
	return s
}
