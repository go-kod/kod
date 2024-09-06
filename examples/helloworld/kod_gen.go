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
		Name:      "github.com/go-kod/kod/examples/helloworld/HelloWorld",
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
	kod.Register(&kod.Registration{
		Name:      "github.com/go-kod/kod/examples/helloworld/HelloWorldInterceptor",
		Interface: reflect.TypeOf((*HelloWorldInterceptor)(nil)).Elem(),
		Impl:      reflect.TypeOf(helloWorldInterceptor{}),
		Refs:      ``,
		LocalStubFn: func(ctx context.Context, info *kod.LocalStubFnInfo) any {
			interceptors := info.Interceptors
			if h, ok := info.Impl.(interface {
				Interceptors() []interceptor.Interceptor
			}); ok {
				interceptors = append(interceptors, h.Interceptors()...)
			}

			return helloWorldInterceptor_local_stub{
				impl:        info.Impl.(HelloWorldInterceptor),
				interceptor: interceptor.Chain(interceptors),
				name:        info.Name,
			}
		},
	})
	kod.Register(&kod.Registration{
		Name:      "github.com/go-kod/kod/examples/helloworld/HelloWorldLazy",
		Interface: reflect.TypeOf((*HelloWorldLazy)(nil)).Elem(),
		Impl:      reflect.TypeOf(lazyHelloWorld{}),
		Refs:      ``,
		LocalStubFn: func(ctx context.Context, info *kod.LocalStubFnInfo) any {
			interceptors := info.Interceptors
			if h, ok := info.Impl.(interface {
				Interceptors() []interceptor.Interceptor
			}); ok {
				interceptors = append(interceptors, h.Interceptors()...)
			}

			return helloWorldLazy_local_stub{
				impl:        info.Impl.(HelloWorldLazy),
				interceptor: interceptor.Chain(interceptors),
				name:        info.Name,
			}
		},
	})
	kod.Register(&kod.Registration{
		Name:      "github.com/go-kod/kod/Main",
		Interface: reflect.TypeOf((*kod.Main)(nil)).Elem(),
		Impl:      reflect.TypeOf(App{}),
		Refs: `⟦bda493e9:KoDeDgE:github.com/go-kod/kod/Main→github.com/go-kod/kod/examples/helloworld/HelloWorld⟧,
⟦b60b3708:KoDeDgE:github.com/go-kod/kod/Main→github.com/go-kod/kod/examples/helloworld/HelloWorldLazy⟧,
⟦c811f6f3:KoDeDgE:github.com/go-kod/kod/Main→github.com/go-kod/kod/examples/helloworld/HelloWorldInterceptor⟧`,
		LocalStubFn: func(ctx context.Context, info *kod.LocalStubFnInfo) any {
			interceptors := info.Interceptors
			if h, ok := info.Impl.(interface {
				Interceptors() []interceptor.Interceptor
			}); ok {
				interceptors = append(interceptors, h.Interceptors()...)
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
var _ kod.InstanceOf[HelloWorldInterceptor] = (*helloWorldInterceptor)(nil)
var _ kod.InstanceOf[HelloWorldLazy] = (*lazyHelloWorld)(nil)
var _ kod.InstanceOf[kod.Main] = (*App)(nil)

// Local stub implementations.

type helloWorld_local_stub struct {
	impl        HelloWorld
	name        string
	interceptor interceptor.Interceptor
}

// Check that helloWorld_local_stub implements the HelloWorld interface.
var _ HelloWorld = (*helloWorld_local_stub)(nil)

func (s helloWorld_local_stub) SayHello(ctx context.Context) {

	if s.interceptor == nil {
		s.impl.SayHello(ctx)
		return
	}

	call := func(ctx context.Context, info interceptor.CallInfo, req, res []any) (err error) {
		s.impl.SayHello(ctx)
		return
	}

	info := interceptor.CallInfo{
		Impl:       s.impl,
		Component:  s.name,
		FullMethod: "github.com/go-kod/kod/examples/helloworld/HelloWorld.SayHello",
		Method:     "SayHello",
	}

	_ = s.interceptor(ctx, info, []any{}, []any{}, call)
}

type helloWorldInterceptor_local_stub struct {
	impl        HelloWorldInterceptor
	name        string
	interceptor interceptor.Interceptor
}

// Check that helloWorldInterceptor_local_stub implements the HelloWorldInterceptor interface.
var _ HelloWorldInterceptor = (*helloWorldInterceptor_local_stub)(nil)

func (s helloWorldInterceptor_local_stub) SayHello(ctx context.Context) {

	if s.interceptor == nil {
		s.impl.SayHello(ctx)
		return
	}

	call := func(ctx context.Context, info interceptor.CallInfo, req, res []any) (err error) {
		s.impl.SayHello(ctx)
		return
	}

	info := interceptor.CallInfo{
		Impl:       s.impl,
		Component:  s.name,
		FullMethod: "github.com/go-kod/kod/examples/helloworld/HelloWorldInterceptor.SayHello",
		Method:     "SayHello",
	}

	_ = s.interceptor(ctx, info, []any{}, []any{}, call)
}

type helloWorldLazy_local_stub struct {
	impl        HelloWorldLazy
	name        string
	interceptor interceptor.Interceptor
}

// Check that helloWorldLazy_local_stub implements the HelloWorldLazy interface.
var _ HelloWorldLazy = (*helloWorldLazy_local_stub)(nil)

func (s helloWorldLazy_local_stub) SayHello(ctx context.Context) {

	if s.interceptor == nil {
		s.impl.SayHello(ctx)
		return
	}

	call := func(ctx context.Context, info interceptor.CallInfo, req, res []any) (err error) {
		s.impl.SayHello(ctx)
		return
	}

	info := interceptor.CallInfo{
		Impl:       s.impl,
		Component:  s.name,
		FullMethod: "github.com/go-kod/kod/examples/helloworld/HelloWorldLazy.SayHello",
		Method:     "SayHello",
	}

	_ = s.interceptor(ctx, info, []any{}, []any{}, call)
}

type main_local_stub struct {
	impl        kod.Main
	name        string
	interceptor interceptor.Interceptor
}

// Check that main_local_stub implements the kod.Main interface.
var _ kod.Main = (*main_local_stub)(nil)

