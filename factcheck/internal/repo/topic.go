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
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status factcheck.StatusTopic) (factcheck.Topic, error)
	UpdateDescription(ctx context.Context, id string, description string) (factcheck.Topic, error)
	UpdateName(ctx context.Context, id string, name string) (factcheck.Topic, error)
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

// Delete deletes a topic by ID using the stringToUUID adapter
func (r *repositoryTopic) Delete(ctx context.Context, id string) error {
	topicID, err := uuid(id)
	if err != nil {
		return err
	}

	return r.queries.DeleteTopic(ctx, topicID)
}

func (r *repositoryTopic) UpdateStatus(ctx context.Context, id string, status factcheck.StatusTopic) (factcheck.Topic, error) {
	topicID, err := uuid(id)
	if err != nil {
		return factcheck.Topic{}, err
	}
	dbTopic, err := r.queries.UpdateTopicStatus(ctx, postgres.UpdateTopicStatusParams{
		ID:     topicID,
		Status: string(status),
	})
	if err != nil {
		return factcheck.Topic{}, err
	}
	return topicDomain(dbTopic), nil
}

func (r *repositoryTopic) UpdateDescription(ctx context.Context, id string, description string) (factcheck.Topic, error) {
	topicID, err := uuid(id)
	if err != nil {
		return factcheck.Topic{}, err
	}
	dbTopic, err := r.queries.UpdateTopicDescription(ctx, postgres.UpdateTopicDescriptionParams{
		ID:          topicID,
		Description: description,
	})
	if err != nil {
		return factcheck.Topic{}, err
	}
	return topicDomain(dbTopic), nil
}

func (r *repositoryTopic) UpdateName(ctx context.Context, id string, name string) (factcheck.Topic, error) {
	topicID, err := uuid(id)
	if err != nil {
		return factcheck.Topic{}, err
	}
	dbTopic, err := r.queries.UpdateTopicName(ctx, postgres.UpdateTopicNameParams{
		ID:   topicID,
		Name: name,
	})
	if err != nil {
		return factcheck.Topic{}, err
	}
	return topicDomain(dbTopic), nil
}
