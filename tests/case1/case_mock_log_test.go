package case1

import (
	"context"
	"log/slog"
	"testing"

	"github.com/go-kod/kod"
	"github.com/stretchr/testify/assert"
)

func TestMockLog(t *testing.T) {
	t.Parallel()
	log, observer := kod.NewLogObserver()

	kod.RunTest(t, func(ctx context.Context, k Test1Component) {

		_, err := k.Foo(ctx, &FooReq{Id: 1})
		assert.Equal(t, "test1:B", err.Error())
		assert.Equal(t, 6, observer.Len())
		assert.Equal(t, 2, observer.Filter(func(r slog.Record) bool {
			return r.Level == slog.LevelError
		}).Len())
		assert.Equal(t, 0, observer.Clean().Len())
		slog.Info("test")
		assert.Equal(t, 1, observer.Len())
	}, kod.WithLogWrapper(log))
}
