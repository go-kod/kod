package kod

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
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

func TestRunRejectsMainValue(t *testing.T) {
	dir := t.TempDir()
	wd, err := os.Getwd()
	require.NoError(t, err)

	require.NoError(t, os.WriteFile(filepath.Join(dir, "go.mod"), []byte(`module runvaluetest

go 1.27

toolchain go1.27rc1

replace github.com/go-kod/kod => `+wd+`

require github.com/go-kod/kod v0.0.0
`), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "run_value_test.go"), []byte(`package runvaluetest

import (
	"context"
	"testing"

	"github.com/go-kod/kod"
)

type app struct {
	kod.Implements[kod.Main]
}

func TestRunValue(t *testing.T) {
	_ = kod.Run(context.Background(), func(context.Context, app) error {
		return nil
	})
}
`), 0o644))

	cmd := exec.Command("go", "test", "-mod=mod", ".")
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	require.Error(t, err, string(out))
	require.Contains(t, string(out), "does not satisfy kod.PointerToMain")
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
