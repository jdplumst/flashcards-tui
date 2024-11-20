package main

import (
	"fmt"
	"os"
	"regexp"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

type create_model struct {
	name    string
	err     error
	created bool
}

func (c *create_model) UpdateCreate(msg tea.Msg) (tea.Cmd, state) {
	model_state := state(create)

	switch c.created {
	case false:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c":
				return tea.Quit, model_state

			case "enter":
				switch c.err {
				case nil:
					if len(c.name) <= 0 {
						c.err = fmt.Errorf("Project name must be at least 1 character long")
					} else if !regexp.MustCompile(`^[A-Za-z]+$`).MatchString(c.name) {
						c.err = fmt.Errorf("Project can only contain alphabetic characters")
					} else if findProject("./", c.name) != nil {
						c.err = fmt.Errorf("Project with name %v already exists", c.name)
					} else {
						_, err := os.Create(c.name + ".db")
						if err != nil {
							c.err = fmt.Errorf("Something went wrong. Please try again.")
						}
						db, err := sqlx.Connect("sqlite3", c.name+".db")
						if err != nil {
							c.err = fmt.Errorf("Error connecting to the database: %v", err)
						}
						_, err = db.Exec(`
							CREATE TABLE IF NOT EXISTS flashcards 
							(key VARCHAR(255) PRIMARY KEY, 
							value VARCHAR(255) NOT NULL)`)
						if err != nil {
							c.err = fmt.Errorf("Error creating database.")
						}
						c.created = true
					}
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
	case true:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c":
				return tea.Quit, model_state

			default:
				c.created = false
				c.name = ""
				c.err = nil
				return nil, 0

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

	if c.created {
		s += "\n"
		s += lipgloss.NewStyle().
			SetString("You have successfully created the project", c.name+"!").
			Foreground(lipgloss.Color("4")).
			Bold(true).
			String()
		s += "\n"
		s += lipgloss.NewStyle().
			SetString("(Press any key to continue)").
			Faint(true).String()
	}

	return s

}
