package main

import (
	"fmt"
	"log"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type model struct {
	step   int
	cursor int
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
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
		}
	}
	return m, nil
}

func (m model) View() string {
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
