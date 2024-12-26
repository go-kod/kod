// Code generated by "kod generate". DO NOT EDIT.
//go:build !ignoreKodGen

package helloworld

import (
	"context"
	"github.com/go-kod/kod"
	"github.com/go-kod/kod/interceptor"
	"reflect"
)

// Full method names for components.
const (
	// HelloWorld_SayHello_FullMethodName is the full name of the method [helloWorld.SayHello].
	HelloWorld_SayHello_FullMethodName = "github.com/go-kod/kod/examples/helloworld/HelloWorld.SayHello"
	// HelloWorldLazy_SayHello_FullMethodName is the full name of the method [lazyHelloWorld.SayHello].
	HelloWorldLazy_SayHello_FullMethodName = "github.com/go-kod/kod/examples/helloworld/HelloWorldLazy.SayHello"
	// HelloWorldInterceptor_SayHello_FullMethodName is the full name of the method [helloWorldInterceptor.SayHello].
	HelloWorldInterceptor_SayHello_FullMethodName = "github.com/go-kod/kod/examples/helloworld/HelloWorldInterceptor.SayHello"
)

func init() {
	kod.Register(&kod.Registration{
		Name:      "github.com/go-kod/kod/Main",
		Interface: reflect.TypeOf((*kod.Main)(nil)).Elem(),
		Impl:      reflect.TypeOf(App{}),
		Refs: `⟦bda493e9:KoDeDgE:github.com/go-kod/kod/Main→github.com/go-kod/kod/examples/helloworld/HelloWorld⟧,
⟦b60b3708:KoDeDgE:github.com/go-kod/kod/Main→github.com/go-kod/kod/examples/helloworld/HelloWorldLazy⟧,
⟦c811f6f3:KoDeDgE:github.com/go-kod/kod/Main→github.com/go-kod/kod/examples/helloworld/HelloWorldInterceptor⟧`,
		LocalStubFn: nil,
	})
	kod.Register(&kod.Registration{
		Name:      "github.com/go-kod/kod/examples/helloworld/HelloWorld",
		Interface: reflect.TypeOf((*HelloWorld)(nil)).Elem(),
		Impl:      reflect.TypeOf(helloWorld{}),
		Refs:      ``,
		LocalStubFn: func(ctx context.Context, info *kod.LocalStubFnInfo) any {
			return helloWorld_local_stub{
				impl:        info.Impl.(HelloWorld),
				interceptor: info.Interceptor,
			}
		},
	})
	kod.Register(&kod.Registration{
		Name:      "github.com/go-kod/kod/examples/helloworld/HelloWorldLazy",
		Interface: reflect.TypeOf((*HelloWorldLazy)(nil)).Elem(),
		Impl:      reflect.TypeOf(lazyHelloWorld{}),
		Refs:      ``,
		LocalStubFn: func(ctx context.Context, info *kod.LocalStubFnInfo) any {
			return helloWorldLazy_local_stub{
				impl:        info.Impl.(HelloWorldLazy),
				interceptor: info.Interceptor,
			}
		},
	})
	kod.Register(&kod.Registration{
		Name:      "github.com/go-kod/kod/examples/helloworld/HelloWorldInterceptor",
		Interface: reflect.TypeOf((*HelloWorldInterceptor)(nil)).Elem(),
		Impl:      reflect.TypeOf(helloWorldInterceptor{}),
		Refs:      ``,
		LocalStubFn: func(ctx context.Context, info *kod.LocalStubFnInfo) any {
			return helloWorldInterceptor_local_stub{
				impl:        info.Impl.(HelloWorldInterceptor),
				interceptor: info.Interceptor,
			}
		},
	})
}

// CodeGen version check.
var _ kod.CodeGenLatestVersion = kod.CodeGenVersion[[0][1]struct{}](`
ERROR: You generated this file with 'kod generate' (codegen
version v0.1.0). The generated code is incompatible with the version of the
github.com/go-kod/kod module that you're using. The kod module
version can be found in your go.mod file or by running the following command.

    go list -m github.com/go-kod/kod

We recommend updating the kod module and the 'kod generate' command by
running the following.

    go get github.com/go-kod/kod@latest
    go install github.com/go-kod/kod/cmd/kod@latest

Then, re-run 'kod generate' and re-build your code. If the problem persists,
please file an issue at https://github.com/go-kod/kod/issues.
`)

// kod.InstanceOf checks.
var _ kod.InstanceOf[kod.Main] = (*App)(nil)
var _ kod.InstanceOf[HelloWorld] = (*helloWorld)(nil)
var _ kod.InstanceOf[HelloWorldLazy] = (*lazyHelloWorld)(nil)
var _ kod.InstanceOf[HelloWorldInterceptor] = (*helloWorldInterceptor)(nil)

// Local stub implementations.
// helloWorld_local_stub is a local stub implementation of [HelloWorld].
type helloWorld_local_stub struct {
	impl        HelloWorld
	interceptor interceptor.Interceptor
}

// Check that [helloWorld_local_stub] implements the [HelloWorld] interface.
var _ HelloWorld = (*helloWorld_local_stub)(nil)

// SayHello wraps the method [helloWorld.SayHello].
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
		FullMethod: HelloWorld_SayHello_FullMethodName,
	}

	_ = s.interceptor(ctx, info, []any{}, []any{}, call)
}

// helloWorldLazy_local_stub is a local stub implementation of [HelloWorldLazy].
type helloWorldLazy_local_stub struct {
	impl        HelloWorldLazy
	interceptor interceptor.Interceptor
}

// Check that [helloWorldLazy_local_stub] implements the [HelloWorldLazy] interface.
var _ HelloWorldLazy = (*helloWorldLazy_local_stub)(nil)

// SayHello wraps the method [lazyHelloWorld.SayHello].
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
		FullMethod: HelloWorldLazy_SayHello_FullMethodName,
	}

	_ = s.interceptor(ctx, info, []any{}, []any{}, call)
}

// helloWorldInterceptor_local_stub is a local stub implementation of [HelloWorldInterceptor].
type helloWorldInterceptor_local_stub struct {
	impl        HelloWorldInterceptor
	interceptor interceptor.Interceptor
}

// Check that [helloWorldInterceptor_local_stub] implements the [HelloWorldInterceptor] interface.
var _ HelloWorldInterceptor = (*helloWorldInterceptor_local_stub)(nil)

// SayHello wraps the method [helloWorldInterceptor.SayHello].
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
		FullMethod: HelloWorldInterceptor_SayHello_FullMethodName,
	}

	_ = s.interceptor(ctx, info, []any{}, []any{}, call)
}
