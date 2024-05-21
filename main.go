package main

import (
	"errors"
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
	cursor int    // initial state
	name   string // create state
	err    error  // create state
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

			case "enter":
				switch m.err {
				case nil:
					if len(m.name) <= 0 {
						m.err = errors.New("Project name must be at least 1 character long")
					}
				// TODO: else create project
				default:
					m.err = nil
					m.name = ""
				}

			case "backspace":
				switch m.err {
				case nil:
					if len(m.name) > 0 {
						m.name = m.name[:len(m.name)-1]
					}
				default:
					m.err = nil
					m.name = ""
				}

			default:
				switch m.err {
				case nil:
					// Prevent things like "ctrl+a" from being appended to the name
					if len(msg.String()) == 1 {
						m.name += msg.String()
					}
				default:
					m.err = nil
					m.name = ""
				}
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

	case create:
		s := ""
		s += lipgloss.NewStyle().
			SetString("CREATING A NEW SET OF FLASHCARDS").
			Foreground(lipgloss.Color("4")).
			Bold(true).
			Italic(true).
			String()
		s += "\n"

		s += lipgloss.NewStyle().
			SetString("What do you want your project to be called?", m.name).
			String()

		if m.err != nil {
			s += "\n"
			s += lipgloss.NewStyle().
				SetString(m.err.Error()).
				Foreground(lipgloss.Color("1")).
				Bold(true).
				String()
			s += "\n"
			s += lipgloss.NewStyle().
				SetString("(Press any key to retry)").
				Faint(true).String()
		}

		return s

	// TODO: set view for existing flashcards
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
