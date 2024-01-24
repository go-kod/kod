package kod

import (
	"context"
	"fmt"
	"log/slog"
	"reflect"

	"github.com/go-kod/kod/internal/hooks"
	"github.com/go-kod/kod/internal/registry"
)

// LocalStubFnInfo is the information passed to LocalStubFn.
type LocalStubFnInfo = registry.LocalStubFnInfo

// Registration is the registration information for a component.
type Registration = registry.Registration

// Register registers the given component implementations.
var Register = registry.Register

// getImpl returns the component for the given implementation type.
func (k *Kod) getImpl(ctx context.Context, t reflect.Type) (any, error) {
	k.mu.Lock()
	defer k.mu.Unlock()

	reg, ok := k.registryByImpl[t]
	if !ok {
		return nil, fmt.Errorf("kod: no component registered for implement %v", t)
	}

	return k.get(ctx, reg)
}

// getIntf returns the component for the given interface type.
func (k *Kod) getIntf(ctx context.Context, t reflect.Type, caller string) (any, error) {

	reg, ok := k.registryByIface[t]
	if !ok {
		return nil, fmt.Errorf("kod: no component registered for interface %v", t)
	}

	intf, ok := k.components[reg.Name]
	if ok {
		return intf, nil
	}

	comp, err := k.get(ctx, reg)
	if err != nil {
		return nil, err
	}

	intf = reg.LocalStubFn(ctx, &LocalStubFnInfo{
		Name: reg.Name,
		Impl: comp,
	})

	k.components[reg.Name] = intf

	return intf, nil
}

// get returns the component for the given registration.
func (k *Kod) get(ctx context.Context, reg *Registration) (any, error) {
	// Check if we already have the component.
	if c, ok := k.components[reg.Name]; ok {
		return c, nil
	}

	if fake, ok := k.opts.fakes[reg.Iface]; ok {
		// We have a fake registered for this component.
		return fake, nil
	}

	// Create a new instance of the component.
	v := reflect.New(reg.Impl)
	obj := v.Interface()

	// Fill config.
	if c, ok := obj.(interface{ getConfig() any }); ok {
		if cfg := c.getConfig(); cfg != nil {
			err := k.viper.UnmarshalKey(reg.Name, cfg)
			if err != nil {
				return nil, err
			}
		}
	}

	// Fill logger.
	if err := fillLog(obj, k.log.With("component", reg.Name)); err != nil {
		return nil, err
	}

	// Fill refs.
	if err := fillRefs(obj, func(t reflect.Type) (any, error) {
		return k.getIntf(ctx, t, reg.Name)
	}); err != nil {
		return nil, err
	}

	// Call Init if available.
	if i, ok := obj.(interface{ Init(context.Context) error }); ok {
		if err := i.Init(ctx); err != nil {
			return nil, fmt.Errorf("component %q initialization failed: %w", reg.Name, err)
		}
	}

	// Call Stop if available.
	if i, ok := obj.(interface{ Stop(context.Context) error }); ok {
		k.hooker.Add(hooks.HookFunc{Name: reg.Name, Fn: i.Stop})
	}

	// Cache the component.
	k.components[reg.Name] = obj

	return obj, nil
}

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
