package kod

import (
	"context"
	"os"
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
	k, err := newKod(context.Background(), options{})
	assert.Nil(t, err)

	assert.EqualError(t, k.parseConfig("nosuffix"), "read config file: Unsupported Config Type \"\"")
}

func TestConfigEnv(t *testing.T) {
	k, err := newKod(context.Background(), options{})
	assert.Nil(t, err)

	assert.Equal(t, k.config.Name, "kod.test")
	assert.Equal(t, k.config.Version, "")
	assert.Equal(t, k.config.Env, "local")

	os.Setenv("KOD_NAME", "test")
	os.Setenv("KOD_VERSION", "1.0.0")
	os.Setenv("KOD_ENV", "dev")

	k, err = newKod(context.Background(), options{})
	assert.Nil(t, err)

	assert.Equal(t, k.config.Name, "test")
	assert.Equal(t, k.config.Version, "1.0.0")
	assert.Equal(t, k.config.Env, "dev")
}
