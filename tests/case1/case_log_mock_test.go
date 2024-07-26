package case1

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-kod/kod"
)

func TestMockLog(t *testing.T) {
	t.Parallel()
	log, observer := kod.NewLogObserver()
	os.Setenv("KOD_LOG_LEVEL", "error")

	kod.RunTest(t, func(ctx context.Context, k Test1Component) {
		_, err := k.Foo(ctx, &FooReq{Id: 1})
		require.Equal(t, "test1:B", err.Error())
		require.Equal(t, 2, observer.Len(), observer.All())
		require.Equal(t, 2, observer.Filter(func(r slog.Record) bool {
			return r.Level == slog.LevelError
		}).Len())
		require.Equal(t, 2, observer.ErrorCount())
		require.Equal(t, 0, observer.Clean().Len())
		slog.Info("test")
		require.Equal(t, 0, observer.Len())
		require.Equal(t, 0, observer.ErrorCount())
	}, kod.WithLogWrapper(log))
}
