package version

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestVersion(t *testing.T) {
	t.Run("version", func(t *testing.T) {
		require.Equal(t, "v1.0.0", SemVer{Major: 1}.String())
	})

	t.Run("self version", func(t *testing.T) {
		require.Equal(t, "(devel)", SelfVersion())
	})
}
