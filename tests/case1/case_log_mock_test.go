package case1

import (
	"context"
	"log/slog"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-kod/kod"
)

func TestMockLog(t *testing.T) {
	log, observer := kod.NewTestLogger()
	slog.SetDefault(log)
	t.Setenv("KOD_LOG_LEVEL", "error")

	kod.RunTest(t, func(ctx context.Context, k Test1Component) {
		_, err := k.Foo(ctx, &FooReq{Id: 1})
		require.Equal(t, "test1:B", err.Error())
		require.Equal(t, 5, observer.Len(), observer.String())
		require.Equal(t, 2, observer.Filter(func(r map[string]any) bool {
			return r["level"] == slog.LevelError.String()
		}).Len())
		require.Equal(t, 2, observer.ErrorCount())
		require.Equal(t, 0, observer.Clean().Len())
		slog.Info("test")
		require.Equal(t, 1, observer.Len())
		require.Equal(t, 0, observer.ErrorCount())
	})
}
