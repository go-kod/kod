// Code generated by "kod generate". DO NOT EDIT.
//go:build !ignoreKodGen

package case2

import (
	"context"
	"github.com/go-kod/kod"
	"github.com/go-kod/kod/interceptor"
	"reflect"
)

// Full method names for components.
const (
	// Test1Component_Foo_FullMethodName is the full name of the method [test1Component.Foo].
	Test1Component_Foo_FullMethodName = "github.com/go-kod/kod/tests/case2/Test1Component.Foo"
	// Test2Component_Foo_FullMethodName is the full name of the method [test2Component.Foo].
	Test2Component_Foo_FullMethodName = "github.com/go-kod/kod/tests/case2/Test2Component.Foo"
)

func init() {
	kod.Register(&kod.Registration{
		Name:      "github.com/go-kod/kod/tests/case2/Test1Component",
		Interface: reflect.TypeOf((*Test1Component)(nil)).Elem(),
		Impl:      reflect.TypeOf(test1Component{}),
		Refs:      `⟦3dc9f060:KoDeDgE:github.com/go-kod/kod/tests/case2/Test1Component→github.com/go-kod/kod/tests/case2/Test2Component⟧`,
		LocalStubFn: func(ctx context.Context, info *kod.LocalStubFnInfo) any {
			return test1Component_local_stub{
				impl:        info.Impl.(Test1Component),
				interceptor: info.Interceptor,
			}
		},
	})
	kod.Register(&kod.Registration{
		Name:      "github.com/go-kod/kod/tests/case2/Test2Component",
		Interface: reflect.TypeOf((*Test2Component)(nil)).Elem(),
		Impl:      reflect.TypeOf(test2Component{}),
		Refs:      `⟦1767cee9:KoDeDgE:github.com/go-kod/kod/tests/case2/Test2Component→github.com/go-kod/kod/tests/case2/Test1Component⟧`,
		LocalStubFn: func(ctx context.Context, info *kod.LocalStubFnInfo) any {
			return test2Component_local_stub{
				impl:        info.Impl.(Test2Component),
				interceptor: info.Interceptor,
			}
		},
	})
	kod.Register(&kod.Registration{
		Name:        "github.com/go-kod/kod/Main",
		Interface:   reflect.TypeOf((*kod.Main)(nil)).Elem(),
		Impl:        reflect.TypeOf(App{}),
		Refs:        `⟦73dc6a0b:KoDeDgE:github.com/go-kod/kod/Main→github.com/go-kod/kod/tests/case2/Test1Component⟧`,
		LocalStubFn: nil,
	})
}

// CodeGen version check.
var _ kod.CodeGenLatestVersion = kod.CodeGenVersion[[0][1]struct{}](`
ERROR: You generated this file with 'kod generate' (devel) (codegen
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
var _ kod.InstanceOf[Test1Component] = (*test1Component)(nil)
var _ kod.InstanceOf[Test2Component] = (*test2Component)(nil)
var _ kod.InstanceOf[kod.Main] = (*App)(nil)

// Local stub implementations.
// test1Component_local_stub is a local stub implementation of [Test1Component].
type test1Component_local_stub struct {
	impl        Test1Component
	interceptor interceptor.Interceptor
}

// Check that [test1Component_local_stub] implements the [Test1Component] interface.
var _ Test1Component = (*test1Component_local_stub)(nil)

// Foo wraps the method [test1Component.Foo].
func (s test1Component_local_stub) Foo(ctx context.Context, a1 *FooReq) (err error) {

	if s.interceptor == nil {
		err = s.impl.Foo(ctx, a1)
		return
	}

	call := func(ctx context.Context, info interceptor.CallInfo, req, res []any) (err error) {
		err = s.impl.Foo(ctx, a1)
		return
	}

	info := interceptor.CallInfo{
		Impl:       s.impl,
		FullMethod: Test1Component_Foo_FullMethodName,
	}

	err = s.interceptor(ctx, info, []any{a1}, []any{}, call)
	return
}

// test2Component_local_stub is a local stub implementation of [Test2Component].
type test2Component_local_stub struct {
	impl        Test2Component
	interceptor interceptor.Interceptor
}

// Check that [test2Component_local_stub] implements the [Test2Component] interface.
var _ Test2Component = (*test2Component_local_stub)(nil)

// Foo wraps the method [test2Component.Foo].
func (s test2Component_local_stub) Foo(ctx context.Context, a1 *FooReq) (err error) {

	if s.interceptor == nil {
		err = s.impl.Foo(ctx, a1)
		return
	}

	call := func(ctx context.Context, info interceptor.CallInfo, req, res []any) (err error) {
		err = s.impl.Foo(ctx, a1)
		return
	}

	info := interceptor.CallInfo{
		Impl:       s.impl,
		FullMethod: Test2Component_Foo_FullMethodName,
	}

	err = s.interceptor(ctx, info, []any{a1}, []any{}, call)
	return
}
