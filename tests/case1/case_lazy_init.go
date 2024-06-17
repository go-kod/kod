package case1

import (
	"context"

	"github.com/go-kod/kod"
)

type lazyInitImpl struct {
	kod.Implements[LazyInitImpl]
	test kod.Ref[LazyInitComponent]
}

func (t *lazyInitImpl) Init(ctx context.Context) error {
	t.L(ctx).Info("lazyInitImpl init...")

	return nil
}

func (t *lazyInitImpl) Try(ctx context.Context) {
	t.L(ctx).Info("Hello, World!")
}

func (t *lazyInitImpl) Shutdown(ctx context.Context) error {
	t.L(ctx).Info("lazyInitImpl shutdown...")

	return nil
}

type lazyInitComponent struct {
	kod.Implements[LazyInitComponent]
	kod.LazyInit
}

func (t *lazyInitComponent) Init(ctx context.Context) error {
	t.L(ctx).Info("lazyInitComponent init...")

	return nil
}

func (t *lazyInitComponent) Try(ctx context.Context) error {
	t.L(ctx).Info("Just do it!")

	return nil
}

func (t *lazyInitComponent) Shutdown(ctx context.Context) error {
	t.L(ctx).Info("lazyInitComponent shutdown...")

	return nil
}
