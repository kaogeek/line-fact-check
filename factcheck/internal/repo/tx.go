package repo

import (
	"context"

	"github.com/kaogeek/line-fact-check/factcheck/data/postgres"
)

type (
	// OptionTx represents transaction-related options
	OptionTx func(*TxOptions)

	// TxOptions contains transaction configuration
	TxOptions struct {
		Tx    Tx
		level IsoLevel
	}

	IsoLevel string
	Tx       postgres.Tx
)

const (
	ReadCommitted  IsoLevel = IsoLevel(postgres.IsoLevelReadCommitted)
	RepeatableRead IsoLevel = IsoLevel(postgres.IsoLevelRepeatableRead)
	Serializable   IsoLevel = IsoLevel(postgres.IsoLevelSerializable)
)

// WithTx sets the transaction for the operation
func WithTx(tx Tx) OptionTx {
	return func(o *TxOptions) {
		if o.Tx != nil {
			panic("tx already set")
		}
		o.Tx = tx
	}
}

// WithIsolationLevel sets the isolation level for the operation
func WithIsolationLevel(level IsoLevel) OptionTx {
	return func(o *TxOptions) {
		if o.level != "" {
			panic("isolation level already set")
		}
		o.level = level
	}
}

func (r *Repository) Begin(ctx context.Context) (Tx, error) {
	return r.TxnManager.Begin(ctx)
}

func (r *Repository) BeginTx(ctx context.Context, level IsoLevel) (Tx, error) {
	return r.TxnManager.BeginTx(ctx, postgres.IsoLevel(level))
}
