// Code generated by "kod generate". DO NOT EDIT.
//go:build !ignoreKodGen

package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-kod/kod"
	"github.com/go-kod/kod/interceptor"
	"net/http"
	"reflect"
)

// Full method names for components.
const (
	// test1Controller_ComponentName is the full name of the component [test1Controller].
	test1Controller_ComponentName = "github.com/go-kod/kod/tests/graphcase/test1Controller"
	// HTTPController_ComponentName is the full name of the component [HTTPController].
	HTTPController_ComponentName = "github.com/go-kod/kod/tests/graphcase/HTTPController"
	// testService_ComponentName is the full name of the component [testService].
	testService_ComponentName = "github.com/go-kod/kod/tests/graphcase/testService"
	// testService_Foo_FullMethodName is the full name of the method [serviceImpl.Foo].
	testService_Foo_FullMethodName = "github.com/go-kod/kod/tests/graphcase/testService.Foo"
	// testModel_ComponentName is the full name of the component [testModel].
	testModel_ComponentName = "github.com/go-kod/kod/tests/graphcase/testModel"
	// testModel_Foo_FullMethodName is the full name of the method [modelImpl.Foo].
	testModel_Foo_FullMethodName = "github.com/go-kod/kod/tests/graphcase/testModel.Foo"
	// Test1Component_ComponentName is the full name of the component [Test1Component].
	Test1Component_ComponentName = "github.com/go-kod/kod/tests/graphcase/Test1Component"
	// Test1Component_Foo_FullMethodName is the full name of the method [test1Component.Foo].
	Test1Component_Foo_FullMethodName = "github.com/go-kod/kod/tests/graphcase/Test1Component.Foo"
)

func init() {
	kod.Register(&kod.Registration{
		Name:      "github.com/go-kod/kod/tests/graphcase/test1Controller",
		Interface: reflect.TypeOf((*test1Controller)(nil)).Elem(),
		Impl:      reflect.TypeOf(test1ControllerImpl{}),
		Refs: `⟦54c533d5:KoDeDgE:github.com/go-kod/kod/tests/graphcase/test1Controller→github.com/go-kod/kod/tests/graphcase/HTTPController⟧,
⟦f932c69a:KoDeDgE:github.com/go-kod/kod/tests/graphcase/test1Controller→github.com/go-kod/kod/tests/graphcase/Test1Component⟧`,
		LocalStubFn: func(ctx context.Context, info *kod.LocalStubFnInfo) any {
			return test1Controller_local_stub{
				impl:        info.Impl.(test1Controller),
				interceptor: info.Interceptor,
			}
		},
	})
	kod.Register(&kod.Registration{
		Name:      "github.com/go-kod/kod/tests/graphcase/HTTPController",
		Interface: reflect.TypeOf((*HTTPController)(nil)).Elem(),
		Impl:      reflect.TypeOf(httpControllerImpl{}),
		Refs:      `⟦38b48264:KoDeDgE:github.com/go-kod/kod/tests/graphcase/HTTPController→github.com/go-kod/kod/tests/graphcase/testService⟧`,
		LocalStubFn: func(ctx context.Context, info *kod.LocalStubFnInfo) any {
			return hTTPController_local_stub{
				impl:        info.Impl.(HTTPController),
				interceptor: info.Interceptor,
			}
		},
	})
	kod.Register(&kod.Registration{
		Name:      "github.com/go-kod/kod/tests/graphcase/testService",
		Interface: reflect.TypeOf((*testService)(nil)).Elem(),
		Impl:      reflect.TypeOf(serviceImpl{}),
		Refs:      `⟦e691e13e:KoDeDgE:github.com/go-kod/kod/tests/graphcase/testService→github.com/go-kod/kod/tests/graphcase/testModel⟧`,
		LocalStubFn: func(ctx context.Context, info *kod.LocalStubFnInfo) any {
			return testService_local_stub{
				impl:        info.Impl.(testService),
				interceptor: info.Interceptor,
			}
		},
	})
	kod.Register(&kod.Registration{
		Name:      "github.com/go-kod/kod/tests/graphcase/testModel",
		Interface: reflect.TypeOf((*testModel)(nil)).Elem(),
		Impl:      reflect.TypeOf(modelImpl{}),
		Refs:      ``,
		LocalStubFn: func(ctx context.Context, info *kod.LocalStubFnInfo) any {
			return testModel_local_stub{
				impl:        info.Impl.(testModel),
				interceptor: info.Interceptor,
			}
		},
	})
	kod.Register(&kod.Registration{
		Name:      "github.com/go-kod/kod/tests/graphcase/Test1Component",
		Interface: reflect.TypeOf((*Test1Component)(nil)).Elem(),
		Impl:      reflect.TypeOf(test1Component{}),
		Refs:      ``,
		LocalStubFn: func(ctx context.Context, info *kod.LocalStubFnInfo) any {
			return test1Component_local_stub{
				impl:        info.Impl.(Test1Component),
				interceptor: info.Interceptor,
			}
		},
	})
	kod.Register(&kod.Registration{
		Name:      "github.com/go-kod/kod/Main",
		Interface: reflect.TypeOf((*kod.Main)(nil)).Elem(),
		Impl:      reflect.TypeOf(App{}),
		Refs: `⟦b628eb85:KoDeDgE:github.com/go-kod/kod/Main→github.com/go-kod/kod/tests/graphcase/test1Controller⟧,
⟦75680c21:KoDeDgE:github.com/go-kod/kod/Main→github.com/go-kod/kod/tests/graphcase/Test1Component⟧`,
		LocalStubFn: nil,
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
var _ kod.InstanceOf[test1Controller] = (*test1ControllerImpl)(nil)
var _ kod.InstanceOf[HTTPController] = (*httpControllerImpl)(nil)
var _ kod.InstanceOf[testService] = (*serviceImpl)(nil)
var _ kod.InstanceOf[testModel] = (*modelImpl)(nil)
var _ kod.InstanceOf[Test1Component] = (*test1Component)(nil)
var _ kod.InstanceOf[kod.Main] = (*App)(nil)

// Local stub implementations.
// test1Controller_local_stub is a local stub implementation of [test1Controller].
type test1Controller_local_stub struct {
	impl        test1Controller
	interceptor interceptor.Interceptor
}

// Check that [test1Controller_local_stub] implements the [test1Controller] interface.
var _ test1Controller = (*test1Controller_local_stub)(nil)

// Foo wraps the method [test1ControllerImpl.Foo].
func (s test1Controller_local_stub) Foo(a0 *gin.Context) {
	// Because the first argument is not context.Context, so interceptors are not supported.
	s.impl.Foo(a0)
	return
}

// hTTPController_local_stub is a local stub implementation of [HTTPController].
type hTTPController_local_stub struct {
	impl        HTTPController
	interceptor interceptor.Interceptor
}

// Check that [hTTPController_local_stub] implements the [HTTPController] interface.
var _ HTTPController = (*hTTPController_local_stub)(nil)

// Foo wraps the method [httpControllerImpl.Foo].
func (s hTTPController_local_stub) Foo(a0 http.ResponseWriter, a1 http.Request) {
	// Because the first argument is not context.Context, so interceptors are not supported.
	s.impl.Foo(a0, a1)
	return
}

// testService_local_stub is a local stub implementation of [testService].
type testService_local_stub struct {
	impl        testService
	interceptor interceptor.Interceptor
}

// Check that [testService_local_stub] implements the [testService] interface.
var _ testService = (*testService_local_stub)(nil)

// Foo wraps the method [serviceImpl.Foo].
func (s testService_local_stub) Foo(ctx context.Context) (err error) {

	if s.interceptor == nil {
		err = s.impl.Foo(ctx)
		return
	}

	call := func(ctx context.Context, info interceptor.CallInfo, req, res []any) (err error) {
		err = s.impl.Foo(ctx)
		return
	}

	info := interceptor.CallInfo{
		Impl:       s.impl,
		FullMethod: testService_Foo_FullMethodName,
	}

	err = s.interceptor(ctx, info, []any{}, []any{}, call)
	return
}

// testModel_local_stub is a local stub implementation of [testModel].
type testModel_local_stub struct {
	impl        testModel
	interceptor interceptor.Interceptor
}

// Check that [testModel_local_stub] implements the [testModel] interface.
var _ testModel = (*testModel_local_stub)(nil)

// Foo wraps the method [modelImpl.Foo].
func (s testModel_local_stub) Foo(ctx context.Context) (err error) {

	if s.interceptor == nil {
		err = s.impl.Foo(ctx)
		return
	}

	call := func(ctx context.Context, info interceptor.CallInfo, req, res []any) (err error) {
		err = s.impl.Foo(ctx)
		return
	}

	info := interceptor.CallInfo{
		Impl:       s.impl,
		FullMethod: testModel_Foo_FullMethodName,
	}

	err = s.interceptor(ctx, info, []any{}, []any{}, call)
	return
}

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
