package main

import (
	"errors"
	"io/fs"
	"path/filepath"
	"strings"
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
