package case1

import (
	"context"
	"log/slog"
	"testing"

	"github.com/go-kod/kod"
	"github.com/stretchr/testify/assert"
)

func TestLogLevel(t *testing.T) {
	obs, log := kod.NewLogObserver()

	kod.RunTest(t, func(ctx context.Context, k *test1Component) {
		log.Clean()

		k.L(ctx).Debug("debug")
		k.L(ctx).WithGroup("group").Info("info")

		assert.Equal(t, 1, log.Len())
		assert.Equal(t, 0, log.Filter(func(r slog.Record) bool {
			return r.Level == slog.LevelDebug
		}).Len())
	}, kod.WithLogWrapper(obs))

	t.Setenv("KOD_LOG_LEVEL", "debug")

	kod.RunTest(t, func(ctx context.Context, k *test1Component) {
		log.Clean()

		k.L(ctx).Debug("debug")
		k.L(ctx).WithGroup("group").Info("info")

		assert.Equal(t, 2, log.Len())
		assert.Equal(t, 1, log.Filter(func(r slog.Record) bool {
			return r.Level == slog.LevelDebug
		}).Len())
	}, kod.WithLogWrapper(obs))
}
