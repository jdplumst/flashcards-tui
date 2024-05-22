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
)

type model struct {
	state   state
	initial initial_model
	create  create_model
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
		switch msg := msg.(type) {
		case tea.KeyMsg:
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
		return m.initial.ViewInitial()

	case create:
		return m.create.ViewCreate()

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
