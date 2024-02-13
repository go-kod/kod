package kod

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFill(t *testing.T) {
	t.Run("case 1", func(t *testing.T) {
		assert.NotNil(t, fillLog(nil, nil))
	})

	t.Run("case 2", func(t *testing.T) {
		assert.NotNil(t, fillRefs(nil, nil))
	})

	t.Run("case 3", func(t *testing.T) {
		i := 0
		assert.NotNil(t, fillRefs(&i, nil))
	})

}
