package kod

import (
	"context"
	"fmt"
	"reflect"
)

// Get returns a component from the Kod value stored in ctx.
func Get[T any](ctx context.Context) (T, error) {
	var zero T
	k := FromContext(ctx)
	if k == nil {
		return zero, fmt.Errorf("kod: no Kod in context")
	}
	return k.Get[T](ctx)
}

// Get returns a component by interface type or implementation pointer type.
func (k *Kod) Get[T any](ctx context.Context) (T, error) {
	var zero T
	if k == nil {
		return zero, fmt.Errorf("kod: nil Kod")
	}

	t := reflect.TypeFor[T]()

	var v any
	var err error
	switch t.Kind() {
	case reflect.Interface:
		v, err = k.getIntf(ctx, t)
	case reflect.Pointer:
		if t.Elem().Kind() != reflect.Struct {
			return zero, fmt.Errorf("kod: component type %v must be an interface or pointer to component implementation", t)
		}
		v, err = k.getImpl(ctx, t.Elem())
	default:
		return zero, fmt.Errorf("kod: component type %v must be an interface or pointer to component implementation", t)
	}
	if err != nil {
		return zero, err
	}

	got, ok := v.(T)
	if !ok {
		return zero, fmt.Errorf("kod: component %v has type %T", t, v)
	}
	return got, nil
}
