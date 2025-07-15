package repo

import (
	"context"

	"github.com/kaogeek/line-fact-check/factcheck/data/postgres"
)

type (
	Options struct{ tx Tx }
	Option  func(Options) Options

	Tx postgres.Tx
)

func (r *Repository) Begin(ctx context.Context) (Tx, error) {
	return r.TxnManager.Begin(ctx)
}

func WithTx(tx Tx) Option {
	return func(o Options) Options {
		o.tx = tx
		return o
	}
}
