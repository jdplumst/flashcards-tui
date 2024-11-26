package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type add int

const (
	off add = iota
	key
	value
)

type edit_model struct {
	project    string
	flashcards []Flashcard
	cursor     int
	add        add
	add_key    string
	add_value  string
	delete     Flashcard
	err        error
}

func (e *edit_model) UpdateEdit(msg tea.Msg) (tea.Cmd, state) {
	model_state := state(edit)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return tea.Quit, 0

		case "up", "k":
			switch e.add {
			case add(off):
				if e.cursor > 0 {
					e.cursor--
				}
			case add(key):
				if len(msg.String()) == 1 {
					e.add_key += msg.String()
				}
			case add(value):
				if len(msg.String()) == 1 {
					e.add_value += msg.String()
				}
			}

		case "down", "j":
			switch e.add {
			case add(off):
				if e.cursor < len(e.flashcards)-1 {
					e.cursor++
				}
			case add(key):
				if len(msg.String()) == 1 {
					e.add_key += msg.String()
				}
			case add(value):
				if len(msg.String()) == 1 {
					e.add_value += msg.String()
				}
			}

		case "enter":
			switch e.add {
			case add(key):
				e.add = add(value)
			case add(value):
				e.err = addFlashcard(e.project, e.add_key, e.add_value)
				if e.err == nil {
					e.flashcards, e.err = getFlashcards(e.project)
					e.add = add(off)
					e.add_key = ""
					e.add_value = ""
				}
			}

		case "backspace":
			switch e.add {
			case add(key):
				if len(e.add_key) > 0 {
					e.add_key = e.add_key[:len(e.add_key)-1]
				}
			case add(value):
				if len(e.add_value) > 0 {
					e.add_value = e.add_value[:len(e.add_value)-1]
				}
			}

		case "a":
			switch e.add {
			case add(off):
				e.add = add(key)
			case add(key):
				e.add_key += msg.String()
			case add(value):
				e.add_value += msg.String()
			}

		case "d":
			switch e.delete.Key {
			case "":
				switch e.add {
				case add(off):
					e.delete = e.flashcards[e.cursor]
				case add(key):
					if len(msg.String()) == 1 {
						e.add_key += msg.String()
					}
				case add(value):
					if len(msg.String()) == 1 {
						e.add_value += msg.String()
					}
				}
			default:
				switch e.add {
				case add(key):
					if len(msg.String()) == 1 {
						e.add_key += msg.String()
					}
				case add(value):
					if len(msg.String()) == 1 {
						e.add_value += msg.String()
					}
				}
			}

		case "y":
			switch e.delete.Key {
			case "":
				switch e.add {
				case add(key):
					if len(msg.String()) == 1 {
						e.add_key += msg.String()
					}
				case add(value):
					if len(msg.String()) == 1 {
						e.add_value += msg.String()
					}
				}
			default:
				e.err = deleteFlashcard(e.project, e.delete.Key)
				if e.err == nil {
					e.delete = Flashcard{Key: "", Value: ""}
				}
			}

		case "n":
			switch e.delete.Key {
			case "":
				switch e.add {
				case add(key):
					if len(msg.String()) == 1 {
						e.add_key += msg.String()
					}
				case add(value):
					if len(msg.String()) == 1 {
						e.add_value += msg.String()
					}
				}
			default:
				e.delete = Flashcard{Key: "", Value: ""}
			}

		default:
			switch e.add {
			case add(key):
				if len(msg.String()) == 1 {
					e.add_key += msg.String()
				}
			case add(value):
				if len(msg.String()) == 1 {
					e.add_value += msg.String()
				}

			}
		}
	}
	return nil, model_state
}

func (e *edit_model) ViewEdit() string {
	s := ""
	s += lipgloss.NewStyle().
		SetString("Edit Flashcards for", e.project).
		Foreground(lipgloss.Color("3")).
		Bold(true).
		Italic(true).String()
	s += "\n"

	s += lipgloss.NewStyle().
		SetString("Press (a) to add a new flashcard").
		Foreground(lipgloss.Color("3")).
		Bold(true).
		String()
	s += "\n"

	s += lipgloss.NewStyle().
		SetString("Press (d) to delete a flashcard").
		Foreground(lipgloss.Color("3")).
		Bold(true).
		String()
	s += "\n"

	s += lipgloss.NewStyle().
		SetString("Press (e) to edit a flashcard").
		Foreground(lipgloss.Color("3")).
		Bold(true).
		String()
	s += "\n"

	s += "\n"

	for idx, flashcard := range e.flashcards {
		cursor := "  "
		if e.cursor == idx {
			cursor = "> "
		}
		s += cursor + flashcard.Key + ", " + flashcard.Value
		s += "\n"
	}

	if e.add != add(off) {
		s += "\n"
		s += lipgloss.NewStyle().
			SetString("You are adding a new flashcard.").
			Foreground(lipgloss.Color("3")).
			Bold(true).String()
		s += "\n"
		s += lipgloss.NewStyle().
			SetString("Key: ").
			String()
		s += lipgloss.NewStyle().
			SetString(e.add_key).
			Blink(e.add == add(key)).
			String()
		s += "\n"
		s += lipgloss.NewStyle().
			SetString("Value: ").
			String()
		s += lipgloss.NewStyle().
			SetString(e.add_value).
			Blink(e.add == add(value)).
			String()
	}

	if e.delete.Key != "" {
		s += "\n"
		s += lipgloss.NewStyle().
			SetString("Are you sure you want to delete {" + e.delete.Key + ", " + e.delete.Value + "}?").
			Foreground(lipgloss.Color("1")).
			String()
		s += "\n"
		s += lipgloss.NewStyle().
			SetString("Press (y) for yes, (n) for no").
			Foreground(lipgloss.Color("1")).
			String()
	}

	if e.err != nil {
		s += "\n"
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
