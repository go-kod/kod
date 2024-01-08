package kod

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/goleak"
)

func TestMain(m *testing.M) {
	goleak.VerifyTestMain(m)
}

func TestConfigNoSuffix(t *testing.T) {
	k, err := newKod(nil, options{})
	assert.Nil(t, err)

	assert.EqualError(t, k.parseConfig("nosuffix"), "read config file: Unsupported Config Type \"\"")

	assert.Equal(t, "info", k.Config().Log.Level)
}
