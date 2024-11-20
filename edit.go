package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type edit_model struct {
	project string
	err     error
}

var flashcards, editErr = getFlashcards("hi")

func (e *edit_model) UpdateEdit(msg tea.Msg) (tea.Cmd, state) {
	model_state := state(edit)

	if editErr != nil {
		e.err = editErr
	}

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
	s += "\n"

	for _, flashcard := range flashcards {
		s += flashcard.Key + ", " + flashcard.Value
		s += "\n"
	}

	if e.err != nil {
		s += lipgloss.NewStyle().
			SetString(e.err.Error()).
			Foreground(lipgloss.Color("1")).
			Bold(true).
			String()
		s += "\n"
		s += lipgloss.NewStyle().
			SetString("(Press any key to retry)").
			Faint(true).String()
	}

	return s
}
