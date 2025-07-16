package repo

import (
	"context"

	"github.com/kaogeek/line-fact-check/factcheck/data/postgres"
)

type (
	Options struct {
		tx    Tx
		level IsoLevel
	}
	Option   func(Options) Options
	IsoLevel string
	Tx       postgres.Tx
)

const (
	ReadCommitted  IsoLevel = IsoLevel(postgres.IsoLevelReadCommitted)
	RepeatableRead IsoLevel = IsoLevel(postgres.IsoLevelRepeatableRead)
	Serializable   IsoLevel = IsoLevel(postgres.IsoLevelSerializable)
)

func WithTx(tx Tx) Option {
	return func(o Options) Options {
		return o.WithTx(tx)
	}
}

func (o Options) WithTx(tx Tx) Options {
	if o.tx != nil {
		panic("tx already set")
	}
	o.tx = tx
	return o
}

func WithIsolationLevel(level IsoLevel) Option {
	return func(o Options) Options {
		if o.level != "" {
			panic("isolation level already set")
		}
		o.level = level
		return o
	}
}

func (r *Repository) Begin(ctx context.Context) (Tx, error) {
	return r.TxnManager.Begin(ctx)
}

func (r *Repository) BeginTx(ctx context.Context, level IsoLevel) (Tx, error) {
	return r.TxnManager.BeginTx(ctx, postgres.IsoLevel(level)) // TODO:
}
