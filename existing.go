package main

import (
	"errors"
	"io/fs"
	"path/filepath"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type existing_model struct {
	name   string
	cursor int
	err    error
}

func (e *existing_model) UpdateExisting(msg tea.Msg) (tea.Cmd, state) {
	model_state := state(existing)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return tea.Quit, model_state

		case "up", "k":
			switch e.err {
			case nil:
				if e.cursor > 0 {
					e.cursor--
				}
			default:
				e.err = nil
				e.name = ""
				e.cursor = 0
			}

		case "down", "j":
			switch e.err {
			case nil:
				projects, err := findProjects()
				if err != nil {
					e.err = err
				} else {
					if e.cursor < len(projects)-1 {
						e.cursor++
					}
				}
			default:
				e.err = nil
				e.name = ""
				e.cursor = 0
			}
		case "p":
			e.err = errors.New("yo an error")
		default:
			e.err = nil
			e.name = ""
			e.cursor = 0
		}
	}

	return nil, model_state
}

func (e *existing_model) ViewExisting() string {
	s := ""
	s += lipgloss.NewStyle().SetString("EXISTING PROJECTS").
		Foreground(lipgloss.Color("2")).
		Bold(true).
		Italic(true).
		String()
	s += "\n"

	projects, err := findProjects()
	if err != nil {
		s += lipgloss.NewStyle().
			SetString(err.Error()).
			Foreground(lipgloss.Color("1")).
			Bold(true).
			String()
		s += "\n"
		s += lipgloss.NewStyle().
			SetString("(Press any key to retry)").
			Faint(true).String()
	}

	for idx, project := range projects {
		x := "[ ] "
		if idx == e.cursor {
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
	err := filepath.WalkDir(".", func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
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