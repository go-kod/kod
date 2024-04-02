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
	kod.Register(&kod.Registration{
		Name:      "github.com/go-kod/kod/Main",
		Interface: reflect.TypeOf((*kod.Main)(nil)).Elem(),
		Impl:      reflect.TypeOf(app{}),
		Refs:      `⟦562af859:KoDeDgE:github.com/go-kod/kod/Main→github.com/go-kod/kod/examples/app/helloworld/HelloWorld⟧`,
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
var _ kod.InstanceOf[kod.Main] = (*app)(nil)

// Local stub implementations.

type helloWorld_local_stub struct {
	impl        HelloWorld
	name        string
	interceptor interceptor.Interceptor
}

// Check that helloWorld_local_stub implements the HelloWorld interface.
var _ HelloWorld = (*helloWorld_local_stub)(nil)

func (s helloWorld_local_stub) SayHello() (r0 string) {
	// Because the first argument is not context.Context, so interceptors are not supported.
	r0 = s.impl.SayHello()
	return
}

type main_local_stub struct {
	impl        kod.Main
	name        string
	interceptor interceptor.Interceptor
}

// Check that main_local_stub implements the kod.Main interface.
var _ kod.Main = (*main_local_stub)(nil)

