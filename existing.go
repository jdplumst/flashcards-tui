package main

import (
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
			if e.cursor > 0 {
				e.cursor--
			}

		case "down", "j":
			if e.cursor < len(findProjects())-1 {
				e.cursor++
			}
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

	projects := findProjects()

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

	return s
}

func findProjects() []string {
	var a []string
	filepath.WalkDir(".", func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}
		if filepath.Ext(d.Name()) == ".db" {
			a = append(a, s)
		}
		return nil
	})
	return a
}
