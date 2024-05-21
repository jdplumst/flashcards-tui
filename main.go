package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type state int

const (
	initial state = iota
	create
	existing
)

type model struct {
	state  state
	cursor int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch m.state {
		case initial:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit

			case "up", "k":
				if m.cursor > 0 {
					m.cursor--
				}

			case "down", "j":
				if m.cursor < 1 {
					m.cursor++
				}

			case "enter":
				m.state = state(m.cursor + 1)
				m.cursor = 0
			}
		case create:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			}
		case existing:
			switch msg.String() {
			case "ctrl+c":
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m model) View() string {
	switch m.state {
	case initial:

		s := ""
		s += lipgloss.NewStyle().
			SetString("Do you want to create a new set of flashcards or use an existing set?").
			Foreground(lipgloss.Color("171")).
			Bold(true).
			String()
		s += "\n" // add this separately from previous line to work better with lipgloss

		choices := [2]string{
			lipgloss.NewStyle().
				SetString("Create a new set of flashcards").
				Foreground(lipgloss.Color("4")).
				Italic(m.cursor == 0).
				String(),
			lipgloss.NewStyle().
				SetString("Use an existing set of flashcards").
				Foreground(lipgloss.Color("2")).
				Italic(m.cursor == 1).
				String(),
		}

		for i, choice := range choices {
			cursor := " "
			if m.cursor == i {
				cursor = ">"
			}

			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}
		return s

	// TODO:
	case create:
		return "Create a new set of flashcards!"

	case existing:
		return "Use an existing set of flashcards!"
	}

	return "Something went wrong, please try again."
}

func main() {
	f, err := tea.LogToFile("debug.log", "debug")
	if err != nil {
		log.Fatalf("err: %w", err)
	}
	defer f.Close()
	p := tea.NewProgram(model{}, tea.WithAltScreen())
	if _, err := p.Run(); err != nil {
		log.Fatal()
	}
}
