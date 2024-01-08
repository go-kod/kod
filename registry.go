package kod

import (
	"context"
	"fmt"
	"reflect"
)

type LocalStubFnInfo struct {
	Impl   any
	Name   string
	Caller string
}

type Registration struct {
	Name        string       // full package-prefixed component name
	Iface       reflect.Type // interface type for the component
	Impl        reflect.Type // implementation type (struct)
	Refs        string
	LocalStubFn func(context.Context, *LocalStubFnInfo) any
}

var (
	regs = make([]*Registration, 0)
)

func Register(reg Registration) {
	regs = append(regs, &reg)
}

func getRegs() []*Registration {
	return regs
}

// func (k *Kod) getRegistryByType(t reflect.Type) *Registration {
// 	k.mu.Lock()
// 	defer k.mu.Unlock()

// 	if reg, ok := k.registryByIface[t]; ok {
// 		return reg
// 	}

// 	if reg, ok := k.registryByImpl[t]; ok {
// 		return reg
// 	}

// 	return nil
// }

func (k *Kod) getImpl(ctx context.Context, t reflect.Type) (any, error) {
	k.mu.Lock()
	defer k.mu.Unlock()

	reg, ok := k.registryByImpl[t]
	if !ok {
		return nil, fmt.Errorf("kod: no component registered for implement %v", t)
	}

	return k.get(ctx, reg)
}

func (k *Kod) getIntf(ctx context.Context, t reflect.Type, caller string) (any, error) {
	reg, ok := k.registryByIface[t]
	if !ok {
		return nil, fmt.Errorf("kod: no component registered for interface %v", t)
	}

	comp, err := k.get(ctx, reg)
	if err != nil {
		return nil, err
	}

	return reg.LocalStubFn(ctx, &LocalStubFnInfo{
		Name:   reg.Name,
		Impl:   comp,
		Caller: caller,
	}), nil
}

func (k *Kod) get(ctx context.Context, reg *Registration) (any, error) {
	if c, ok := k.components[reg.Name]; ok {
		return c, nil
	}

	if fake, ok := k.opts.fakes[reg.Iface]; ok {
		// We have a fake registered for this component.
		return fake, nil
	}

	v := reflect.New(reg.Impl)
	obj := v.Interface()

	// Fill config.
	if cfg := config(v); cfg != nil {
		err := k.viper.UnmarshalKey(reg.Name, cfg)
		if err != nil {
			return nil, err
		}
	}

	if err := fillLog(obj, k.log.With("component", reg.Name)); err != nil {
		return nil, err
	}

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
		k.deferFn(deferFn{Name: reg.Name, Fn: i.Stop})
	}

	k.components[reg.Name] = obj

	return obj, nil
}
