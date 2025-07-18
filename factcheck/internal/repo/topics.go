package repo

import (
	"context"
	"fmt"
	"strings"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/data/postgres"
	"github.com/kaogeek/line-fact-check/factcheck/internal/utils"
)

// Topics defines the interface for topic data operations
type Topics interface {
	Create(ctx context.Context, topic factcheck.Topic, opts ...Option) (factcheck.Topic, error)
	GetByID(ctx context.Context, id string, opts ...Option) (factcheck.Topic, error)
	List(ctx context.Context, limit, offset int, opts ...Option) ([]factcheck.Topic, error)
	ListDynamic(ctx context.Context, limit, offset int, opts ...OptionTopicDynamic) ([]factcheck.Topic, error)
	ListInIDs(ctx context.Context, ids []string, opts ...Option) ([]factcheck.Topic, error)
	ListByStatus(ctx context.Context, status factcheck.StatusTopic, limit, offset int, opts ...Option) ([]factcheck.Topic, error)
	CountByStatus(ctx context.Context, opts ...Option) (map[factcheck.StatusTopic]int64, error)
	CountByStatusHome(ctx context.Context, opts ...OptionTopic) (map[factcheck.StatusTopic]int64, error)
	Delete(ctx context.Context, id string, opts ...Option) error
	UpdateStatus(ctx context.Context, id string, status factcheck.StatusTopic, opts ...Option) (factcheck.Topic, error)
	UpdateDescription(ctx context.Context, id string, description string, opts ...Option) (factcheck.Topic, error)
	UpdateName(ctx context.Context, id string, name string, opts ...Option) (factcheck.Topic, error)

	// TODO: deprecate, replace with ListDynamic
	ListHome(ctx context.Context, limit, offset int, opts ...OptionTopic) ([]factcheck.Topic, error)
	ListLikeMessageText(ctx context.Context, pattern string, limit, offset int, opts ...Option) ([]factcheck.Topic, error)
	ListLikeID(ctx context.Context, idPattern string, limit, offset int, opts ...Option) ([]factcheck.Topic, error)
	ListLikeIDLikeMessageText(ctx context.Context, idPattern string, pattern string, limit, offset int, opts ...Option) ([]factcheck.Topic, error)
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
func (t *topics) List(ctx context.Context, limit, offset int, opts ...Option) ([]factcheck.Topic, error) {
	queries := queries(t.queries, options(opts...))
	rows, err := queries.ListTopics(ctx, postgres.ListTopicsParams{
		Column1: limit,
		Column2: offset,
	})
	if err != nil {
		return nil, err
	}
	return utils.MapSliceNoError(rows, postgres.ToTopicFromRow), nil
}

// ListAll retrieves all topics (backward compatibility)
func (t *topics) ListAll(ctx context.Context) ([]factcheck.Topic, error) {
	return t.List(ctx, 0, 0)
}

func (t *topics) ListDynamic(ctx context.Context, limit, offset int, opts ...OptionTopicDynamic) ([]factcheck.Topic, error) {
	limit, offset = sanitizeV2(limit, offset)
	options := options(opts...)
	queries := queries(t.queries, options.Options)
	rows, err := queries.ListTopicsDynamic(ctx, options.ListDynamicParams((offset), (limit)))
	if err != nil {
		return nil, err
	}
	return utils.MapSliceNoError(rows, postgres.ToTopic), nil
}

// ListByStatus retrieves topics by status with pagination
func (t *topics) ListByStatus(ctx context.Context, status factcheck.StatusTopic, limit, offset int, opts ...Option) ([]factcheck.Topic, error) {
	limit, offset = sanitize(limit, offset)
	queries := queries(t.queries, options(opts...))
	rows, err := queries.ListTopicsByStatus(ctx, postgres.ListTopicsByStatusParams{
		Status:  string(status),
		Column2: limit,
		Column3: offset,
	})
	if err != nil {
		return nil, err
	}
	return utils.MapSliceNoError(rows, postgres.ToTopicFromStatusRow), nil
}

// ListLikeMessageText retrieves topics that have messages containing the given substring with pagination
func (t *topics) ListLikeMessageText(ctx context.Context, pattern string, limit, offset int, opts ...Option) ([]factcheck.Topic, error) {
	limit, offset = sanitize(limit, offset)
	queries := queries(t.queries, options(opts...))
	likePattern := substring(pattern)
	rows, err := queries.ListTopicsLikeMessageText(ctx, postgres.ListTopicsLikeMessageTextParams{
		Text:    likePattern,
		Column2: limit,
		Column3: offset,
	})
	if err != nil {
		return nil, err
	}
	return utils.MapSliceNoError(rows, postgres.ToTopicFromMessageTextRow), nil
}

// ListLikeID retrieves topics by ID pattern matching using SQL LIKE with pagination
func (t *topics) ListLikeID(ctx context.Context, idPattern string, limit, offset int, opts ...Option) ([]factcheck.Topic, error) {
	limit, offset = sanitize(limit, offset)
	queries := queries(t.queries, options(opts...))
	// Add wildcards for LIKE query if not already present
	likePattern := idPattern
	if !strings.Contains(likePattern, "%") {
		likePattern = substring(idPattern)
	}
	rows, err := queries.ListTopicsLikeID(ctx, postgres.ListTopicsLikeIDParams{
		Column1: likePattern,
		Column2: limit,
		Column3: offset,
	})
	if err != nil {
		return nil, err
	}
	return utils.MapSliceNoError(rows, postgres.ToTopicFromIDRow), nil
}

func (t *topics) ListHome(
	ctx context.Context,
	limit, offset int,
	opts ...OptionTopic,
) (
	[]factcheck.Topic,
	error,
) {
	options := options(opts...)
	queries := queries(t.queries, options.Options)
	limit, offset = sanitize(limit, offset)

	switch {
	case empty(options.LikeID) && empty(options.LikeMessageText) && empty(options.Status):
		return t.List(ctx, limit, offset, options.Clone()...)

	case empty(options.LikeID) && empty(options.LikeMessageText):
		return t.ListByStatus(ctx, options.Status, limit, offset, options.Clone()...)

	case empty(options.LikeID) && empty(options.Status):
		return t.ListLikeMessageText(ctx, options.LikeMessageText, limit, offset, options.Clone()...)

	case empty(options.LikeMessageText) && empty(options.Status):
		return t.ListLikeID(ctx, options.LikeID, limit, offset, options.Clone()...)

	case empty(options.Status):
		// ID pattern + message text filter
		return t.ListLikeIDLikeMessageText(ctx, options.LikeID, options.LikeMessageText, limit, offset, options.Clone()...)

	case empty(options.LikeID):
		// Status + message text filter
		likePattern := substring(options.LikeMessageText)
		rows, err := queries.ListTopicsByStatusLikeMessageText(ctx, postgres.ListTopicsByStatusLikeMessageTextParams{
			Status:  string(options.Status),
			Text:    likePattern,
			Column3: limit,
			Column4: offset,
		})
		if err != nil {
			return nil, err
		}
		return utils.MapSliceNoError(rows, postgres.ToTopicFromStatusLikeMessageTextRow), nil

	case empty(options.LikeMessageText):
		// Status + ID pattern filter
		idPattern := options.LikeID
		if !strings.Contains(idPattern, "%") {
			idPattern = substring(idPattern)
		}
		rows, err := queries.ListTopicsByStatusLikeID(ctx, postgres.ListTopicsByStatusLikeIDParams{
			Status:  string(options.Status),
			Column2: idPattern,
			Column3: limit,
			Column4: offset,
		})
		if err != nil {
			return nil, err
		}
		return utils.MapSliceNoError(rows, postgres.ToTopicFromStatusLikeIDRow), nil
	}

	// All three filters
	idPattern := options.LikeID
	if !strings.Contains(idPattern, "%") {
		idPattern = substring(idPattern)
	}
	messageLikePattern := substring(options.LikeMessageText)
	rows, err := queries.ListTopicsByStatusLikeIDLikeMessageText(ctx, postgres.ListTopicsByStatusLikeIDLikeMessageTextParams{
		Status:  string(options.Status),
		Column2: idPattern,
		Text:    messageLikePattern,
		Column4: limit,
		Column5: offset,
	})
	if err != nil {
		return nil, err
	}
	return utils.MapSliceNoError(rows, postgres.ToTopicFromStatusLikeIDLikeMessageTextRow), nil
}

func (t *topics) CountByStatusHome(ctx context.Context, opts ...OptionTopic) (map[factcheck.StatusTopic]int64, error) {
	options := options(opts...)
	queries := queries(t.queries, options.Options)

	switch {
	case empty(options.LikeID) && empty(options.LikeMessageText):
		return t.CountByStatus(ctx)

	case empty(options.LikeID):
		likePattern := substring(options.LikeMessageText)
		result, err := queries.CountTopicsGroupByStatusLikeMessageText(ctx, likePattern)
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
		result, err := queries.CountTopicsGroupByStatusLikeID(ctx, idPattern)
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
	result, err := queries.CountTopicsGroupByStatusLikeIDLikeMessageText(ctx, postgres.CountTopicsGroupByStatusLikeIDLikeMessageTextParams{
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
func (t *topics) Create(ctx context.Context, top factcheck.Topic, opts ...Option) (factcheck.Topic, error) {
	queries := queries(t.queries, options(opts...))
	params, err := postgres.TopicCreator(top)
	if err != nil {
		return factcheck.Topic{}, err
	}
	created, err := queries.CreateTopic(ctx, params)
	if err != nil {
		return factcheck.Topic{}, err
	}
	return postgres.ToTopic(created), nil
}

// GetByID retrieves a topic by ID using the topicDomain adapter
func (t *topics) GetByID(ctx context.Context, id string, opts ...Option) (factcheck.Topic, error) {
	queries := queries(t.queries, options(opts...))
	uuid, err := postgres.UUID(id)
	if err != nil {
		return factcheck.Topic{}, err
	}
	result, err := queries.GetTopic(ctx, uuid)
	if err != nil {
		return factcheck.Topic{}, handleNotFound(err, map[string]string{"id": id})
	}
	return postgres.ToTopic(result), nil
}

// ListInIDs retrieves topics by IDs using the topicDomain adapter
func (t *topics) ListInIDs(ctx context.Context, ids []string, opts ...Option) ([]factcheck.Topic, error) {
	queries := queries(t.queries, options(opts...))
	if len(ids) == 0 {
		return nil, nil
	}
	uuids, err := postgres.UUIDs(ids)
	if err != nil {
		return nil, err
	}
	rows, err := queries.ListTopicsInIDs(ctx, uuids)
	if err != nil {
		return nil, err
	}
	return utils.MapSliceNoError(rows, postgres.ToTopic), nil
}

// ListLikeIDLikeMessageText retrieves topics by ID pattern and message text using SQL LIKE with pagination
func (t *topics) ListLikeIDLikeMessageText(ctx context.Context, idPattern string, pattern string, limit, offset int, opts ...Option) ([]factcheck.Topic, error) {
	limit, offset = sanitize(limit, offset)
	queries := queries(t.queries, options(opts...))
	// Add wildcards for LIKE queries if not already present
	idLikePattern := idPattern
	if !strings.Contains(idLikePattern, "%") {
		idLikePattern = substring(idPattern)
	}
	messageLikePattern := substring(pattern)
	rows, err := queries.ListTopicsLikeIDLikeMessageText(ctx, postgres.ListTopicsLikeIDLikeMessageTextParams{
		Column1: idLikePattern,
		Text:    messageLikePattern,
		Column3: limit,
		Column4: offset,
	})
	if err != nil {
		return nil, err
	}
	return utils.MapSliceNoError(rows, postgres.ToTopicFromIDLikeMessageTextRow), nil
}

// ListLikeIDLikeMessageTextAll retrieves all topics by ID pattern and message text using SQL LIKE (backward compatibility)
func (t *topics) ListLikeIDLikeMessageTextAll(ctx context.Context, idPattern string, pattern string) ([]factcheck.Topic, error) {
	return t.ListLikeIDLikeMessageText(ctx, idPattern, pattern, 0, 0)
}

func (t *topics) CountByStatus(ctx context.Context, opts ...Option) (map[factcheck.StatusTopic]int64, error) {
	queries := queries(t.queries, options(opts...))
	rows, err := queries.CountTopicsGroupedByStatus(ctx)
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
func (t *topics) Delete(ctx context.Context, id string, opts ...Option) error {
	options := options(opts...)
	queries := queries(t.queries, options)
	uuid, err := postgres.UUID(id)
	if err != nil {
		return err
	}
	err = queries.DeleteTopic(ctx, uuid)
	if err != nil {
		return handleNotFound(err, map[string]string{"id": id})
	}
	return nil
}

func (t *topics) UpdateStatus(ctx context.Context, id string, status factcheck.StatusTopic, opts ...Option) (factcheck.Topic, error) {
	queries := queries(t.queries, options(opts...))
	uuid, err := postgres.UUID(id)
	if err != nil {
		return factcheck.Topic{}, err
	}
	dbTopic, err := queries.UpdateTopicStatus(ctx, postgres.UpdateTopicStatusParams{
		ID:     uuid,
		Status: string(status),
	})
	if err != nil {
		return factcheck.Topic{}, handleNotFound(err, map[string]string{"id": id})
	}
	return postgres.ToTopic(dbTopic), nil
}

func (t *topics) UpdateDescription(ctx context.Context, id string, description string, opts ...Option) (factcheck.Topic, error) {
	queries := queries(t.queries, options(opts...))
	uuid, err := postgres.UUID(id)
	if err != nil {
		return factcheck.Topic{}, err
	}
	updated, err := queries.UpdateTopicDescription(ctx, postgres.UpdateTopicDescriptionParams{
		ID:          uuid,
		Description: description,
	})
	if err != nil {
		return factcheck.Topic{}, handleNotFound(err, map[string]string{"id": id})
	}
	return postgres.ToTopic(updated), nil
}

func (t *topics) UpdateName(ctx context.Context, id string, name string, opts ...Option) (factcheck.Topic, error) {
	queries := queries(t.queries, options(opts...))
	uuid, err := postgres.UUID(id)
	if err != nil {
		return factcheck.Topic{}, err
	}
	updated, err := queries.UpdateTopicName(ctx, postgres.UpdateTopicNameParams{
		ID:   uuid,
		Name: name,
	})
	if err != nil {
		return factcheck.Topic{}, handleNotFound(err, map[string]string{"id": id})
	}
	return postgres.ToTopic(updated), nil
}

func empty[S ~string](s S) bool {
	return s == ""
}
