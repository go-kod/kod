package kod

import (
	"context"
	"errors"
	"fmt"
	"reflect"

	"github.com/dominikbraun/graph"

	"github.com/go-kod/kod/interceptor"
	"github.com/go-kod/kod/internal/callgraph"
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
func (k *Kod) getIntf(ctx context.Context, t reflect.Type) (any, error) {
	reg, ok := k.registryByInterface[t]
	if !ok {
		return nil, fmt.Errorf("kod: no component registered for interface %v", t)
	}

	intf, ok := k.components[reg.Name]
	if ok {
		return intf, nil
	}

	impl, err := k.get(ctx, reg)
	if err != nil {
		return nil, err
	}

	interceptors := k.opts.interceptors
	if h, ok := impl.(interface {
		Interceptors() []interceptor.Interceptor
	}); ok {
		interceptors = append(interceptors, h.Interceptors()...)
	}

	info := &LocalStubFnInfo{
		Impl:        impl,
		Interceptor: interceptor.Chain(interceptors),
	}

	intf = reg.LocalStubFn(ctx, info)

	k.components[reg.Name] = intf

	return intf, nil
}

// get returns the component for the given registration.
func (k *Kod) get(ctx context.Context, reg *Registration) (any, error) {
	// Check if we already have the component.
	if c, ok := k.components[reg.Name]; ok {
		return c, nil
	}

	if fake, ok := k.opts.fakes[reg.Interface]; ok {
		// We have a fake registered for this component.
		return fake, nil
	}

	// Create a new instance of the component.
	v := reflect.New(reg.Impl)
	obj := v.Interface()

	// Fill global config.
	if c, ok := obj.(interface{ getGlobalConfig() any }); ok {
		if cfg := c.getGlobalConfig(); cfg != nil {
			err := k.Unmarshal("", cfg)
			if err != nil {
				return nil, err
			}
		}
	}

	// Fill config.
	if c, ok := obj.(interface{ getConfig() any }); ok {
		if cfg := c.getConfig(); cfg != nil {
			err := k.Unmarshal(reg.Name, cfg)
			if err != nil {
				return nil, err
			}
		}
	}

	// // Fill logger.
	// if err := fillLog(reg.Name, obj, slog.Default()); err != nil {
	// 	return nil, err
	// }

	// Fill refs.
	if err := fillRefs(obj, k.lazyInitComponents,
		func(t reflect.Type) componentGetter {
			return func() (any, error) {
				return k.getIntf(ctx, t)
			}
		}); err != nil {
		return nil, err
	}

	// Call Init if available.
	if i, ok := obj.(interface{ Init(context.Context) error }); ok {
		if err := i.Init(ctx); err != nil {
			return nil, fmt.Errorf("component %q initialization failed: %w", reg.Name, err)
		}
	}

	// Call Shutdown if available.
	if i, ok := obj.(interface{ Shutdown(context.Context) error }); ok {
		k.hooker.Add(hooks.HookFunc{Name: reg.Name, Fn: i.Shutdown})
	}

	// Cache the component.
	k.components[reg.Name] = obj

	return obj, nil
}

func fillRefs(impl any, lazyInit map[reflect.Type]bool, get func(reflect.Type) componentGetter) error {
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
		x, ok := p.(interface {
			setRef(bool, componentGetter)
		})
		if !ok {
			continue
		}

		// Set the component.
		valueField := f.Field(0)
		componentGetter := get(valueField.Type())
		isLazyInit := lazyInit[valueField.Type()]

		x.setRef(isLazyInit, componentGetter)
	}
	return nil
}

// checkCircularDependency checks that there are no circular dependencies
// between registered components.
func checkCircularDependency(reg []*Registration) error {
	g := graph.New(graph.StringHash, graph.Directed(), graph.PreventCycles())

	for _, reg := range reg {
		if err := g.AddVertex(reg.Name); err != nil {
			return fmt.Errorf("components [%s], error %s", reg.Name, err)
		}
	}

	var errs []error
	for _, reg := range reg {
		edges := callgraph.ParseEdges([]byte(reg.Refs))
		for _, edge := range edges {
			err := g.AddEdge(edge[0], edge[1])
			if err != nil {
				switch err {
				case graph.ErrEdgeAlreadyExists, graph.ErrEdgeCreatesCycle:
					err = fmt.Errorf("components [%s] and [%s] have cycle Ref", edge[0], edge[1])
				}
				errs = append(errs, err)
			}
		}
	}

	return errors.Join(errs...)
}

// processRegistrations checks that all registered component interfaces are
// implemented by a registered component implementation struct.
func processRegistrations(regs []*Registration) (map[reflect.Type]bool, error) {
	// Gather the set of registered interfaces.
	intfs := map[reflect.Type]struct{}{}
	for _, reg := range regs {
		intfs[reg.Interface] = struct{}{}
	}

	lazyInitComponents := make(map[reflect.Type]bool)

	// Check that for every kod.Ref[T] field in a component implementation
	// struct, T is a registered interface.
	var errs []error
	for _, reg := range regs {
		for i := 0; i < reg.Impl.NumField(); i++ {
			f := reg.Impl.Field(i)
			switch {
			case f.Type.Implements(reflect.TypeFor[interface{ isRef() }]()):
				// f is a kod.Ref[T].
				v := f.Type.Field(0) // a Ref[T]'s value field
				// v是func 类型，取它的第一个返回值的类型

				// check if v is interface or not
				if v.Type.Kind() != reflect.Interface {
					err := fmt.Errorf(
						"component implementation struct %v has field %v, but field type %v is not an interface",
						reg.Impl, f.Type, v.Type,
					)
					errs = append(errs, err)
					continue
				}

				if _, ok := intfs[v.Type]; !ok {
					// T is not a registered component interface.
					err := fmt.Errorf(
						"component implementation struct %v has field %v, but component %v was not registered; maybe you forgot to run 'kod generate'",
						reg.Impl, f.Type, v.Type,
					)
					errs = append(errs, err)
				}
			case f.Type.Implements(reflect.TypeFor[interface{ isLazyInit() }]()):
				// f is a kod.LazyInit.
				lazyInitComponents[reg.Interface] = true
			}
		}
	}
	return lazyInitComponents, errors.Join(errs...)
}
