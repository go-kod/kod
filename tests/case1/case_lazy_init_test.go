package case1

import (
	"context"
	"testing"

	"github.com/go-kod/kod"
	"github.com/stretchr/testify/require"
)

func TestLazyInit(t *testing.T) {
	t.Parallel()
	log, observer := kod.NewLogObserver()

	kod.RunTest(t, func(ctx context.Context, k *lazyInitImpl) {
		require.Equal(t, 1, observer.Len(), observer.All())

		k.Try(ctx)

		require.Equal(t, 2, observer.Len(), observer.All())

		k.test.Get()
		require.Equal(t, 3, observer.Len(), observer.All())

		require.Nil(t, k.test.Get().Try(ctx))
		require.Equal(t, 4, observer.Len(), observer.All())
	}, kod.WithLogWrapper(log))
}

func TestLazyInitTest(t *testing.T) {
	t.Parallel()
	log, observer := kod.NewLogObserver()

	kod.RunTest2(t, func(ctx context.Context, k *lazyInitImpl, comp *lazyInitComponent) {
		k.Try(ctx)

		require.Equal(t, 3, observer.Len(), observer.All())

		require.Equal(t, k.test.Get(), k.test.Get())

		require.Equal(t, 3, observer.Len(), observer.All())
	}, kod.WithLogWrapper(log))
}

func TestLazyInitTest2(t *testing.T) {
	t.Parallel()
	log, observer := kod.NewLogObserver()

	kod.RunTest2(t, func(ctx context.Context, k LazyInitImpl, comp LazyInitComponent) {
		require.Equal(t, 2, observer.Len(), observer.All())

		k.Try(ctx)

		require.Equal(t, 3, observer.Len(), observer.All())
	}, kod.WithLogWrapper(log))
}

func TestLazyInitTest3(t *testing.T) {
	t.Parallel()
	log, observer := kod.NewLogObserver()

	kod.RunTest2(t, func(ctx context.Context, k *lazyInitImpl, comp LazyInitComponent) {
		k.Try(ctx)

		require.Equal(t, 3, observer.Len(), observer.All())

		require.Equal(t, comp, k.test.Get())
		require.Equal(t, 3, observer.Len(), observer.All())

		require.Nil(t, k.test.Get().Try(ctx))
		require.Equal(t, 4, observer.Len(), observer.All())
	}, kod.WithLogWrapper(log))
}
