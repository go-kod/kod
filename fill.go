package kod

import (
	"fmt"
	"log/slog"
	"reflect"
)

func fillLog(obj any, log *slog.Logger) error {
	x, ok := obj.(interface{ setLogger(*slog.Logger) })
	if !ok {
		return fmt.Errorf("fillLog: %T does not implement kod.Implements", obj)
	}

	x.setLogger(log)
	return nil
}

func fillRefs(impl any, get func(reflect.Type) (any, error)) error {
	p := reflect.ValueOf(impl)
	if p.Kind() != reflect.Pointer {
		return fmt.Errorf("fillRefs: %T not a pointer", impl)
	}

	s := p.Elem()
	if s.Kind() != reflect.Struct {
		return fmt.Errorf("fillRefs: %T not a struct pointer", impl)
	}

	for i, n := 0, s.NumField(); i < n; i++ {
		f := s.Field(i)
		if !f.CanAddr() {
			continue
		}
		p := reflect.NewAt(f.Type(), f.Addr().UnsafePointer()).Interface()
		x, ok := p.(interface{ setRef(any) })
		if !ok {
			continue
		}

		// Set the component.
		valueField := f.Field(0)
		component, err := get(valueField.Type())
		if err != nil {
			return fmt.Errorf("fillRefs: setting field %v.%s: %w", s.Type(), s.Type().Field(i).Name, err)
		}
		x.setRef(component)
	}
	return nil
}
