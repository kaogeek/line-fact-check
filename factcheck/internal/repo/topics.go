package repo

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/internal/data/postgres"
	"github.com/kaogeek/line-fact-check/factcheck/internal/utils"
)

// Topics defines the interface for topic data operations
type Topics interface {
	Create(ctx context.Context, topic factcheck.Topic, opts ...Option) (factcheck.Topic, error)
	Resolve(ctx context.Context, id string, answerText string, opts ...Option) (factcheck.Topic, error)
	GetByID(ctx context.Context, id string, opts ...Option) (factcheck.Topic, error)
	GetStatus(ctx context.Context, id string, opts ...Option) (factcheck.StatusTopic, error)
	Exists(ctx context.Context, id string, opts ...Option) (bool, error)
	List(ctx context.Context, limit, offset int, opts ...Option) ([]factcheck.Topic, error)
	ListDynamicV2(ctx context.Context, limit, offset int, opts ...OptionTopic) ([]factcheck.Topic, error)
	ListInIDs(ctx context.Context, ids []string, opts ...Option) ([]factcheck.Topic, error)
	ListByStatus(ctx context.Context, status factcheck.StatusTopic, limit, offset int, opts ...Option) ([]factcheck.Topic, error)
	CountByStatus(ctx context.Context, opts ...Option) (map[factcheck.StatusTopic]int64, error)
	CountByStatusDynamicV2(ctx context.Context, opts ...OptionTopic) (map[factcheck.StatusTopic]int64, error)
	Delete(ctx context.Context, id string, opts ...Option) error
	UpdateStatus(ctx context.Context, id string, status factcheck.StatusTopic, opts ...Option) (factcheck.Topic, error)
	UpdateDescription(ctx context.Context, id string, description string, opts ...Option) (factcheck.Topic, error)
	UpdateName(ctx context.Context, id string, name string, opts ...Option) (factcheck.Topic, error)
}

type topics struct {
	queries *postgres.Queries
}

func NewTopics(queries *postgres.Queries) Topics {
	return &topics{queries: queries}
}

func (t *topics) List(ctx context.Context, limit, offset int, opts ...Option) ([]factcheck.Topic, error) {
	queries := queries(t.queries, options(opts...))
	rows, err := queries.ListTopics(ctx, postgres.ListTopicsParams{
		Column1: limit,
		Column2: offset,
	})
	if err != nil {
		return nil, err
	}
	return utils.MapNoError(rows, postgres.ToTopicFromRow), nil
}

// ListAll retrieves all topics (backward compatibility)
func (t *topics) ListAll(ctx context.Context) ([]factcheck.Topic, error) {
	return t.List(ctx, 0, 0)
}

func (t *topics) Exists(ctx context.Context, id string, opts ...Option) (bool, error) {
	queries := queries(t.queries, options(opts...))
	uuid, err := postgres.UUID(id)
	if err != nil {
		return false, err
	}
	exists, err := queries.TopicExists(ctx, uuid)
	if err != nil {
		return false, handleNotFound(err, filter{"id": id})
	}
	return exists, nil
}

func (t *topics) GetStatus(ctx context.Context, id string, opts ...Option) (factcheck.StatusTopic, error) {
	queries := queries(t.queries, options(opts...))
	uuid, err := postgres.UUID(id)
	if err != nil {
		return "", err
	}
	row, err := queries.GetTopicStatus(ctx, uuid)
	if err != nil {
		return "", handleNotFound(err, filter{"id": id})
	}
	return factcheck.StatusTopic(row), nil
}

func (t *topics) Resolve(ctx context.Context, id string, answerText string, opts ...Option) (factcheck.Topic, error) {
	queries := queries(t.queries, options(opts...))
	uuid, err := postgres.UUID(id)
	if err != nil {
		return factcheck.Topic{}, err
	}
	result, err := postgres.Text(answerText)
	if err != nil {
		return factcheck.Topic{}, err
	}
	status, err := postgres.Text(factcheck.StatusTopicResolved)
	if err != nil {
		return factcheck.Topic{}, err
	}
	resolved, err := queries.ResolveTopic(ctx, postgres.ResolveTopicParams{
		ID:           uuid,
		Result:       result,
		ResultStatus: status,
	})
	if err != nil {
		return factcheck.Topic{}, handleNotFound(err, filter{"id": id})
	}
	return postgres.ToTopic(resolved), nil
}

func (t *topics) ListDynamicV2(ctx context.Context, limit, offset int, opts ...OptionTopic) ([]factcheck.Topic, error) {
	limit, offset = sanitize(limit, offset)
	options := options(opts...)
	queries := queries(t.queries, options.Options)
	rows, err := queries.ListTopicsDynamicV2(ctx, postgres.ListTopicsDynamicV2Params{
		Column1: options.LikeID,
		Column2: utils.MapNoError(options.Statuses, utils.String[factcheck.StatusTopic, string]),
		Column3: options.LikeMessageText,
		Column4: int32(limit),  //nolint:gosec
		Column5: int32(offset), //nolint:gosec
	})
	if err != nil {
		return nil, err
	}
	return postgres.ToTopics(rows), nil
}

func (t *topics) CountByStatusDynamicV2(ctx context.Context, opts ...OptionTopic) (map[factcheck.StatusTopic]int64, error) {
	options := options(opts...)
	queries := queries(t.queries, options.Options)
	if len(options.Statuses) != 0 {
		slog.WarnContext(ctx, "Statuses is not supported in CountByStatusDynamic", "statuses", options.Statuses)
	}
	rows, err := queries.CountTopicsGroupByStatusDynamicV2(ctx, postgres.CountTopicsGroupByStatusDynamicV2Params{
		Column1: options.LikeID,
		Column2: options.LikeMessageText,
	})
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
	return utils.MapNoError(rows, postgres.ToTopicFromStatusRow), nil
}

func (t *topics) ListLikeID(ctx context.Context, idPattern string, limit, offset int, opts ...Option) ([]factcheck.Topic, error) {
	limit, offset = sanitize(limit, offset)
	queries := queries(t.queries, options(opts...))
	rows, err := queries.ListTopicsLikeID(ctx, postgres.ListTopicsLikeIDParams{
		Column1: substringAuto(idPattern),
		Column2: limit,
		Column3: offset,
	})
	if err != nil {
		return nil, err
	}
	return utils.MapNoError(rows, postgres.ToTopicFromIDRow), nil
}

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
	return utils.MapNoError(rows, postgres.ToTopic), nil
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
