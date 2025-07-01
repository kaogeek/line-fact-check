package repo

import (
	"github.com/kaogeek/line-fact-check/factcheck/postgres/sqlcgen"
)

// Repository combines all repository interfaces
type Repository struct {
	Topic       RepositoryTopic
	Message     RepositoryMessage
	UserMessage RepositoryUserMessage
}

// NewRepository creates a new repository with all implementations
func NewRepository(queries *postgres.Queries) *Repository {
	return &Repository{
		Topic:       NewRepositoryTopic(queries),
		Message:     NewRepositoryMessage(queries),
		UserMessage: NewRepositoryUserMessage(queries),
	}
}
