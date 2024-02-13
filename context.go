package kod

import (
	"context"
)

var (
	kodContextKey = struct{}{}
)

// FromContext returns the Kod value stored in ctx, if any.
func FromContext(ctx context.Context) *Kod {
	if v, ok := ctx.Value(kodContextKey).(*Kod); ok {
		return v
	}
	return nil
}

// newContext returns a new Context that carries value kod.
func newContext(ctx context.Context, kod *Kod) context.Context {
	return context.WithValue(ctx, kodContextKey, kod)
}
