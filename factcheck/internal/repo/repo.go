package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5"

	"github.com/kaogeek/line-fact-check/factcheck/data/postgres"
)

// Repository combines all repository interfaces
type Repository struct {
	Topics       Topics
	Messages     Messages
	UserMessages UserMessages
	TxnManager   postgres.TxnManager
}

type Tx interface {
	Commit(context.Context) error
	Rollback(context.Context) error
}

// ErrNotFound is returned when a requested resource is not found
type ErrNotFound struct {
	Err    error `json:"-"` // Prevent leaks?
	Filter any   `json:"filter"`
}

// New creates a new repository with all implementations
func New(queries *postgres.Queries, conn *pgx.Conn) Repository {
	return Repository{
		Topics:       NewTopics(queries),
		Messages:     NewMessages(queries),
		UserMessages: NewUserMessages(queries),
		TxnManager:   postgres.NewTxnManager(conn),
	}
}

func (r *Repository) Begin(ctx context.Context) (Tx, error) {
	return r.TxnManager.Begin(ctx)
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

// substring surrounds the pattern with % for LIKE queries
func substring(pattern string) string {
	return "%" + pattern + "%"
}
