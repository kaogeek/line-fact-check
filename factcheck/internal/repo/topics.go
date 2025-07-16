package repo

import (
	"context"
	"fmt"
	"strings"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/data/postgres"
)

// Topics defines the interface for topic data operations
type Topics interface {
	Create(ctx context.Context, topic factcheck.Topic, opts ...OptionTx) (factcheck.Topic, error)
	GetByID(ctx context.Context, id string) (factcheck.Topic, error)
	List(ctx context.Context, limit, offset int) ([]factcheck.Topic, error)
	ListHome(ctx context.Context, limit, offset int, opts ...OptionTopic) ([]factcheck.Topic, error)
	ListByStatus(ctx context.Context, status factcheck.StatusTopic, limit, offset int) ([]factcheck.Topic, error)
	ListInIDs(ctx context.Context, ids []string) ([]factcheck.Topic, error)
	ListLikeMessageText(ctx context.Context, pattern string, limit, offset int) ([]factcheck.Topic, error)
	ListLikeID(ctx context.Context, idPattern string, limit, offset int) ([]factcheck.Topic, error)
	ListLikeIDLikeMessageText(ctx context.Context, idPattern string, pattern string, limit, offset int) ([]factcheck.Topic, error)
	ListLikeIDLikeMessageTextAll(ctx context.Context, idPattern string, pattern string) ([]factcheck.Topic, error) // Backward compatibility
	CountStatus(ctx context.Context, status factcheck.StatusTopic) (int64, error)
	CountByStatus(ctx context.Context) (map[factcheck.StatusTopic]int64, error)
	CountByStatusHome(ctx context.Context, opts ...OptionTopic) (map[factcheck.StatusTopic]int64, error)
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

// List retrieves topics with pagination using the topicDomain adapter
func (t *topics) List(ctx context.Context, limit, offset int) ([]factcheck.Topic, error) {
	dbTopics, err := t.queries.ListTopics(ctx, postgres.ListTopicsParams{
		Column1: limit,
		Column2: offset,
	})
	if err != nil {
		return nil, err
	}
	topics := make([]factcheck.Topic, len(dbTopics))
	for i, dbTopic := range dbTopics {
		topics[i] = postgres.ToTopicFromRow(dbTopic)
	}
	return topics, nil
}

// ListAll retrieves all topics (backward compatibility)
func (t *topics) ListAll(ctx context.Context) ([]factcheck.Topic, error) {
	return t.List(ctx, 0, 0)
}

// ListByStatus retrieves topics by status with pagination
func (t *topics) ListByStatus(ctx context.Context, status factcheck.StatusTopic, limit, offset int) ([]factcheck.Topic, error) {
	limit, offset = sanitize(limit, offset)
	dbTopics, err := t.queries.ListTopicsByStatus(ctx, postgres.ListTopicsByStatusParams{
		Status:  string(status),
		Column2: limit,
		Column3: offset,
	})
	if err != nil {
		return nil, err
	}
	topics := make([]factcheck.Topic, len(dbTopics))
	for i, dbTopic := range dbTopics {
		topics[i] = postgres.ToTopicFromStatusRow(dbTopic)
	}
	return topics, nil
}

// ListLikeMessageText retrieves topics that have messages containing the given substring with pagination
func (t *topics) ListLikeMessageText(ctx context.Context, pattern string, limit, offset int) ([]factcheck.Topic, error) {
	limit, offset = sanitize(limit, offset)
	likePattern := substring(pattern)
	dbTopics, err := t.queries.ListTopicsLikeMessageText(ctx, postgres.ListTopicsLikeMessageTextParams{
		Text:    likePattern,
		Column2: limit,
		Column3: offset,
	})
	if err != nil {
		return nil, err
	}
	topics := make([]factcheck.Topic, len(dbTopics))
	for i, dbTopic := range dbTopics {
		topics[i] = postgres.ToTopicFromMessageTextRow(dbTopic)
	}
	return topics, nil
}

// ListLikeID retrieves topics by ID pattern matching using SQL LIKE with pagination
func (t *topics) ListLikeID(ctx context.Context, idPattern string, limit, offset int) ([]factcheck.Topic, error) {
	limit, offset = sanitize(limit, offset)
	// Add wildcards for LIKE query if not already present
	likePattern := idPattern
	if !strings.Contains(likePattern, "%") {
		likePattern = substring(idPattern)
	}

	dbTopics, err := t.queries.ListTopicsLikeID(ctx, postgres.ListTopicsLikeIDParams{
		Column1: likePattern,
		Column2: limit,
		Column3: offset,
	})
	if err != nil {
		return nil, err
	}
	topics := make([]factcheck.Topic, len(dbTopics))
	for i, dbTopic := range dbTopics {
		topics[i] = postgres.ToTopicFromIDRow(dbTopic)
	}
	return topics, nil
}

// TODO: handle subquery transactions
func (t *topics) ListHome(ctx context.Context, limit, offset int, opts ...OptionTopic) ([]factcheck.Topic, error) {
	limit, offset = sanitize(limit, offset)

	options := &OptionsTopic{}
	for _, opt := range opts {
		opt(options)
	}

	switch {
	case empty(options.LikeID) && empty(options.LikeMessageText) && empty(options.Status):
		return t.List(ctx, limit, offset)

	case empty(options.LikeID) && empty(options.LikeMessageText):
		return t.ListByStatus(ctx, options.Status, limit, offset)

	case empty(options.LikeID) && empty(options.Status):
		return t.ListLikeMessageText(ctx, options.LikeMessageText, limit, offset)

	case empty(options.LikeMessageText) && empty(options.Status):
		return t.ListLikeID(ctx, options.LikeID, limit, offset)

	case empty(options.LikeID):
		// Status + message text filter
		likePattern := substring(options.LikeMessageText)
		topics, err := t.queries.ListTopicsByStatusLikeMessageText(ctx, postgres.ListTopicsByStatusLikeMessageTextParams{
			Status:  string(options.Status),
			Text:    likePattern,
			Column3: limit,
			Column4: offset,
		})
		if err != nil {
			return nil, err
		}
		result := make([]factcheck.Topic, len(topics))
		for i, dbTopic := range topics {
			result[i] = postgres.ToTopicFromStatusLikeMessageTextRow(dbTopic)
		}
		return result, nil

	case empty(options.LikeMessageText):
		// Status + ID pattern filter
		idPattern := options.LikeID
		if !strings.Contains(idPattern, "%") {
			idPattern = substring(idPattern)
		}
		topics, err := t.queries.ListTopicsByStatusLikeID(ctx, postgres.ListTopicsByStatusLikeIDParams{
			Status:  string(options.Status),
			Column2: idPattern,
			Column3: limit,
			Column4: offset,
		})
		if err != nil {
			return nil, err
		}
		result := make([]factcheck.Topic, len(topics))
		for i, dbTopic := range topics {
			result[i] = postgres.ToTopicFromStatusLikeIDRow(dbTopic)
		}
		return result, nil

	case empty(options.Status):
		// ID pattern + message text filter
		return t.ListLikeIDLikeMessageText(ctx, options.LikeID, options.LikeMessageText, limit, offset)
	}

	// All three filters
	idPattern := options.LikeID
	if !strings.Contains(idPattern, "%") {
		idPattern = substring(idPattern)
	}
	messageLikePattern := substring(options.LikeMessageText)
	topics, err := t.queries.ListTopicsByStatusLikeIDLikeMessageText(ctx, postgres.ListTopicsByStatusLikeIDLikeMessageTextParams{
		Status:  string(options.Status),
		Column2: idPattern,
		Text:    messageLikePattern,
		Column4: limit,
		Column5: offset,
	})
	if err != nil {
		return nil, err
	}
	result := make([]factcheck.Topic, len(topics))
	for i, dbTopic := range topics {
		result[i] = postgres.ToTopicFromStatusLikeIDLikeMessageTextRow(dbTopic)
	}
	return result, nil
}

func (t *topics) CountByStatusHome(ctx context.Context, opts ...OptionTopic) (map[factcheck.StatusTopic]int64, error) {
	options := &OptionsTopic{}
	for _, opt := range opts {
		opt(options)
	}

	switch {
	case empty(options.LikeID) && empty(options.LikeMessageText):
		return t.CountByStatus(ctx)

	case empty(options.LikeID):
		likePattern := substring(options.LikeMessageText)
		result, err := t.queries.CountTopicsGroupByStatusLikeMessageText(ctx, likePattern)
		if err != nil {
			return nil, err
		}
		m := make(map[factcheck.StatusTopic]int64)
		for i := range result {
			m[factcheck.StatusTopic(result[i].Status)] = result[i].Count
		}
		return m, nil

	case empty(options.LikeMessageText):
		idPattern := options.LikeID
		if !strings.Contains(idPattern, "%") {
			idPattern = substring(idPattern)
		}
		result, err := t.queries.CountTopicsGroupByStatusLikeID(ctx, idPattern)
		if err != nil {
			return nil, err
		}
		m := make(map[factcheck.StatusTopic]int64)
		for i := range result {
			m[factcheck.StatusTopic(result[i].Status)] = result[i].Count
		}
		return m, nil
	}
	idPattern := options.LikeID
	if !strings.Contains(idPattern, "%") {
		idPattern = substring(idPattern)
	}
	messageLikePattern := substring(options.LikeMessageText)
	result, err := t.queries.CountTopicsGroupByStatusLikeIDLikeMessageText(ctx, postgres.CountTopicsGroupByStatusLikeIDLikeMessageTextParams{
		Column1: idPattern,
		Text:    messageLikePattern,
	})
	if err != nil {
		return nil, err
	}
	m := make(map[factcheck.StatusTopic]int64)
	for i := range result {
		m[factcheck.StatusTopic(result[i].Status)] = result[i].Count
	}
	return m, nil
}

// Create creates a new topic using the topic adapter
func (t *topics) Create(ctx context.Context, top factcheck.Topic, opts ...OptionTx) (factcheck.Topic, error) {
	txOptions := &TxOptions{}
	for _, opt := range opts {
		opt(txOptions)
	}

	params, err := postgres.TopicCreator(top)
	if err != nil {
		return factcheck.Topic{}, err
	}

	query := t.queries.CreateTopic
	if txOptions.Tx != nil {
		query = t.queries.WithTx(txOptions.Tx).CreateTopic
	}

	dbTopic, err := query(ctx, params)
	if err != nil {
		return factcheck.Topic{}, err
	}

	return postgres.ToTopic(dbTopic), nil
}

// GetByID retrieves a topic by ID using the topicDomain adapter
func (t *topics) GetByID(ctx context.Context, id string) (factcheck.Topic, error) {
	topicID, err := postgres.UUID(id)
	if err != nil {
		return factcheck.Topic{}, err
	}
	dbTopic, err := t.queries.GetTopic(ctx, topicID)
	if err != nil {
		return factcheck.Topic{}, handleNotFound(err, map[string]string{"id": id})
	}
	return postgres.ToTopic(dbTopic), nil
}

// ListInIDs retrieves topics by IDs using the topicDomain adapter
func (t *topics) ListInIDs(ctx context.Context, ids []string) ([]factcheck.Topic, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	uuids, err := postgres.UUIDs(ids)
	if err != nil {
		return nil, err
	}
	dbTopics, err := t.queries.ListTopicsInIDs(ctx, uuids)
	if err != nil {
		return nil, err
	}
	topics := make([]factcheck.Topic, len(dbTopics))
	for i, dbTopic := range dbTopics {
		topics[i] = postgres.ToTopic(dbTopic)
	}
	return topics, nil
}

// ListLikeIDLikeMessageText retrieves topics by ID pattern and message text using SQL LIKE with pagination
func (t *topics) ListLikeIDLikeMessageText(ctx context.Context, idPattern string, pattern string, limit, offset int) ([]factcheck.Topic, error) {
	limit, offset = sanitize(limit, offset)
	// Add wildcards for LIKE queries if not already present
	idLikePattern := idPattern
	if !strings.Contains(idLikePattern, "%") {
		idLikePattern = substring(idPattern)
	}

	messageLikePattern := substring(pattern)

	dbTopics, err := t.queries.ListTopicsLikeIDLikeMessageText(ctx, postgres.ListTopicsLikeIDLikeMessageTextParams{
		Column1: idLikePattern,
		Text:    messageLikePattern,
		Column3: limit,
		Column4: offset,
	})
	if err != nil {
		return nil, err
	}
	topics := make([]factcheck.Topic, len(dbTopics))
	for i, dbTopic := range dbTopics {
		topics[i] = postgres.ToTopicFromIDLikeMessageTextRow(dbTopic)
	}
	return topics, nil
}

// ListLikeIDLikeMessageTextAll retrieves all topics by ID pattern and message text using SQL LIKE (backward compatibility)
func (t *topics) ListLikeIDLikeMessageTextAll(ctx context.Context, idPattern string, pattern string) ([]factcheck.Topic, error) {
	return t.ListLikeIDLikeMessageText(ctx, idPattern, pattern, 0, 0)
}

func (t *topics) CountStatus(ctx context.Context, status factcheck.StatusTopic) (int64, error) {
	return t.queries.CountTopicsByStatus(ctx, string(status))
}

func (t *topics) CountByStatus(ctx context.Context) (map[factcheck.StatusTopic]int64, error) {
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
	topicID, err := postgres.UUID(id)
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
	topicID, err := postgres.UUID(id)
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
	return postgres.ToTopic(dbTopic), nil
}

func (t *topics) UpdateDescription(ctx context.Context, id string, description string) (factcheck.Topic, error) {
	topicID, err := postgres.UUID(id)
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
	return postgres.ToTopic(dbTopic), nil
}

func (t *topics) UpdateName(ctx context.Context, id string, name string) (factcheck.Topic, error) {
	topicID, err := postgres.UUID(id)
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
	return postgres.ToTopic(dbTopic), nil
}

func empty[S ~string](s S) bool {
	return s == ""
}
