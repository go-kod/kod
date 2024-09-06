package case1

import (
	"context"
	"log/slog"
	"testing"

	"github.com/go-kod/kod"
	"github.com/stretchr/testify/assert"
)

func TestLogLevel(t *testing.T) {
	log, observer := kod.NewTestLogger()

	kod.RunTest(t, func(ctx context.Context, k *test1Component) {
		observer.Clean()

		k.L(ctx).Debug("debug")
		k.L(ctx).WithGroup("group").Info("info")

		assert.Equal(t, 1, observer.Len())
		assert.Equal(t, 0, observer.Filter(func(r map[string]any) bool {
			return r["level"] == slog.LevelDebug.String()
		}).Len())
	}, kod.WithLogger(log))

	t.Setenv("KOD_LOG_LEVEL", "debug")

	kod.RunTest(t, func(ctx context.Context, k *test1Component) {
		observer.Clean()

		k.L(ctx).Debug("debug")
		k.L(ctx).WithGroup("group").Info("info")

		assert.Equal(t, 1, observer.Len())
		assert.Equal(t, 0, observer.Filter(func(r map[string]any) bool {
			return r["level"] == slog.LevelDebug.String()
		}).Len())
	}, kod.WithLogger(log))
}
