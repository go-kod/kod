package kod

import (
	"context"
	"testing"
	"time"

	"github.com/knadh/koanf/v2"
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

func TestConfigNoSuffix(t *testing.T) {
	k, err := newKod(context.Background())
	assert.Nil(t, err)

	assert.EqualError(t, k.parseConfig("nosuffix"), "read config file: Unsupported Config Type \"\"")
}

func TestConfigNoFile(t *testing.T) {
	k, err := newKod(context.Background())
	assert.Nil(t, err)

	assert.EqualError(t, k.parseConfig("notfound.yaml"), "read config file: open notfound.yaml: no such file or directory")
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

type testComponent struct {
	Implements[testInterface]
	WithConfig[testConfig]
	initialized bool
	initErr     error
	shutdown    bool
	shutdownErr error
}

type testConfig struct {
	Value string `default:"default"`
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

func TestConfigurationLoading(t *testing.T) {
	tests := []struct {
		name     string
		koanf    *koanf.Koanf
		filename string
		wantErr  bool
	}{
		{
			name:  "custom koanf",
			koanf: koanf.New("."), // 使用 koanf.New() 替代空实例
		},
		{
			name:     "invalid file extension",
			filename: "config.invalid",
			wantErr:  true,
		},
		{
			name:     "missing file",
			filename: "notexist.yaml",
			wantErr:  true, // Should use defaults
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			opts := []func(*options){
				WithConfigFile(tt.filename),
			}
			if tt.koanf != nil {
				opts = append(opts, WithKoanf(tt.koanf))
			}

			k, err := newKod(context.Background(), opts...)
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			require.NoError(t, err)

			cfg := k.Config()
			assert.NotEmpty(t, cfg.Name)
			assert.NotEmpty(t, cfg.Env)
			assert.Equal(t, 5*time.Second, cfg.ShutdownTimeout)
		})
	}
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
