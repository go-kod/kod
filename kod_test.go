package kod

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m,
		goleak.IgnoreAnyFunction("github.com/go-kod/kod/interceptor/internal/ratelimit.cpuproc"),
		goleak.IgnoreAnyFunction("go.opentelemetry.io/otel/sdk/metric.(*PeriodicReader).run"),
		goleak.IgnoreAnyFunction("go.opentelemetry.io/otel/sdk/trace.(*batchSpanProcessor).processQueue"),
		goleak.IgnoreAnyFunction("go.opentelemetry.io/otel/sdk/log.exportSync.func1"),
		goleak.IgnoreAnyFunction("go.opentelemetry.io/otel/sdk/log.(*BatchProcessor).poll.func1"),
	)
}

type testComponent struct {
	Implements[testInterface]
	initialized bool
	initErr     error
	shutdown    bool
	shutdownErr error
}

type testInterface interface {
	IsInitialized() bool
}

func (c *testComponent) Init(context.Context) error {
	c.initialized = true
	return c.initErr
}

func (c *testComponent) Shutdown(context.Context) error {
	c.shutdown = true
	return c.shutdownErr
}

func (c *testComponent) IsInitialized() bool {
	return c.initialized
}

func (c *testComponent) implements(testInterface) {}

func TestWithShutdownTimeout(t *testing.T) {
	k, err := newKod(context.Background(), WithShutdownTimeout(time.Second))
	require.NoError(t, err)
	assert.Equal(t, time.Second, k.shutdownTimeout)
}

func TestDeferHooks(t *testing.T) {
	k, err := newKod(context.Background())
	require.NoError(t, err)

	executed := false
	k.Defer("test", func(context.Context) error {
		executed = true
		return nil
	})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	k.hooker.Do(ctx)
	assert.True(t, executed)
}
