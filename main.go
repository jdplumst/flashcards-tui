package main

import (
	"log"

	tea "github.com/charmbracelet/bubbletea"
)

type state int

const (
	initial state = iota
	create
	existing
	review
	edit
)

type model struct {
	state    state
	initial  initial_model
	create   create_model
	existing existing_model
	review   review_model
	edit     edit_model
	project  string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case initial:
		cmd, s := m.initial.UpdateInitial(msg)
		m.state = s
		return m, cmd

	case create:
		cmd, s := m.create.UpdateCreate(msg)
		m.state = s
		return m, cmd

	case existing:
		cmd, s, project := m.existing.UpdateExisting(msg)
		m.state = s
		m.project = project
		return m, cmd

	case review:
		cmd, s := m.review.UpdateReview(msg)
		m.state = s
		return m, cmd

	case edit:
		cmd, s := m.edit.UpdateEdit(msg)
		m.state = s
		return m, cmd

	default:
		m.state = 0
	}

	return m, nil
}

func (m model) View() string {
	switch m.state {
	case initial:
		return m.initial.ViewInitial()

	case create:
		return m.create.ViewCreate()

	case existing:
		return m.existing.ViewExisting()

	case review:
		m.review.project = m.project
		return m.review.ViewReview()

	case edit:
		m.edit.project = m.project
		return m.edit.ViewEdit()
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
