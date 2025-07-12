package repo

import (
	"context"
	"fmt"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/data/postgres"
)

// Topics defines the interface for topic data operations
type Topics interface {
	Create(ctx context.Context, topic factcheck.Topic) (factcheck.Topic, error)
	GetByID(ctx context.Context, id string) (factcheck.Topic, error)
	List(ctx context.Context) ([]factcheck.Topic, error)
	ListByStatus(ctx context.Context, status factcheck.StatusTopic) ([]factcheck.Topic, error)
	CountByStatus(ctx context.Context, status factcheck.StatusTopic) (int64, error)
	CountByStatuses(ctx context.Context) (map[factcheck.StatusTopic]int64, error)
	Delete(ctx context.Context, id string) error
	UpdateStatus(ctx context.Context, id string, status factcheck.StatusTopic) (factcheck.Topic, error)
	UpdateDescription(ctx context.Context, id string, description string) (factcheck.Topic, error)
	UpdateName(ctx context.Context, id string, name string) (factcheck.Topic, error)
}

// topics implements RepositoryTopic
type topics struct {
	queries *postgres.Queries
}

// NewRepositoryTopic creates a new topic repository
func NewRepositoryTopic(queries *postgres.Queries) Topics {
	return &topics{
		queries: queries,
	}
}

// Create creates a new topic using the topic adapter
func (t *topics) Create(ctx context.Context, top factcheck.Topic) (factcheck.Topic, error) {
	params, err := topic(top)
	if err != nil {
		return factcheck.Topic{}, err
	}
	dbTopic, err := t.queries.CreateTopic(ctx, params)
	if err != nil {
		return factcheck.Topic{}, err
	}

	return topicDomain(dbTopic), nil
}

// GetByID retrieves a topic by ID using the topicDomain adapter
func (t *topics) GetByID(ctx context.Context, id string) (factcheck.Topic, error) {
	topicID, err := uuid(id)
	if err != nil {
		return factcheck.Topic{}, err
	}
	dbTopic, err := t.queries.GetTopic(ctx, topicID)
	if err != nil {
		return factcheck.Topic{}, handleNotFound(err, map[string]string{"id": id})
	}
	return topicDomain(dbTopic), nil
}

// List retrieves all topics using the topicDomain adapter
func (t *topics) List(ctx context.Context) ([]factcheck.Topic, error) {
	dbTopics, err := t.queries.ListTopics(ctx)
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
func (t *topics) ListByStatus(ctx context.Context, status factcheck.StatusTopic) ([]factcheck.Topic, error) {
	dbTopics, err := t.queries.ListTopicsByStatus(ctx, string(status))
	if err != nil {
		return nil, err
	}
	topics := make([]factcheck.Topic, len(dbTopics))
	for i, dbTopic := range dbTopics {
		topics[i] = topicDomain(dbTopic)
	}
	return topics, nil
}

func (t *topics) CountByStatus(ctx context.Context, status factcheck.StatusTopic) (int64, error) {
	return t.queries.CountTopicsByStatus(ctx, string(status))
}

func (t *topics) CountByStatuses(ctx context.Context) (map[factcheck.StatusTopic]int64, error) {
	rows, err := t.queries.CountTopicsGroupedByStatus(ctx)
	if err != nil {
		return nil, err
	}
	result := make(map[factcheck.StatusTopic]int64)
	for i := range rows {
		r := &rows[i]
		s := factcheck.StatusTopic(r.Status)
		if !s.IsValid() {
			return nil, fmt.Errorf("unexpected invalid status '%s' with %d count", s, r.Count)
		}
		result[s] = r.Count
	}
	return result, nil
}

// Delete deletes a topic by ID using the stringToUUID adapter
func (t *topics) Delete(ctx context.Context, id string) error {
	topicID, err := uuid(id)
	if err != nil {
		return err
	}
	err = t.queries.DeleteTopic(ctx, topicID)
	if err != nil {
		return handleNotFound(err, map[string]string{"id": id})
	}
	return nil
}

func (t *topics) UpdateStatus(ctx context.Context, id string, status factcheck.StatusTopic) (factcheck.Topic, error) {
	topicID, err := uuid(id)
	if err != nil {
		return factcheck.Topic{}, err
	}
	dbTopic, err := t.queries.UpdateTopicStatus(ctx, postgres.UpdateTopicStatusParams{
		ID:     topicID,
		Status: string(status),
	})
	if err != nil {
		return factcheck.Topic{}, handleNotFound(err, map[string]string{"id": id})
	}
	return topicDomain(dbTopic), nil
}

func (t *topics) UpdateDescription(ctx context.Context, id string, description string) (factcheck.Topic, error) {
	topicID, err := uuid(id)
	if err != nil {
		return factcheck.Topic{}, err
	}
	dbTopic, err := t.queries.UpdateTopicDescription(ctx, postgres.UpdateTopicDescriptionParams{
		ID:          topicID,
		Description: description,
	})
	if err != nil {
		return factcheck.Topic{}, handleNotFound(err, map[string]string{"id": id})
	}
	return topicDomain(dbTopic), nil
}

func (t *topics) UpdateName(ctx context.Context, id string, name string) (factcheck.Topic, error) {
	topicID, err := uuid(id)
	if err != nil {
		return factcheck.Topic{}, err
	}
	dbTopic, err := t.queries.UpdateTopicName(ctx, postgres.UpdateTopicNameParams{
		ID:   topicID,
		Name: name,
	})
	if err != nil {
		return factcheck.Topic{}, handleNotFound(err, map[string]string{"id": id})
	}
	return topicDomain(dbTopic), nil
}
