package repo

import (
	"context"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/data/postgres"
)

// RepositoryTopic defines the interface for topic data operations
type RepositoryTopic interface {
	Create(ctx context.Context, topic factcheck.Topic) (factcheck.Topic, error)
	GetByID(ctx context.Context, id string) (factcheck.Topic, error)
	List(ctx context.Context) ([]factcheck.Topic, error)
	ListByStatus(ctx context.Context, status factcheck.StatusTopic) ([]factcheck.Topic, error)
	Update(ctx context.Context, topic factcheck.Topic) (factcheck.Topic, error)
	Delete(ctx context.Context, id string) error
}

// repositoryTopic implements RepositoryTopic
type repositoryTopic struct {
	queries *postgres.Queries
}

// NewRepositoryTopic creates a new topic repository
func NewRepositoryTopic(queries *postgres.Queries) RepositoryTopic {
	return &repositoryTopic{
		queries: queries,
	}
}

// Create creates a new topic using the topic adapter
func (r *repositoryTopic) Create(ctx context.Context, t factcheck.Topic) (factcheck.Topic, error) {
	params, err := topic(t)
	if err != nil {
		return factcheck.Topic{}, err
	}

	dbTopic, err := r.queries.CreateTopic(ctx, params)
	if err != nil {
		return factcheck.Topic{}, err
	}

	return topicDomain(dbTopic), nil
}

// GetByID retrieves a topic by ID using the topicDomain adapter
func (r *repositoryTopic) GetByID(ctx context.Context, id string) (factcheck.Topic, error) {
	topicID, err := uuid(id)
	if err != nil {
		return factcheck.Topic{}, err
	}

	dbTopic, err := r.queries.GetTopic(ctx, topicID)
	if err != nil {
		return factcheck.Topic{}, err
	}

	return topicDomain(dbTopic), nil
}

// List retrieves all topics using the topicDomain adapter
func (r *repositoryTopic) List(ctx context.Context) ([]factcheck.Topic, error) {
	dbTopics, err := r.queries.ListTopics(ctx)
	if err != nil {
		return nil, err
	}

	topics := make([]factcheck.Topic, len(dbTopics))
	for i, dbTopic := range dbTopics {
		topics[i] = topicDomain(dbTopic)
	}

	return topics, nil
}

// ListByStatus retrieves topics by status using the topicDomain adapter
func (r *repositoryTopic) ListByStatus(ctx context.Context, status factcheck.StatusTopic) ([]factcheck.Topic, error) {
	dbTopics, err := r.queries.ListTopicsByStatus(ctx, string(status))
	if err != nil {
		return nil, err
	}

	topics := make([]factcheck.Topic, len(dbTopics))
	for i, dbTopic := range dbTopics {
		topics[i] = topicDomain(dbTopic)
	}

	return topics, nil
}

// Update updates a topic using the topicUpdate adapter
func (r *repositoryTopic) Update(ctx context.Context, t factcheck.Topic) (factcheck.Topic, error) {
	params, err := topicUpdate(t)
	if err != nil {
		return factcheck.Topic{}, err
	}

	dbTopic, err := r.queries.UpdateTopic(ctx, params)
	if err != nil {
		return factcheck.Topic{}, err
	}

	return topicDomain(dbTopic), nil
}

// Delete deletes a topic by ID using the stringToUUID adapter
func (r *repositoryTopic) Delete(ctx context.Context, id string) error {
	topicID, err := uuid(id)
	if err != nil {
		return err
	}

	return r.queries.DeleteTopic(ctx, topicID)
}
