package main

import tea "github.com/charmbracelet/bubbletea"

type review_model struct {
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
	return "This is the review view!"
}
