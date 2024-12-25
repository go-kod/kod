package case1

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-kod/kod"
)

func TestLazyInit(t *testing.T) {
	log, observer := kod.NewTestLogger()
	slog.SetDefault(log)

	kod.RunTest(t, func(ctx context.Context, k *lazyInitImpl) {
		require.Equal(t, 1, observer.Len(), observer.String())

		k.Try(ctx)

		require.Equal(t, 2, observer.Len(), observer.String())

		k.test.Get()
		require.Equal(t, 3, observer.Len(), observer.String())

		require.Nil(t, k.test.Get().Try(ctx))
		require.Equal(t, 4, observer.Len(), observer.String())
	})
}

func TestLazyInitTest(t *testing.T) {
	log, observer := kod.NewTestLogger()
	slog.SetDefault(log)

	kod.RunTest2(t, func(ctx context.Context, k *lazyInitImpl, comp *lazyInitComponent) {
		k.Try(ctx)

		require.Equal(t, 3, observer.Len(), observer.String())

		require.Equal(t, k.test.Get(), k.test.Get())

		require.Equal(t, 3, observer.Len(), observer.String())
	})
}

func TestLazyInitTest2(t *testing.T) {
	log, observer := kod.NewTestLogger()
	slog.SetDefault(log)

	kod.RunTest2(t, func(ctx context.Context, k LazyInitImpl, comp LazyInitComponent) {
		require.Equal(t, 2, observer.Len(), observer.String())

		k.Try(ctx)

		require.Equal(t, 3, observer.Len(), observer.String())
	})
}

func TestLazyInitTest3(t *testing.T) {
	log, observer := kod.NewTestLogger()
	slog.SetDefault(log)

	kod.RunTest2(t, func(ctx context.Context, k *lazyInitImpl, comp LazyInitComponent) {
		k.Try(ctx)

		require.Equal(t, 3, observer.Len(), observer.String())

		require.Equal(t, comp, k.test.Get())
		require.Equal(t, 3, observer.Len(), observer.String())

		require.Nil(t, k.test.Get().Try(ctx))
		require.Equal(t, 4, observer.Len(), observer.String())
	})
}
