package main

import tea "github.com/charmbracelet/bubbletea"

type edit_model struct {
}

func (e *edit_model) UpdateEdit(msg tea.Msg) (tea.Cmd, state) {
	model_state := state(edit)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return tea.Quit, 0
		}
	}
	return nil, model_state
}

func (e *edit_model) ViewEdit() string {
	return "This is the edit view!"
}
