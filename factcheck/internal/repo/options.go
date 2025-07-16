package repo

type (
	// Option is a function that modifies the Options struct
	Option func(*Options)

	// Options is a base option for repository operations/methods
	// It can be embedded inside of other options
	Options struct {
		tx Tx
	}
)

func (o Options) Clone() []Option {
	if o.tx == nil {
		return nil
	}
	return []Option{
		WithTx(o.tx),
	}
}

func (o Options) Tx() Tx {
	return o.tx
}

func options[O any, F ~func(*O)](opts ...F) O {
	var option O
	return apply(option, opts...)
}

func apply[O any, F ~func(*O)](o O, opts ...F) O {
	for _, opt := range opts {
		opt(&o)
	}
	return o
}
