// Code generated by "kod generate". DO NOT EDIT.
//go:build !ignoreKodGen

package main

import (
	"context"
	"github.com/go-kod/kod"
	"github.com/go-kod/kod/interceptor"
	"reflect"
)

func init() {
	kod.Register(&kod.Registration{
		Name:  "github.com/go-kod/kod/examples/helloworld/HelloWorld",
		Iface: reflect.TypeOf((*HelloWorld)(nil)).Elem(),
		Impl:  reflect.TypeOf(helloWorld{}),
		Refs:  ``,
		LocalStubFn: func(ctx context.Context, info *kod.LocalStubFnInfo) any {
			var interceptors []kod.Interceptor
			if h, ok := info.Impl.(interface{ Interceptors() []kod.Interceptor }); ok {
				interceptors = h.Interceptors()
			}

			return helloWorld_local_stub{
				impl:        info.Impl.(HelloWorld),
				interceptor: interceptor.Chain(interceptors),
				name:        info.Name,
			}
		},
	})
	kod.Register(&kod.Registration{
		Name:  "github.com/go-kod/kod/Main",
		Iface: reflect.TypeOf((*kod.Main)(nil)).Elem(),
		Impl:  reflect.TypeOf(app{}),
		Refs:  `⟦bda493e9:KoDeDgE:github.com/go-kod/kod/Main→github.com/go-kod/kod/examples/helloworld/HelloWorld⟧`,
		LocalStubFn: func(ctx context.Context, info *kod.LocalStubFnInfo) any {
			var interceptors []kod.Interceptor
			if h, ok := info.Impl.(interface{ Interceptors() []kod.Interceptor }); ok {
				interceptors = h.Interceptors()
			}

			return main_local_stub{
				impl:        info.Impl.(kod.Main),
				interceptor: interceptor.Chain(interceptors),
				name:        info.Name,
			}
		},
	})
}

// kod.InstanceOf checks.
var _ kod.InstanceOf[HelloWorld] = (*helloWorld)(nil)
var _ kod.InstanceOf[kod.Main] = (*app)(nil)

// Local stub implementations.

type helloWorld_local_stub struct {
	impl        HelloWorld
	name        string
	interceptor kod.Interceptor
}

// Check that helloWorld_local_stub implements the HelloWorld interface.
var _ HelloWorld = (*helloWorld_local_stub)(nil)

func (s helloWorld_local_stub) SayHello() (r0 string) {

	if s.interceptor == nil {
		r0 = s.impl.SayHello()
		return
	}

	call := func(ctx context.Context, info kod.CallInfo, req, res []any) (err error) {
		r0 = s.impl.SayHello()
		res[0] = r0
		return
	}

	info := kod.CallInfo{
		Impl:       s.impl,
		Component:  s.name,
		FullMethod: "github.com/go-kod/kod/examples/helloworld/HelloWorld.SayHello",
		Method:     "SayHello",
	}

	ctx := context.Background()
	_ = s.interceptor(ctx, info, []any{}, []any{r0}, call)
	return
}

type main_local_stub struct {
	impl        kod.Main
	name        string
	interceptor kod.Interceptor
}

// Check that main_local_stub implements the kod.Main interface.
var _ kod.Main = (*main_local_stub)(nil)

