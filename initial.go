package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type initial_model struct {
	cursor int
}

func (i *initial_model) UpdateInitial(msg tea.Msg) (tea.Cmd, state) {
	model_state := state(initial)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return tea.Quit, model_state

		case "up", "k":
			if i.cursor > 0 {
				i.cursor--
			}

		case "down", "j":
			if i.cursor < 1 {
				i.cursor++
			}

		case "enter":
			switch i.cursor {
			case 0:
				model_state = state(create)
			case 1:
				model_state = state(existing)
			}
			i.cursor = 0
		}

	}

	return nil, model_state
}

func (i *initial_model) ViewInitial() string {
	s := ""
	s += lipgloss.NewStyle().
		SetString(" What would you like to do?").
		Foreground(lipgloss.Color("171")).
		Bold(true).
		String()
	s += "\n" // add this separately from previous line to work better with lipgloss

	choices := [5]string{
		lipgloss.NewStyle().
			SetString("Create a new set of flashcards").
			Foreground(lipgloss.Color("4")).
			Italic(i.cursor == 0).
			String(),
		lipgloss.NewStyle().
			SetString("Use an existing set of flashcards").
			Foreground(lipgloss.Color("2")).
			Italic(i.cursor == 1).
			String(),
	}

	for idx, choice := range choices {
		cursor := " "
		if i.cursor == idx {
			cursor = ">"
		}

		s += fmt.Sprintf("%s %s\n", cursor, choice)
	}
	return s

}
