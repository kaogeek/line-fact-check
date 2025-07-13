package repo

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/kaogeek/line-fact-check/factcheck/data/postgres"
)

// Repository combines all repository interfaces
type Repository struct {
	Topics       Topics
	Messages     Messages
	UserMessages UserMessages
}

// ErrNotFound is returned when a requested resource is not found
type ErrNotFound struct {
	Err    error `json:"-"` // Prevent leaks?
	Filter any   `json:"filter"`
}

// New creates a new repository with all implementations
func New(queries *postgres.Queries) Repository {
	return Repository{
		Topics:       NewTopics(queries),
		Messages:     NewMessages(queries),
		UserMessages: NewUserMessages(queries),
	}
}

// IsNotFound checks if the error is a not found error
func IsNotFound(err error) bool {
	return errors.Is(err, &ErrNotFound{})
}

func (e *ErrNotFound) Error() string {
	return fmt.Sprintf("not found for filter %+v", e.Filter)
}

func (e *ErrNotFound) Unwrap() error {
	return e.Err
}

// Is allows errors.Is to work with *ErrNotFound
func (e *ErrNotFound) Is(target error) bool {
	_, ok := target.(*ErrNotFound)
	return ok
}

func handleNotFound(err error, filter map[string]string) error {
	if errors.Is(err, sql.ErrNoRows) {
		return &ErrNotFound{
			Err:    err,
			Filter: filter,
		}
	}
	return err
}
