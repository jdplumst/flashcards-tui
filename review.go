package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type review_model struct {
	project string
}

func (r *review_model) UpdateReview(msg tea.Msg) (tea.Cmd, state) {
	model_state := state(review)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return tea.Quit, 0
		}
	}
	return nil, model_state
}

func (r *review_model) ViewReview() string {
	s := ""
	s += lipgloss.NewStyle().
		SetString("Review Flashcards for", r.project).
		Foreground(lipgloss.Color("201")).
		Bold(true).
		Italic(true).String()
	return s
}
