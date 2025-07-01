package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/postgres/sqlcgen"
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

func (r *repositoryTopic) Create(ctx context.Context, topic factcheck.Topic) (factcheck.Topic, error) {
	// Convert string ID to UUID
	var topicID pgtype.UUID
	if err := topicID.Scan(topic.ID); err != nil {
		return factcheck.Topic{}, err
	}

	// Convert timestamps
	createdAt := pgtype.Timestamptz{}
	if err := createdAt.Scan(topic.CreatedAt); err != nil {
		return factcheck.Topic{}, err
	}

	var updatedAt pgtype.Timestamptz
	if topic.UpdatedAt != nil {
		if err := updatedAt.Scan(*topic.UpdatedAt); err != nil {
			return factcheck.Topic{}, err
		}
	}

	// Convert optional fields
	var result pgtype.Text
	if topic.Result != "" {
		if err := result.Scan(topic.Result); err != nil {
			return factcheck.Topic{}, err
		}
	}

	var resultStatus pgtype.Text
	if topic.ResultStatus != "" {
		if err := resultStatus.Scan(string(topic.ResultStatus)); err != nil {
			return factcheck.Topic{}, err
		}
	}

	params := postgres.CreateTopicParams{
		ID:           topicID,
		Name:         topic.Name,
		Status:       string(topic.Status),
		Result:       result,
		ResultStatus: resultStatus,
		CreatedAt:    createdAt,
		UpdatedAt:    updatedAt,
	}

	dbTopic, err := r.queries.CreateTopic(ctx, params)
	if err != nil {
		return factcheck.Topic{}, err
	}

	return r.convertToDomainTopic(dbTopic), nil
}

func (r *repositoryTopic) GetByID(ctx context.Context, id string) (factcheck.Topic, error) {
	var topicID pgtype.UUID
	if err := topicID.Scan(id); err != nil {
		return factcheck.Topic{}, err
	}

	dbTopic, err := r.queries.GetTopic(ctx, topicID)
	if err != nil {
		return factcheck.Topic{}, err
	}

	return r.convertToDomainTopic(dbTopic), nil
}

func (r *repositoryTopic) List(ctx context.Context) ([]factcheck.Topic, error) {
	dbTopics, err := r.queries.ListTopics(ctx)
	if err != nil {
		return nil, err
	}

	topics := make([]factcheck.Topic, len(dbTopics))
	for i, dbTopic := range dbTopics {
		topics[i] = r.convertToDomainTopic(dbTopic)
	}

	return topics, nil
}

func (r *repositoryTopic) ListByStatus(ctx context.Context, status factcheck.StatusTopic) ([]factcheck.Topic, error) {
	dbTopics, err := r.queries.ListTopicsByStatus(ctx, string(status))
	if err != nil {
		return nil, err
	}

	topics := make([]factcheck.Topic, len(dbTopics))
	for i, dbTopic := range dbTopics {
		topics[i] = r.convertToDomainTopic(dbTopic)
	}

	return topics, nil
}

func (r *repositoryTopic) Update(ctx context.Context, topic factcheck.Topic) (factcheck.Topic, error) {
	// Convert string ID to UUID
	var topicID pgtype.UUID
	if err := topicID.Scan(topic.ID); err != nil {
		return factcheck.Topic{}, err
	}

	// Convert timestamps
	var updatedAt pgtype.Timestamptz
	if topic.UpdatedAt != nil {
		if err := updatedAt.Scan(*topic.UpdatedAt); err != nil {
			return factcheck.Topic{}, err
		}
	}

	// Convert optional fields
	var result pgtype.Text
	if topic.Result != "" {
		if err := result.Scan(topic.Result); err != nil {
			return factcheck.Topic{}, err
		}
	}

	var resultStatus pgtype.Text
	if topic.ResultStatus != "" {
		if err := resultStatus.Scan(string(topic.ResultStatus)); err != nil {
			return factcheck.Topic{}, err
		}
	}

	params := postgres.UpdateTopicParams{
		ID:           topicID,
		Name:         topic.Name,
		Status:       string(topic.Status),
		Result:       result,
		ResultStatus: resultStatus,
		UpdatedAt:    updatedAt,
	}

	dbTopic, err := r.queries.UpdateTopic(ctx, params)
	if err != nil {
		return factcheck.Topic{}, err
	}

	return r.convertToDomainTopic(dbTopic), nil
}

func (r *repositoryTopic) Delete(ctx context.Context, id string) error {
	var topicID pgtype.UUID
	if err := topicID.Scan(id); err != nil {
		return err
	}

	return r.queries.DeleteTopic(ctx, topicID)
}

// convertToDomainTopic converts a database topic to domain topic
func (r *repositoryTopic) convertToDomainTopic(dbTopic postgres.Topic) factcheck.Topic {
	topic := factcheck.Topic{
		Name:   dbTopic.Name,
		Status: factcheck.StatusTopic(dbTopic.Status),
	}

	// Convert UUID to string - using a simpler approach
	if dbTopic.ID.Valid {
		// Try to get the string representation directly
		topic.ID = dbTopic.ID.String()
	}

	// Convert optional text fields
	if dbTopic.Result.Valid {
		topic.Result = dbTopic.Result.String
	}

	if dbTopic.ResultStatus.Valid {
		topic.ResultStatus = factcheck.StatusTopicResult(dbTopic.ResultStatus.String)
	}

	// Convert timestamps
	if dbTopic.CreatedAt.Valid {
		topic.CreatedAt = dbTopic.CreatedAt.Time
	}

	if dbTopic.UpdatedAt.Valid {
		topic.UpdatedAt = &dbTopic.UpdatedAt.Time
	}

	return topic
}
