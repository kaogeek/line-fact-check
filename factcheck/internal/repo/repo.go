package repo

import (
	"github.com/kaogeek/line-fact-check/factcheck/models/postgres"
)

// Repository combines all repository interfaces
type Repository struct {
	Topic       RepositoryTopic
	Message     RepositoryMessage
	UserMessage RepositoryUserMessage
}

// New creates a new repository with all implementations
func New(queries *postgres.Queries) *Repository {
	return &Repository{
		Topic:       NewRepositoryTopic(queries),
		Message:     NewRepositoryMessage(queries),
		UserMessage: NewRepositoryUserMessage(queries),
	}
}
