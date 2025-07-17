package repo

import (
	"context"

	"github.com/kaogeek/line-fact-check/factcheck/data/postgres"
)

type (
	IsoLevel string
	Tx       postgres.Tx
)

const (
	ReadCommitted  IsoLevel = IsoLevel(postgres.IsoLevelReadCommitted)
	RepeatableRead IsoLevel = IsoLevel(postgres.IsoLevelRepeatableRead)
	Serializable   IsoLevel = IsoLevel(postgres.IsoLevelSerializable)
)

func (r *Repository) Begin(ctx context.Context) (Tx, error) {
	return r.TxnManager.Begin(ctx)
}

func (r *Repository) BeginTx(ctx context.Context, level IsoLevel) (Tx, error) {
	return r.TxnManager.BeginTx(ctx, postgres.IsoLevel(level))
}

// WithTx sets the transaction for the operation
func WithTx(tx Tx) Option {
	return func(o *Options) {
		if o.tx != nil {
			panic("tx already set")
		}
		o.tx = tx
	}
}

func queries(q *postgres.Queries, o Options) *postgres.Queries {
	if o.tx == nil {
		return q
	}
	return q.WithTx(o.tx)
}
