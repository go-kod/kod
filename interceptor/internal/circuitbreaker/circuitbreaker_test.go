package circuitbreaker

import (
	"context"
	"testing"
	"time"

	"github.com/sony/gobreaker"
	"github.com/stretchr/testify/assert"
)

func TestCircuitBreaker(t *testing.T) {
	t.Run("case 1", func(t *testing.T) {
		cb := NewCircuitBreaker(context.Background(), "case1")
		done, err := cb.Allow()
		assert.Nil(t, err)
		done(context.Canceled)
		done, err = cb.Allow()
		assert.Nil(t, err)
		done(context.Canceled)
		done, err = cb.Allow()
		assert.Nil(t, err)
		done(context.Canceled)
		_, err = cb.Allow()
		assert.Equal(t, gobreaker.ErrOpenState, err)
		_, err = cb.Allow()
		assert.Equal(t, gobreaker.ErrOpenState, err)
		_, err = cb.Allow()
		assert.Equal(t, gobreaker.ErrOpenState, err)
		_, err = cb.Allow()
		assert.Equal(t, gobreaker.ErrOpenState, err)
		_, err = cb.Allow()
		assert.Equal(t, gobreaker.ErrOpenState, err)
		_, err = cb.Allow()
		assert.Equal(t, gobreaker.ErrOpenState, err)

		time.Sleep(time.Second)
		_, err = cb.Allow()
		assert.Equal(t, gobreaker.ErrOpenState, err)
	})
}
