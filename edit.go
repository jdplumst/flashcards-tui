package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type add int

const (
	off add = iota
	key
	value
)

type edit_model struct {
	project    string
	flashcards []Flashcard
	add        add
	add_key    string
	add_value  string
	err        error
}

func (e *edit_model) UpdateEdit(msg tea.Msg) (tea.Cmd, state) {
	model_state := state(edit)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return tea.Quit, 0

		case "enter":
			switch e.add {
			case add(key):
				e.add = add(value)
			case add(value):
				e.err = addFlashcard(e.project, e.add_key, e.add_value)
				if e.err == nil {
					e.flashcards, e.err = getFlashcards(e.project)
					e.add = add(off)
					e.add_key = ""
					e.add_value = ""
				}
			}

		case "backspace":
			switch e.add {
			case add(key):
				if len(e.add_key) > 0 {
					e.add_key = e.add_key[:len(e.add_key)-1]
				}
			case add(value):
				if len(e.add_value) > 0 {
					e.add_value = e.add_value[:len(e.add_value)-1]
				}
			}

		case "a":
			switch e.add {
			case add(off):
				e.add = add(key)
			case add(key):
				e.add_key += msg.String()
			case add(value):
				e.add_value += msg.String()
			}

		default:
			switch e.add {
			case add(key):
				if len(msg.String()) == 1 {
					e.add_key += msg.String()
				}
			case add(value):
				if len(msg.String()) == 1 {
					e.add_value += msg.String()
				}

			}
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

	for _, flashcard := range e.flashcards {
		s += flashcard.Key + ", " + flashcard.Value
		s += "\n"
	}

	if e.add != add(off) {
		s += "\n"
		s += lipgloss.NewStyle().
			SetString("You are adding a new flashcard.").
			Foreground(lipgloss.Color("3")).
			Bold(true).String()
		s += "\n"
		s += lipgloss.NewStyle().
			SetString("Key: ").
			String()
		s += lipgloss.NewStyle().
			SetString(e.add_key).
			Blink(e.add == add(key)).
			String()
		s += "\n"
		s += lipgloss.NewStyle().
			SetString("Value: ").
			String()
		s += lipgloss.NewStyle().
			SetString(e.add_value).
			Blink(e.add == add(value)).
			String()
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
