package main

import (
	"fmt"
	"strconv"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

type phase int

const (
	start phase = iota
	quiz
	end
)

type score int

const (
	unmarked score = iota
	correct
	incorrect
)

type Reviewcard struct {
	flashcard Flashcard
	guess     string
	score     score
}

type review_model struct {
	project     string
	phase       phase
	reviewcards []Reviewcard
	index       int
	guess       string
	err         error
}

func (r *review_model) UpdateReview(msg tea.Msg) (tea.Cmd, state) {
	model_state := state(review)

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return tea.Quit, 0

		case "enter":
			if r.err != nil {
				r.err = nil
				r.index = 0
				r.guess = ""
				r.reviewcards = nil
				r.phase = phase(start)
			} else {
				if r.phase != phase(quiz) {
					r.reviewcards = nil
					flashcards, err := getReview(r.project)
					if err != nil {
						r.err = err
						fmt.Println(err)
					} else {
						for _, flashcard := range flashcards {
							r.reviewcards = append(r.reviewcards, Reviewcard{flashcard: flashcard, score: unmarked})
						}
						r.guess = ""
						r.index = 0
						r.phase = quiz
					}
				} else {
					if r.guess == r.reviewcards[r.index].flashcard.Value {
						r.reviewcards[r.index].score = correct
					} else {
						r.reviewcards[r.index].score = incorrect
					}

					r.reviewcards[r.index].guess = r.guess
					r.guess = ""

					if r.index == len(r.reviewcards)-1 {
						r.phase = end
					} else {
						r.index++
					}
				}
			}
		default:
			if r.err != nil {
				r.err = nil
				r.index = 0
				r.guess = ""
				r.reviewcards = nil
				r.phase = phase(start)
			} else {
				if len(msg.String()) == 1 {
					r.guess += msg.String()
				}
			}
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

	s += "\n"

	if r.phase == phase(start) {
		s += lipgloss.NewStyle().
			SetString("Press Enter to start").
			String()
		return s
	}

	for i := 0; i <= r.index; i++ {
		s += lipgloss.NewStyle().
			SetString(r.reviewcards[i].flashcard.Key).
			Foreground(lipgloss.Color("205")).
			Bold(true).
			String()
		switch r.reviewcards[i].score {
		case unmarked:
			s += " - " + r.guess
		case correct:
			s += " " + r.reviewcards[i].guess + " ✅"
		case incorrect:
			s += " " + r.reviewcards[i].guess + " ❌"
		}
		s += "\n"
	}

	if r.phase == phase(end) {
		score := 0
		for _, reviewcard := range r.reviewcards {
			if reviewcard.score == correct {
				score++
			}
		}
		s += "You got " + strconv.Itoa(score) + "/" + strconv.Itoa(len(r.reviewcards)) + " correct\n"
		s += lipgloss.NewStyle().
			SetString("Press Enter to restart").
			String()
	}

	if r.err != nil {
		s += "\n"
		s += lipgloss.NewStyle().
			SetString(r.err.Error()).
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
