// Code generated by "kod generate". DO NOT EDIT.
//go:build !ignoreKodGen

package case3

import (
	"context"
	"github.com/go-kod/kod"
	"github.com/go-kod/kod/interceptor"
	"reflect"
)

// Full method names for components.
const (
	// Test1Component_Foo_FullMethodName is the full name of the method [test1Component.Foo].
	Test1Component_Foo_FullMethodName = "github.com/go-kod/kod/tests/case3/Test1Component.Foo"
	// Test2Component_Foo_FullMethodName is the full name of the method [test2Component.Foo].
	Test2Component_Foo_FullMethodName = "github.com/go-kod/kod/tests/case3/Test2Component.Foo"
	// Test3Component_Foo_FullMethodName is the full name of the method [test3Component.Foo].
	Test3Component_Foo_FullMethodName = "github.com/go-kod/kod/tests/case3/Test3Component.Foo"
)

func init() {
	kod.Register(&kod.Registration{
		Name:      "github.com/go-kod/kod/tests/case3/Test1Component",
		Interface: reflect.TypeOf((*Test1Component)(nil)).Elem(),
		Impl:      reflect.TypeOf(test1Component{}),
		Refs:      `⟦27eee423:KoDeDgE:github.com/go-kod/kod/tests/case3/Test1Component→github.com/go-kod/kod/tests/case3/Test2Component⟧`,
		LocalStubFn: func(ctx context.Context, info *kod.LocalStubFnInfo) any {
			return test1Component_local_stub{
				impl:        info.Impl.(Test1Component),
				interceptor: info.Interceptor,
				name:        info.Name,
			}
		},
	})
	kod.Register(&kod.Registration{
		Name:      "github.com/go-kod/kod/tests/case3/Test2Component",
		Interface: reflect.TypeOf((*Test2Component)(nil)).Elem(),
		Impl:      reflect.TypeOf(test2Component{}),
		Refs:      `⟦9e26e2b7:KoDeDgE:github.com/go-kod/kod/tests/case3/Test2Component→github.com/go-kod/kod/tests/case3/Test3Component⟧`,
		LocalStubFn: func(ctx context.Context, info *kod.LocalStubFnInfo) any {
			return test2Component_local_stub{
				impl:        info.Impl.(Test2Component),
				interceptor: info.Interceptor,
				name:        info.Name,
			}
		},
	})
	kod.Register(&kod.Registration{
		Name:      "github.com/go-kod/kod/tests/case3/Test3Component",
		Interface: reflect.TypeOf((*Test3Component)(nil)).Elem(),
		Impl:      reflect.TypeOf(test3Component{}),
		Refs:      `⟦ab2f1bfd:KoDeDgE:github.com/go-kod/kod/tests/case3/Test3Component→github.com/go-kod/kod/tests/case3/Test1Component⟧`,
		LocalStubFn: func(ctx context.Context, info *kod.LocalStubFnInfo) any {
			return test3Component_local_stub{
				impl:        info.Impl.(Test3Component),
				interceptor: info.Interceptor,
				name:        info.Name,
			}
		},
	})
	kod.Register(&kod.Registration{
		Name:      "github.com/go-kod/kod/Main",
		Interface: reflect.TypeOf((*kod.Main)(nil)).Elem(),
		Impl:      reflect.TypeOf(App{}),
		Refs: `⟦964a80ec:KoDeDgE:github.com/go-kod/kod/Main→github.com/go-kod/kod/tests/case3/Test1Component⟧,
⟦e679c332:KoDeDgE:github.com/go-kod/kod/Main→github.com/go-kod/kod/tests/case3/Test2Component⟧`,
		LocalStubFn: nil,
	})
}

// kod.InstanceOf checks.
var _ kod.InstanceOf[Test1Component] = (*test1Component)(nil)
var _ kod.InstanceOf[Test2Component] = (*test2Component)(nil)
var _ kod.InstanceOf[Test3Component] = (*test3Component)(nil)
var _ kod.InstanceOf[kod.Main] = (*App)(nil)

// Local stub implementations.

// test1Component_local_stub is a local stub implementation of [Test1Component].
type test1Component_local_stub struct {
	impl        Test1Component
	name        string
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
		Component:  s.name,
		FullMethod: Test1Component_Foo_FullMethodName,
	}

	err = s.interceptor(ctx, info, []any{a1}, []any{}, call)
	return
}

// test2Component_local_stub is a local stub implementation of [Test2Component].
type test2Component_local_stub struct {
	impl        Test2Component
	name        string
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
		Component:  s.name,
		FullMethod: Test2Component_Foo_FullMethodName,
	}

	err = s.interceptor(ctx, info, []any{a1}, []any{}, call)
	return
}

// test3Component_local_stub is a local stub implementation of [Test3Component].
type test3Component_local_stub struct {
	impl        Test3Component
	name        string
	interceptor interceptor.Interceptor
}

// Check that [test3Component_local_stub] implements the [Test3Component] interface.
var _ Test3Component = (*test3Component_local_stub)(nil)

// Foo wraps the method [test3Component.Foo].
func (s test3Component_local_stub) Foo(ctx context.Context, a1 *FooReq) (err error) {

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
		Component:  s.name,
		FullMethod: Test3Component_Foo_FullMethodName,
	}

	err = s.interceptor(ctx, info, []any{a1}, []any{}, call)
	return
}

// main_local_stub is a local stub implementation of [kod.Main].
type main_local_stub struct {
	impl        kod.Main
	name        string
	interceptor interceptor.Interceptor
}

// Check that [main_local_stub] implements the [kod.Main] interface.
var _ kod.Main = (*main_local_stub)(nil)

