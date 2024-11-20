package main

import (
	"errors"
	"fmt"
	"io/fs"
	"path/filepath"
	"strings"

	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

// Finds if project with name exists, returns error if it does exist
func findProject(root, name string) error {
	return filepath.WalkDir(root, func(s string, d fs.DirEntry, e error) error {
		if e != nil {
			return e
		}

		if strings.TrimSuffix(d.Name(), ".db") == strings.ToLower(name) {
			return errors.New("Project already exists")
		}

		return nil
	})
}

// Returns an array of all the projects
func findProjects() ([]string, error) {
	var a []string

	err := filepath.WalkDir(".", func(s string, d fs.DirEntry, fErr error) error {
		if fErr != nil {
			return fErr
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

type Flashcard struct {
	Key   string
	Value string
}

func getFlashcards(project string) ([]Flashcard, error) {
	projectName := "./" + project + ".db"
	db, err := sqlx.Connect("sqlite3", projectName)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to the database: %v", err)
	}

	// TODO: remove this
	_, _ = db.Exec(`DELETE FROM flashcards`)

	// TODO: remove this
	_, err = db.Exec(`INSERT INTO flashcards (key, value)
		VALUES ("test", "value")`)
	if err != nil {
		return nil, fmt.Errorf("error 2", err)
	}

	var flashcards []Flashcard
	err = db.Select(&flashcards,
		`SELECT key, value 
		FROM flashcards`)
	if err != nil {
		return nil, err
	}

	return flashcards, nil
}
