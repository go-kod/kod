package registry

import (
	"context"
	"reflect"

	"github.com/go-kod/kod/interceptor"
)

// LocalStubFnInfo is the information passed to LocalStubFn.
type LocalStubFnInfo struct {
	Impl        any
	Interceptor interceptor.Interceptor
}

// Registration is the registration information for a component.
type Registration struct {
	Name        string       // full package-prefixed component name
	Interface   reflect.Type // interface type for the component
	Impl        reflect.Type // implementation type (struct)
	Refs        string
	LocalStubFn func(context.Context, *LocalStubFnInfo) any
}

// regs is the list of registered components.
var regs = make([]*Registration, 0)

// Register registers the given component implementations.
func Register(reg *Registration) {
	regs = append(regs, reg)
}

func All() []*Registration {
	return regs
}
