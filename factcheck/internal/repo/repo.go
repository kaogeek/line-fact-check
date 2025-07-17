// Package repo defines our common repository.
// It abstracts over sqlc generated code by providing interface and code
// to work with types defined in package factcheck
package repo

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/data/postgres"
	"github.com/kaogeek/line-fact-check/factcheck/internal/utils"
)

// Repository combines all repository interfaces
type Repository struct {
	Topics       Topics
	Messages     Messages
	UserMessages UserMessages
	TxnManager   postgres.TxnManager
}

// ErrNotFound is returned when a requested resource is not found
type ErrNotFound struct {
	Err    error `json:"-"` // Prevent leaks?
	Filter any   `json:"filter"`
}

// New creates a new repository with all implementations
func New(queries *postgres.Queries, pool *pgxpool.Pool) Repository {
	return Repository{
		Topics:       NewTopics(queries),
		Messages:     NewMessages(queries),
		UserMessages: NewUserMessages(queries),
		TxnManager:   postgres.NewTxnManager(pool),
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

// substring surrounds the pattern with % for LIKE queries
func substring(pattern string) string {
	return "%" + pattern + "%"
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

func paginationNone[T any](offset, limit int) factcheck.Pagination[T] {
	return factcheck.Pagination[T]{
		Data:   nil,
		Offset: offset,
		Limit:  limit,
		Total:  0,
	}
}

func paginate[D any, T any](
	offset int,
	limit int,
	rows []D,
	mapFn func(D) T, // Maps []row to data
	totalFromList0 func(D) int64, // Extracts total_count from row
) factcheck.Pagination[T] {
	if len(rows) == 0 {
		return paginationNone[T](offset, limit)
	}
	data := utils.MapSliceNoError(rows, mapFn)
	return factcheck.Pagination[T]{
		Data:   data,
		Offset: offset,
		Limit:  limit,
		Total:  totalFromList0(rows[0]),
	}
}
