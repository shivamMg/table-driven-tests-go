package api

import "errors"

type TODO struct {
	Name     string `json:"name"`
	Category string `json:"category"`
}

func (todo *TODO) Validate() error {
	if todo.Name == "" {
		return errors.New("empty name")
	}
	if todo.Category == "" {
		return errors.New("empty category")
	}
	return nil
}
