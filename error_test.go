package kod

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPanicError(t *testing.T) {
	t.Run("panic", func(t *testing.T) {
		defer func() {
			r := recover()
			assert.Contains(t, RecoverFrom(r).Error(), "panic caught: testpanic")
		}()

		panic("testpanic")
	})

}
