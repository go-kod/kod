// Code generated by "kod generate". DO NOT EDIT.
//go:build !ignoreKodGen

package helloworld

import (
	"context"
	"github.com/go-kod/kod"
	"github.com/go-kod/kod/interceptor"
	"reflect"
)

func init() {
	kod.Register(&kod.Registration{
		Name:      "github.com/go-kod/kod/examples/app/helloworld/HelloWorld",
		Interface: reflect.TypeOf((*HelloWorld)(nil)).Elem(),
		Impl:      reflect.TypeOf(helloWorld{}),
		Refs:      ``,
		LocalStubFn: func(ctx context.Context, info *kod.LocalStubFnInfo) any {
			interceptors := info.Interceptors
			if h, ok := info.Impl.(interface {
				Interceptors() []interceptor.Interceptor
			}); ok {
				interceptors = append(interceptors, h.Interceptors()...)
			}

			return helloWorld_local_stub{
				impl:        info.Impl.(HelloWorld),
				interceptor: interceptor.Chain(interceptors),
				name:        info.Name,
			}
		},
	})
}

// kod.InstanceOf checks.
var _ kod.InstanceOf[HelloWorld] = (*helloWorld)(nil)

// Local stub implementations.

type helloWorld_local_stub struct {
	impl        HelloWorld
	name        string
	interceptor interceptor.Interceptor
}

// Check that helloWorld_local_stub implements the HelloWorld interface.
var _ HelloWorld = (*helloWorld_local_stub)(nil)

func (s helloWorld_local_stub) SayHello(ctx context.Context) (r0 string) {

	if s.interceptor == nil {
		r0 = s.impl.SayHello(ctx)
		return
	}

	call := func(ctx context.Context, info interceptor.CallInfo, req, res []any) (err error) {
		r0 = s.impl.SayHello(ctx)
		res[0] = r0
		return
	}

	info := interceptor.CallInfo{
		Impl:       s.impl,
		Component:  s.name,
		FullMethod: "github.com/go-kod/kod/examples/app/helloworld/HelloWorld.SayHello",
		Method:     "SayHello",
	}

	_ = s.interceptor(ctx, info, []any{}, []any{r0}, call)
	return
}
