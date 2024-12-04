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

// Struct for flashcards in databases
type Flashcard struct {
	Key   string
	Value string
}

// Returns all flashcards for a given project
func getFlashcards(project string) ([]Flashcard, error) {
	db, err := sqlx.Connect("sqlite3", project)
	if err != nil {
		return nil, fmt.Errorf("Error connecting to the database: %v", err)
	}

	var flashcards []Flashcard
	err = db.Select(&flashcards, `
		SELECT key, value 
		FROM flashcards
		`)
	if err != nil {
		return nil, err
	}

	db.Close()

	return flashcards, nil
}

// Inserts key and value into db for project
func addFlashcard(project, key, value string) error {
	db, err := sqlx.Connect("sqlite3", project)
	if err != nil {
		return fmt.Errorf("Error connecting to the database: %v", err)
	}

	_, err = db.Exec(`
		INSERT INTO flashcards (key, value)
		VALUES (?, ?)
		`,
		key,
		value)
	if err != nil {
		return fmt.Errorf("Error adding flashcard: %v", err)
	}

	db.Close()

	return nil
}

// Deletes flashcard with key from project
func deleteFlashcard(project, key string) error {
	db, err := sqlx.Connect("sqlite3", project)
	if err != nil {
		return fmt.Errorf("Error connecting to the database: %v", err)
	}

	_, err = db.Exec(`
		DELETE FROM flashcards
		WHERE key = ?
		`,
		key)
	if err != nil {
		return fmt.Errorf("Error deleting flashcard: %v", err)
	}

	db.Close()

	return nil
}

// Edits flashcard with key from project with new_key and new_value
func editFlashcard(project, key, new_key, new_value string) error {
	db, err := sqlx.Connect("sqlite3", project)
	if err != nil {
		return fmt.Errorf("Error connecting to the database: %v", err)
	}

	_, err = db.Exec(`
		UPDATE flashcards
		SET key = ?, value = ?
		WHERE key = ?
		`,
		new_key,
		new_value,
		key)
	if err != nil {
		return fmt.Errorf("Error editing flashcard: %v", err)
	}

	db.Close()

	return nil
}
