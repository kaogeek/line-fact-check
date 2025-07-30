// Package repo defines our common repository.
// It abstracts over sqlc generated code by providing interface and code
// to work with types defined in package factcheck
package repo

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kaogeek/line-fact-check/factcheck/internal/data/postgres"
)

// Repository combines all repository interfaces
// and provides a transaction manager for beginning a transaction
type Repository struct {
	Topics        Topics
	MessagesV2    MessagesV2
	MessageGroups MessageGroups
	Answers       Answers

	TxnManager postgres.TxnManager
}

// ErrNotFound is returned when a requested resource is not found
type ErrNotFound struct {
	Err    error `json:"-"` // Prevent leaks?
	Filter any   `json:"filter"`
}

// New creates a new repository with all implementations
func New(queries *postgres.Queries, pool *pgxpool.Pool) Repository {
	return Repository{
		Topics:        NewTopics(queries),
		MessagesV2:    NewMessagesV2(queries),
		MessageGroups: NewMessageGroups(queries),
		Answers:       NewAnswers(queries),
		TxnManager:    postgres.NewTxnManager(pool),
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

type filter map[string]any

func handleNotFound(err error, filter any) error {
	if errors.Is(err, sql.ErrNoRows) {
		return &ErrNotFound{
			Err:    err,
			Filter: filter,
		}
	}
	return err
}

// substring surrounds the pattern with % for LIKE queries
func substring(pattern string) string {
	return "%" + pattern + "%"
}

func substringAuto(pattern string) string {
	if pattern != "" && !strings.Contains(pattern, "%") {
		pattern = substring(pattern)
	}
	return pattern
}

func sanitize(limit, offset int) (int, int) {
	if limit < 0 {
		limit = 0
	}
	if offset < 0 {
		offset = 0
	}
	return limit, offset
}
