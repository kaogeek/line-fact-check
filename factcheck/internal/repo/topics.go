package repo

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgtype"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/data/postgres"
)

// Topics defines the interface for topic data operations
type Topics interface {
	Create(ctx context.Context, topic factcheck.Topic) (factcheck.Topic, error)
	GetByID(ctx context.Context, id string) (factcheck.Topic, error)
	List(ctx context.Context) ([]factcheck.Topic, error)
	ListByStatus(ctx context.Context, status factcheck.StatusTopic) ([]factcheck.Topic, error)
	ListFiltered(ctx context.Context, ids []string, messageText string) ([]factcheck.Topic, error)
	ListInIDs(ctx context.Context, ids []string) ([]factcheck.Topic, error)
	ListByMessageText(ctx context.Context, substring string) ([]factcheck.Topic, error)
	ListInIDsAndMessageText(ctx context.Context, ids []string, substring string) ([]factcheck.Topic, error)
	ListLikeID(ctx context.Context, idPattern string) ([]factcheck.Topic, error)
	ListLikeIDAndMessageText(ctx context.Context, idPattern string, messageText string) ([]factcheck.Topic, error)
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

// NewTopics creates a new topic repository
func NewTopics(queries *postgres.Queries) Topics {
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

// ListFiltered retrieves topics with optional filtering by IDs and/or message text
func (t *topics) ListFiltered(ctx context.Context, ids []string, messageText string) ([]factcheck.Topic, error) {
	if len(ids) > 0 && messageText != "" {
		return t.ListInIDsAndMessageText(ctx, ids, messageText)
	}
	if len(ids) > 0 {
		return t.ListInIDs(ctx, ids)
	}
	if messageText != "" {
		return t.ListByMessageText(ctx, messageText)
	}
	return t.List(ctx)
}

// ListInIDs retrieves topics by IDs using the topicDomain adapter
func (t *topics) ListInIDs(ctx context.Context, ids []string) ([]factcheck.Topic, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	uuidIDs := make([]pgtype.UUID, len(ids))
	for i, id := range ids {
		uuidID, err := uuid(id)
		if err != nil {
			return nil, err
		}
		uuidIDs[i] = uuidID
	}
	dbTopics, err := t.queries.ListTopicsInIDs(ctx, uuidIDs)
	if err != nil {
		return nil, err
	}
	topics := make([]factcheck.Topic, len(dbTopics))
	for i, dbTopic := range dbTopics {
		topics[i] = topicDomain(dbTopic)
	}
	return topics, nil
}

// ListByMessageText retrieves topics that have messages containing the given substring
func (t *topics) ListByMessageText(ctx context.Context, substring string) ([]factcheck.Topic, error) {
	likePattern := "%" + substring + "%"
	dbTopics, err := t.queries.ListTopicsByMessageText(ctx, likePattern)
	if err != nil {
		return nil, err
	}
	topics := make([]factcheck.Topic, len(dbTopics))
	for i, dbTopic := range dbTopics {
		topics[i] = topicDomain(dbTopic)
	}
	return topics, nil
}

// ListInIDsAndMessageText retrieves topics by IDs that also have messages containing the given substring
func (t *topics) ListInIDsAndMessageText(ctx context.Context, ids []string, substring string) ([]factcheck.Topic, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	uuidIDs := make([]pgtype.UUID, len(ids))
	for i, id := range ids {
		uuidID, err := uuid(id)
		if err != nil {
			return nil, err
		}
		uuidIDs[i] = uuidID
	}
	likePattern := "%" + substring + "%"
	dbTopics, err := t.queries.ListTopicsInIDsAndMessageText(ctx, postgres.ListTopicsInIDsAndMessageTextParams{
		Column1: uuidIDs,
		Text:    likePattern,
	})
	if err != nil {
		return nil, err
	}
	topics := make([]factcheck.Topic, len(dbTopics))
	for i, dbTopic := range dbTopics {
		topics[i] = topicDomain(dbTopic)
	}
	return topics, nil
}

// ListLikeID retrieves topics by ID pattern matching using SQL LIKE
func (t *topics) ListLikeID(ctx context.Context, idPattern string) ([]factcheck.Topic, error) {
	// Add wildcards for LIKE query if not already present
	likePattern := idPattern
	if !strings.Contains(likePattern, "%") {
		likePattern = "%" + idPattern + "%"
	}

	dbTopics, err := t.queries.ListTopicsLikeID(ctx, likePattern)
	if err != nil {
		return nil, err
	}
	topics := make([]factcheck.Topic, len(dbTopics))
	for i, dbTopic := range dbTopics {
		topics[i] = topicDomain(dbTopic)
	}
	return topics, nil
}

// ListLikeIDAndMessageText retrieves topics by ID pattern and message text using SQL LIKE
func (t *topics) ListLikeIDAndMessageText(ctx context.Context, idPattern string, messageText string) ([]factcheck.Topic, error) {
	// Add wildcards for LIKE queries if not already present
	idLikePattern := idPattern
	if !strings.Contains(idLikePattern, "%") {
		idLikePattern = "%" + idPattern + "%"
	}

	messageLikePattern := "%" + messageText + "%"

	dbTopics, err := t.queries.ListTopicsLikeIDAndMessageText(ctx, postgres.ListTopicsLikeIDAndMessageTextParams{
		Column1: idLikePattern,
		Text:    messageLikePattern,
	})
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
