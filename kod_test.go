package kod

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
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

func TestConfigNoSuffix(t *testing.T) {
	k, err := newKod(context.Background())
	assert.Nil(t, err)

	assert.EqualError(t, k.parseConfig("nosuffix"), "read config file: Unsupported Config Type \"\"")
}

func TestConfigEnv(t *testing.T) {
	k, err := newKod(context.Background())
	assert.Nil(t, err)

	assert.Equal(t, k.config.Name, "kod.test")
	assert.Equal(t, k.config.Version, "")
	assert.Equal(t, k.config.Env, "local")

	t.Setenv("KOD_NAME", "test")
	t.Setenv("KOD_VERSION", "1.0.0")
	t.Setenv("KOD_ENV", "dev")

	k, err = newKod(context.Background())
	assert.Nil(t, err)

	assert.Equal(t, k.config.Name, "test")
	assert.Equal(t, k.config.Version, "1.0.0")
	assert.Equal(t, k.config.Env, "dev")
}
