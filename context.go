package kod

import (
	"context"
)

var (
	kodContextKey = struct{}{}
)

func FromContext(ctx context.Context) *Kod {
	if v, ok := ctx.Value(kodContextKey).(*Kod); ok {
		return v
	}
	return nil
}

func newContext(ctx context.Context, kod *Kod) context.Context {
	return context.WithValue(ctx, kodContextKey, kod)
}
