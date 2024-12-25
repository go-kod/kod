package case1

import (
	"context"
	"log/slog"
	"testing"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/internal/mock"
	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	log, observer := kod.NewTestLogger()
	slog.SetDefault(log)

	kod.RunTest(t, func(ctx context.Context, k *test1Component) {
		observer.Clean()

		k.L(ctx).Info("hello", "config", k.Config())

		assert.Equal(t, 1, observer.Len())
		assert.Equal(t, "{\"component\":\"github.com/go-kod/kod/tests/case1/Test1Component\",\"config\":{\"A\":\"B\",\"Redis\":{\"Addr\":\"localhost:6379\",\"Timeout\":2000000000}},\"level\":\"INFO\",\"msg\":\"hello\"}\n",
			observer.RemoveKeys("time").String())
	})
}

func TestDefaultConfig2(t *testing.T) {
	log, observer := kod.NewTestLogger()
	slog.SetDefault(log)

	kod.RunTest(t, func(ctx context.Context, k *test1Component) {
		observer.Clean()

		k.L(ctx).Info("hello", "config", k.Config())

		assert.Equal(t, 1, observer.Len())
		assert.Equal(t, "{\"component\":\"github.com/go-kod/kod/tests/case1/Test1Component\",\"config\":{\"A\":\"B2\",\"Redis\":{\"Addr\":\"localhost:6379\",\"Timeout\":1000000000}},\"level\":\"INFO\",\"msg\":\"hello\"}\n",
			observer.RemoveKeys("time").String())
	}, kod.WithConfigFile("./kod2.toml"))
}

func TestDefaultConfigError(t *testing.T) {
	mock.ExpectFailure(t, func(tb testing.TB) {
		kod.RunTest(tb, func(ctx context.Context, k *test1ComponentDefaultErrorImpl) {
			k.L(ctx).Info("hello", "config", k.Config())
		}, kod.WithConfigFile("./kod2.toml"))
	})
}

func TestDefaultGlobalConfigError(t *testing.T) {
	mock.ExpectFailure(t, func(tb testing.TB) {
		kod.RunTest(tb, func(ctx context.Context, k *test1ComponentGlobalDefaultErrorImpl) {
			k.L(ctx).Info("hello", "config", k.Config())
		}, kod.WithConfigFile("./kod2.toml"))
	})
}
