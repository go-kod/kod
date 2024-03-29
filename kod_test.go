package kod

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m,
		goleak.IgnoreAnyFunction("github.com/go-kod/kod/internal/ratelimit.cpuproc"),
		goleak.IgnoreAnyFunction("gopkg.in/natefinch/lumberjack%2ev2.(*Logger).millRun"),
	)
}

func TestConfigNoSuffix(t *testing.T) {
	k, err := newKod(options{})
	assert.Nil(t, err)

	assert.EqualError(t, k.parseConfig("nosuffix"), "read config file: Unsupported Config Type \"\"")

	assert.Equal(t, "info", k.Config().Log.Level)
}
