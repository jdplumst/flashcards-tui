package main

import (
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type existing_model struct {
	name          string
	cursor_list   int
	cursor_prompt int
	err           error
}

var projects, err = findProjects()

func (e *existing_model) UpdateExisting(msg tea.Msg) (tea.Cmd, state, string) {
	model_state := state(existing)

	if err != nil {
		e.err = err
	}

	switch e.name {
	case "":

		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c":
				return tea.Quit, model_state, ""

			case "up", "k":
				switch e.err {
				case nil:
					if e.cursor_list > 0 {
						e.cursor_list--
					}
				default:
					e.err = nil
					e.name = ""
					e.cursor_list = 0
					e.cursor_prompt = 0
				}

			case "down", "j":
				switch e.err {
				case nil:
					if e.cursor_list < len(projects)-1 {
						e.cursor_list++

					}
				default:
					e.err = nil
					e.name = ""
					e.cursor_list = 0
					e.cursor_prompt = 0
				}

			case "enter":
				e.err = nil
				e.name = projects[e.cursor_list]
				e.cursor_prompt = 0

			default:
				e.err = nil
				e.name = ""
				e.cursor_list = 0
				e.cursor_prompt = 0
			}
		}

	default:
		switch msg := msg.(type) {
		case tea.KeyMsg:
			switch msg.String() {
			case "ctrl+c":
				return tea.Quit, model_state, ""

			case "up", "k":
				switch e.err {
				case nil:
					if e.cursor_prompt > 0 {
						e.cursor_prompt--
					}
				default:
					e.err = nil
					e.name = ""
					e.cursor_list = 0
					e.cursor_prompt = 0
				}

			case "down", "j":
				switch e.err {
				case nil:
					if e.cursor_prompt < 1 {
						e.cursor_prompt++
					}
				default:
					e.err = nil
					e.name = ""
					e.cursor_list = 0
					e.cursor_prompt = 0
				}

			case "enter":
				switch e.cursor_prompt {
				case 0:
					model_state = 3
				case 1:
					model_state = 4
				}
				project := e.name
				e.err = nil
				e.name = ""
				e.cursor_list = 0
				e.cursor_prompt = 0
				return nil, model_state, project

			}
		}

	}
	return nil, model_state, ""
}

func (e *existing_model) ViewExisting() string {
	s := ""
	s += lipgloss.NewStyle().SetString("EXISTING PROJECTS").
		Foreground(lipgloss.Color("2")).
		Bold(true).
		Italic(true).
		String()
	s += "\n"

	for idx, project := range projects {
		x := "[ ] "
		if idx == e.cursor_list {
			x = "[X] "
		}
		selector := lipgloss.NewStyle().
			SetString(x).
			Foreground(lipgloss.Color("4")).
			String()

		name := lipgloss.NewStyle().
			SetString(strings.TrimSuffix(project, ".db")).
			Foreground(lipgloss.Color("4")).
			String()

		s += selector
		s += name
		s += "\n"
	}

	if e.name != "" {
		s += "\n"
		s += lipgloss.NewStyle().
			SetString("You have selected", strings.TrimSuffix(e.name, ".db")).
			Foreground(lipgloss.Color("2")).
			Bold(true).String()
		s += "\n"
		s += lipgloss.NewStyle().
			SetString("Do you want to perform a review or edit these flashcards?").
			Foreground(lipgloss.Color("2")).
			Bold(true).
			String()
		s += "\n"

		choices := [2]string{
			lipgloss.NewStyle().
				SetString("Perform a review").
				Foreground(lipgloss.Color("201")).
				Italic(e.cursor_prompt == 0).
				String(),
			lipgloss.NewStyle().
				SetString("Edit the flashcards").
				Foreground(lipgloss.Color("3")).
				Italic(e.cursor_prompt == 1).
				String(),
		}

		for idx, choice := range choices {
			cursor := " "
			if e.cursor_prompt == idx {
				cursor = ">"
			}

			s += fmt.Sprintf("%s %s\n", cursor, choice)
		}
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

func findProjects() ([]string, error) {
	var a []string

	err := filepath.WalkDir(".", func(s string, d fs.DirEntry, fErr error) error {
		if fErr != nil {
			return fErr
		}
		if filepath.Ext(d.Name()) == ".db" {
			a = append(a, s)
		}
		return nil
	})

	if err != nil {
		return nil, err
	}

	return a, nil
}
