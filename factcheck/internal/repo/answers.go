package repo

import (
	"context"
	"log/slog"

	"github.com/kaogeek/line-fact-check/factcheck"
	"github.com/kaogeek/line-fact-check/factcheck/data/postgres"
)

type Answers interface {
	Create(context.Context, factcheck.Answer, ...Option) (factcheck.Answer, error)
	GetByID(ctx context.Context, id string, opts ...Option) (factcheck.Answer, error)
	GetByTopicID(ctx context.Context, topicID string, opts ...Option) (factcheck.Answer, error)
	ListByTopicID(ctx context.Context, topicID string, opts ...Option) ([]factcheck.Answer, error)
	Delete(ctx context.Context, id string, opts ...Option) error
}

func NewAnswers(queries *postgres.Queries) Answers {
	return &answers{queries: queries}
}

type answers struct {
	queries *postgres.Queries
}

func (a *answers) Create(ctx context.Context, answer factcheck.Answer, opts ...Option) (factcheck.Answer, error) {
	queries := queries(a.queries, options(opts...))
	if answer.Text == "" {
		slog.WarnContext(ctx, "empty answer.text", "answer_id", answer.ID)
	}
	params, err := postgres.AnswerCreator(answer)
	if err != nil {
		return factcheck.Answer{}, err
	}
	created, err := queries.CreateAnswer(ctx, params)
	if err != nil {
		return factcheck.Answer{}, err
	}
	return postgres.ToAnswer(created)
}

func (a *answers) GetByID(ctx context.Context, id string, opts ...Option) (factcheck.Answer, error) {
	queries := queries(a.queries, options(opts...))
	uuid, err := postgres.UUID(id)
	if err != nil {
		return factcheck.Answer{}, err
	}
	result, err := queries.GetAnswerByID(ctx, uuid)
	if err != nil {
		return factcheck.Answer{}, handleNotFound(err, map[string]string{"id": id})
	}
	return postgres.ToAnswer(result)
}

func (a *answers) GetByTopicID(ctx context.Context, topicID string, opts ...Option) (factcheck.Answer, error) {
	queries := queries(a.queries, options(opts...))
	topicUUID, err := postgres.UUID(topicID)
	if err != nil {
		return factcheck.Answer{}, err
	}
	result, err := queries.GetAnswerByTopicID(ctx, topicUUID)
	if err != nil {
		return factcheck.Answer{}, handleNotFound(err, map[string]string{"topic_id": topicID})
	}
	return postgres.ToAnswer(result)
}

func (a *answers) ListByTopicID(ctx context.Context, topicID string, opts ...Option) ([]factcheck.Answer, error) {
	queries := queries(a.queries, options(opts...))
	topicUUID, err := postgres.UUID(topicID)
	if err != nil {
		return nil, err
	}
	result, err := queries.ListAnswersByTopicID(ctx, topicUUID)
	if err != nil {
		return nil, err
	}
	return postgres.ToAnswers(result)
}

func (a *answers) Delete(ctx context.Context, id string, opts ...Option) error {
	queries := queries(a.queries, options(opts...))
	uuid, err := postgres.UUID(id)
	if err != nil {
		return err
	}
	err = queries.DeleteAnswer(ctx, uuid)
	return handleNotFound(err, filter{"id": id})
}
