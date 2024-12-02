package main

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type mode int

const (
	off mode = iota
	adding
	editting
)

type prompt int

const (
	key = iota
	value
)

type edit_model struct {
	project    string
	flashcards []Flashcard
	cursor     int
	mode       mode
	prompt     prompt
	key        string
	value      string
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
			switch e.mode {
			case mode(off):
				if e.cursor > 0 {
					e.cursor--
				}
			default:
				switch e.prompt {

				case prompt(key):
					if len(msg.String()) == 1 {
						e.key += msg.String()
					}
				case prompt(value):
					if len(msg.String()) == 1 {
						e.value += msg.String()
					}
				}
			}

		case "down", "j":
			switch e.mode {
			case mode(off):
				if e.cursor < len(e.flashcards)-1 {
					e.cursor++
				}
			default:
				switch e.prompt {

				case prompt(key):
					if len(msg.String()) == 1 {
						e.key += msg.String()
					}
				case prompt(value):
					if len(msg.String()) == 1 {
						e.value += msg.String()
					}
				}
			}

		case "enter":
			switch e.mode {
			case mode(off):
				break
			default:
				switch e.prompt {
				case prompt(key):
					e.prompt = prompt(value)
				case prompt(value):
					e.err = addFlashcard(e.project, e.key, e.value)
					if e.err == nil {
						e.flashcards, e.err = getFlashcards(e.project)
						e.mode = mode(off)
						e.prompt = prompt(key)
						e.key = ""
						e.value = ""
					}
				}
			}

		case "backspace":
			switch e.mode {
			case mode(off):
				break
			default:
				switch e.prompt {

				case prompt(key):
					if len(e.key) > 0 {
						e.key = e.key[:len(e.key)-1]
					}
				case prompt(value):
					if len(e.value) > 0 {
						e.value = e.value[:len(e.value)-1]
					}
				}
			}

		case "a":
			switch e.mode {
			case mode(off):
				e.mode = mode(adding)
				e.prompt = prompt(key)
			case mode(adding):
				switch e.prompt {

				case prompt(key):
					e.key += msg.String()
				case prompt(value):
					e.value += msg.String()
				}
			}

		case "d":
			switch e.delete.Key {
			case "":
				switch e.mode {
				case mode(off):
					e.delete = e.flashcards[e.cursor]
				default:
					switch e.prompt {

					case prompt(key):
						if len(msg.String()) == 1 {
							e.key += msg.String()
						}
					case prompt(value):
						if len(msg.String()) == 1 {
							e.value += msg.String()
						}
					}
				}
			}

		case "y":
			switch e.delete.Key {
			case "":
				switch e.mode {
				case mode(off):
					break
				default:
					switch e.prompt {

					case prompt(key):
						if len(msg.String()) == 1 {
							e.key += msg.String()
						}
					case prompt(value):
						if len(msg.String()) == 1 {
							e.value += msg.String()
						}
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
				switch e.mode {
				case mode(off):
					break
				default:
					switch e.prompt {

					case prompt(key):
						if len(msg.String()) == 1 {
							e.key += msg.String()
						}
					case prompt(value):
						if len(msg.String()) == 1 {
							e.value += msg.String()
						}
					}
				}
			default:
				e.delete = Flashcard{Key: "", Value: ""}
			}

		default:
			switch e.mode {
			case mode(off):
				break
			default:
				switch e.prompt {

				case prompt(key):
					if len(msg.String()) == 1 {
						e.key += msg.String()
					}
				case prompt(value):
					if len(msg.String()) == 1 {
						e.value += msg.String()
					}
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

	if e.mode == mode(adding) {
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
			SetString(e.key).
			Blink(e.prompt == prompt(key)).
			String()
		s += "\n"
		s += lipgloss.NewStyle().
			SetString("Value: ").
			String()
		s += lipgloss.NewStyle().
			SetString(e.value).
			Blink(e.prompt == prompt(value)).
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
