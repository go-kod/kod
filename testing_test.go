package kod

import (
	"context"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/go-kod/kod/internal/mock"
)

func Test_testRunner_sub(t *testing.T) {
	t.Run("failure", func(t *testing.T) {
		mock.ExpectFailure(t, func(tb testing.TB) {
			tb.Helper()

			r := &runner{}
			err := r.sub(tb, nil)
			if err != nil {
				tb.FailNow()
			}
		})
	})
}

func Test_checkRunFunc(t *testing.T) {
	t.Run("not a func", func(t *testing.T) {
		_, _, err := checkRunFunc(context.Background(), 0)
		assert.EqualError(t, err, "not a func")
	})

	t.Run("must not be variadic", func(t *testing.T) {
		_, _, err := checkRunFunc(context.Background(), func(t *testing.T, a ...int) {
		})
		assert.EqualError(t, err, "must not be variadic")
	})

	// must have no return outputs
	t.Run("must have no return outputs", func(t *testing.T) {
		_, _, err := checkRunFunc(context.Background(), func(t *testing.T, a int) int {
			return 0
		})
		assert.EqualError(t, err, "must have no return outputs")
	})

	t.Run("must have at least two args", func(t *testing.T) {
		_, _, err := checkRunFunc(context.Background(), func() {
		})
		assert.EqualError(t, err, "must have at least two args")
	})

	// function first argument type *testing.T does not match first kod.Run argument context.Context
	t.Run("function first argument type *testing.T does not match first kod.Run argument context.Context", func(t *testing.T) {
		_, _, err := checkRunFunc(context.Background(), func(t *testing.T, a *testing.T, b *testing.T) {
		})
		assert.EqualError(t, err, "function first argument type *testing.T does not match first kod.Run argument context.Context")
	})

	t.Run("function argument %d type %v must be a component interface or pointer to component implementation", func(t *testing.T) {
		_, _, err := checkRunFunc(context.Background(), func(ctx context.Context, t int) {
		})
		assert.EqualError(t, err, "function argument 1 type int must be a component interface or pointer to component implementation")
	})

	t.Run("ok", func(t *testing.T) {
		_, _, err := checkRunFunc(context.Background(), func(ctx context.Context, a *testComponent) {
		})
		assert.Nil(t, err)
	})
}

// extractComponentInterfaceType
func Test_extractComponentInterfaceType(t *testing.T) {
	t.Run("not a pointer", func(t *testing.T) {
		_, err := extractComponentInterfaceType(reflect.TypeOf(0))
		assert.EqualError(t, err, "type int is not a struct")
	})

	t.Run("not a struct pointer", func(t *testing.T) {
		_, err := extractComponentInterfaceType(reflect.TypeOf(&testing.T{}))
		assert.EqualError(t, err, "type *testing.T is not a struct")
	})

	t.Run("not a component interface", func(t *testing.T) {
		_, err := extractComponentInterfaceType(reflect.TypeOf(&struct{}{}))
		assert.EqualError(t, err, "type *struct {} is not a struct")
	})

	t.Run("not a struct", func(t *testing.T) {
		_, err := extractComponentInterfaceType(reflect.TypeOf(&struct {
			Implements[testing.T]
		}{}))
		assert.EqualError(t, err, "type *struct { kod.Implements[testing.T] } is not a struct")
	})

	t.Run("type struct {} does not embed kod.Implements", func(t *testing.T) {
		_, err := extractComponentInterfaceType(reflect.TypeOf(struct{}{}))
		assert.EqualError(t, err, "type struct {} does not embed kod.Implements")
	})

	t.Run("ok", func(t *testing.T) {
		_, err := extractComponentInterfaceType(reflect.TypeOf(struct {
			Implements[testing.T]
		}{}))
		assert.Nil(t, err)
	})
}
