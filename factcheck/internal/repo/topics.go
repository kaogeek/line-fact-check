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
	ListHomePage(ctx context.Context, opts ...OptionListTopicHome) ([]factcheck.Topic, error)
	ListByStatus(ctx context.Context, status factcheck.StatusTopic) ([]factcheck.Topic, error)
	ListInIDs(ctx context.Context, ids []string) ([]factcheck.Topic, error)
	ListLikeMessageText(ctx context.Context, pattern string) ([]factcheck.Topic, error)
	ListLikeID(ctx context.Context, idPattern string) ([]factcheck.Topic, error)
	ListLikeIDLikeMessageText(ctx context.Context, idPattern string, pattern string) ([]factcheck.Topic, error)
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

type OptionListTopicHome func(f FilterListTopicsHome) FilterListTopicsHome

type FilterListTopicsHome struct {
	likeID          string
	likeMessageText string
	status          factcheck.StatusTopic
}

func LikeTopicID(id string) OptionListTopicHome {
	return func(f FilterListTopicsHome) FilterListTopicsHome {
		f.likeID = id
		return f
	}
}

func LikeTopicMessageText(s string) OptionListTopicHome {
	return func(f FilterListTopicsHome) FilterListTopicsHome {
		f.likeMessageText = s
		return f
	}
}

func WithTopicStatus(s factcheck.StatusTopic) OptionListTopicHome {
	return func(f FilterListTopicsHome) FilterListTopicsHome {
		f.status = s
		return f
	}
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

func empty[S ~string](s S) bool {
	return s == ""
}

func (t *topics) ListHomePage(ctx context.Context, opts ...OptionListTopicHome) ([]factcheck.Topic, error) {
	f := FilterListTopicsHome{}
	for i := range opts {
		f = opts[i](f)
	}

	switch {
	case empty(f.likeID) && empty(f.likeMessageText) && empty(f.status):
		return t.List(ctx)

	case empty(f.likeID) && empty(f.likeMessageText):
		return t.ListByStatus(ctx, f.status)

	case empty(f.likeID) && empty(f.status):
		return t.ListLikeMessageText(ctx, f.likeMessageText)

	case empty(f.likeMessageText) && empty(f.status):
		return t.ListLikeID(ctx, f.likeID)

	case empty(f.likeID):
		// Status + message text filter
		likePattern := substring(f.likeMessageText)
		dbTopics, err := t.queries.ListTopicsByStatusAndMessageText(ctx, postgres.ListTopicsByStatusAndMessageTextParams{
			Status: string(f.status),
			Text:   likePattern,
		})
		if err != nil {
			return nil, err
		}
		topics := make([]factcheck.Topic, len(dbTopics))
		for i, dbTopic := range dbTopics {
			topics[i] = topicDomain(dbTopic)
		}
		return topics, nil

	case empty(f.likeMessageText):
		// Status + ID pattern filter
		idPattern := f.likeID
		if !strings.Contains(idPattern, "%") {
			idPattern = substring(idPattern)
		}
		dbTopics, err := t.queries.ListTopicsByStatusAndLikeID(ctx, postgres.ListTopicsByStatusAndLikeIDParams{
			Status:  string(f.status),
			Column2: idPattern,
		})
		if err != nil {
			return nil, err
		}
		topics := make([]factcheck.Topic, len(dbTopics))
		for i, dbTopic := range dbTopics {
			topics[i] = topicDomain(dbTopic)
		}
		return topics, nil

	case empty(f.status):
		// ID pattern + message text filter
		return t.ListLikeIDLikeMessageText(ctx, f.likeID, f.likeMessageText)

	default:
		// All three filters
		idPattern := f.likeID
		if !strings.Contains(idPattern, "%") {
			idPattern = substring(idPattern)
		}
		messageLikePattern := substring(f.likeMessageText)
		dbTopics, err := t.queries.ListTopicsByStatusAndLikeIDAndMessageText(ctx, postgres.ListTopicsByStatusAndLikeIDAndMessageTextParams{
			Status:  string(f.status),
			Column2: idPattern,
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

// ListLikeMessageText retrieves topics that have messages containing the given substring
func (t *topics) ListLikeMessageText(ctx context.Context, pattern string) ([]factcheck.Topic, error) {
	likePattern := substring(pattern)
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

// ListLikeID retrieves topics by ID pattern matching using SQL LIKE
func (t *topics) ListLikeID(ctx context.Context, idPattern string) ([]factcheck.Topic, error) {
	// Add wildcards for LIKE query if not already present
	likePattern := idPattern
	if !strings.Contains(likePattern, "%") {
		likePattern = substring(idPattern)
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

// ListLikeIDLikeMessageText retrieves topics by ID pattern and message text using SQL LIKE
func (t *topics) ListLikeIDLikeMessageText(ctx context.Context, idPattern string, pattern string) ([]factcheck.Topic, error) {
	// Add wildcards for LIKE queries if not already present
	idLikePattern := idPattern
	if !strings.Contains(idLikePattern, "%") {
		idLikePattern = substring(idPattern)
	}

	messageLikePattern := substring(pattern)

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
