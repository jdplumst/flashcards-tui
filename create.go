package main

import (
	"errors"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type create_model struct {
	name string
	err  error
}

func (c *create_model) UpdateCreate(msg tea.Msg) (tea.Cmd, state) {
	model_state := state(create)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return tea.Quit, model_state

		case "enter":
			switch c.err {
			case nil:
				if len(c.name) <= 0 {
					c.err = errors.New("Project name must be at least 1 character long")
				}
			// TODO: else create project
			default:
				c.err = nil
				c.name = ""
			}

		case "backspace":
			switch c.err {
			case nil:
				if len(c.name) > 0 {
					c.name = c.name[:len(c.name)-1]
				}
			default:
				c.err = nil
				c.name = ""
			}

		default:
			switch c.err {
			case nil:
				// Prevent things like "ctrl+a" from being appended to the name
				if len(msg.String()) == 1 {
					c.name += msg.String()
				}
			default:
				c.err = nil
				c.name = ""
			}
		}

	}

	return nil, model_state

}

func (c *create_model) ViewCreate() string {
	s := ""
	s += lipgloss.NewStyle().
		SetString("CREATING A NEW SET OF FLASHCARDS").
		Foreground(lipgloss.Color("4")).
		Bold(true).
		Italic(true).
		String()
	s += "\n"

	s += lipgloss.NewStyle().
		SetString("What do you want your project to be called?", c.name).
		String()

	if c.err != nil {
		s += "\n"
		s += lipgloss.NewStyle().
			SetString(c.err.Error()).
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
