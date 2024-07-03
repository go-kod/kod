// Code generated by "kod generate". DO NOT EDIT.
//go:build !ignoreKodGen

package case1

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-kod/kod"
	"github.com/go-kod/kod/interceptor"
	"github.com/labstack/echo/v4"
	"net/http"
	"reflect"
)

func init() {
	kod.Register(&kod.Registration{
		Name:      "github.com/go-kod/kod/tests/case1/HTTPController",
		Interface: reflect.TypeOf((*HTTPController)(nil)).Elem(),
		Impl:      reflect.TypeOf(httpControllerImpl{}),
		Refs:      ``,
		LocalStubFn: func(ctx context.Context, info *kod.LocalStubFnInfo) any {
			interceptors := info.Interceptors
			if h, ok := info.Impl.(interface {
				Interceptors() []interceptor.Interceptor
			}); ok {
				interceptors = append(interceptors, h.Interceptors()...)
			}

			return hTTPController_local_stub{
				impl:        info.Impl.(HTTPController),
				interceptor: interceptor.Chain(interceptors),
				name:        info.Name,
			}
		},
	})
	kod.Register(&kod.Registration{
		Name:      "github.com/go-kod/kod/tests/case1/InterceptorRetry",
		Interface: reflect.TypeOf((*InterceptorRetry)(nil)).Elem(),
		Impl:      reflect.TypeOf(interceptorRetry{}),
		Refs:      ``,
		LocalStubFn: func(ctx context.Context, info *kod.LocalStubFnInfo) any {
			interceptors := info.Interceptors
			if h, ok := info.Impl.(interface {
				Interceptors() []interceptor.Interceptor
			}); ok {
				interceptors = append(interceptors, h.Interceptors()...)
			}

			return interceptorRetry_local_stub{
				impl:        info.Impl.(InterceptorRetry),
				interceptor: interceptor.Chain(interceptors),
				name:        info.Name,
			}
		},
	})
	kod.Register(&kod.Registration{
		Name:      "github.com/go-kod/kod/tests/case1/LazyInitComponent",
		Interface: reflect.TypeOf((*LazyInitComponent)(nil)).Elem(),
		Impl:      reflect.TypeOf(lazyInitComponent{}),
		Refs:      ``,
		LocalStubFn: func(ctx context.Context, info *kod.LocalStubFnInfo) any {
			interceptors := info.Interceptors
			if h, ok := info.Impl.(interface {
				Interceptors() []interceptor.Interceptor
			}); ok {
				interceptors = append(interceptors, h.Interceptors()...)
			}

			return lazyInitComponent_local_stub{
				impl:        info.Impl.(LazyInitComponent),
				interceptor: interceptor.Chain(interceptors),
				name:        info.Name,
			}
		},
	})
	kod.Register(&kod.Registration{
		Name:      "github.com/go-kod/kod/tests/case1/LazyInitImpl",
		Interface: reflect.TypeOf((*LazyInitImpl)(nil)).Elem(),
		Impl:      reflect.TypeOf(lazyInitImpl{}),
		Refs:      `⟦8e153348:KoDeDgE:github.com/go-kod/kod/tests/case1/LazyInitImpl→github.com/go-kod/kod/tests/case1/LazyInitComponent⟧`,
		LocalStubFn: func(ctx context.Context, info *kod.LocalStubFnInfo) any {
			interceptors := info.Interceptors
			if h, ok := info.Impl.(interface {
				Interceptors() []interceptor.Interceptor
			}); ok {
				interceptors = append(interceptors, h.Interceptors()...)
			}

			return lazyInitImpl_local_stub{
				impl:        info.Impl.(LazyInitImpl),
				interceptor: interceptor.Chain(interceptors),
				name:        info.Name,
			}
		},
	})
	kod.Register(&kod.Registration{
		Name:      "github.com/go-kod/kod/Main",
		Interface: reflect.TypeOf((*kod.Main)(nil)).Elem(),
		Impl:      reflect.TypeOf(App{}),
		Refs:      `⟦d40a644a:KoDeDgE:github.com/go-kod/kod/Main→github.com/go-kod/kod/tests/case1/Test1Component⟧`,
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
	kod.Register(&kod.Registration{
		Name:      "github.com/go-kod/kod/tests/case1/Test1Component",
		Interface: reflect.TypeOf((*Test1Component)(nil)).Elem(),
		Impl:      reflect.TypeOf(test1Component{}),
		Refs:      ``,
		LocalStubFn: func(ctx context.Context, info *kod.LocalStubFnInfo) any {
			interceptors := info.Interceptors
			if h, ok := info.Impl.(interface {
				Interceptors() []interceptor.Interceptor
			}); ok {
				interceptors = append(interceptors, h.Interceptors()...)
			}

			return test1Component_local_stub{
				impl:        info.Impl.(Test1Component),
				interceptor: interceptor.Chain(interceptors),
				name:        info.Name,
			}
		},
	})
	kod.Register(&kod.Registration{
		Name:      "github.com/go-kod/kod/tests/case1/Test2Component",
		Interface: reflect.TypeOf((*Test2Component)(nil)).Elem(),
		Impl:      reflect.TypeOf(test2Component{}),
		Refs:      ``,
		LocalStubFn: func(ctx context.Context, info *kod.LocalStubFnInfo) any {
			interceptors := info.Interceptors
			if h, ok := info.Impl.(interface {
				Interceptors() []interceptor.Interceptor
			}); ok {
				interceptors = append(interceptors, h.Interceptors()...)
			}

			return test2Component_local_stub{
				impl:        info.Impl.(Test2Component),
				interceptor: interceptor.Chain(interceptors),
				name:        info.Name,
			}
		},
	})
	kod.Register(&kod.Registration{
		Name:      "github.com/go-kod/kod/tests/case1/ctxInterface",
		Interface: reflect.TypeOf((*ctxInterface)(nil)).Elem(),
		Impl:      reflect.TypeOf(ctxImpl{}),
		Refs:      ``,
		LocalStubFn: func(ctx context.Context, info *kod.LocalStubFnInfo) any {
			interceptors := info.Interceptors
			if h, ok := info.Impl.(interface {
				Interceptors() []interceptor.Interceptor
			}); ok {
				interceptors = append(interceptors, h.Interceptors()...)
			}

			return ctxInterface_local_stub{
				impl:        info.Impl.(ctxInterface),
				interceptor: interceptor.Chain(interceptors),
				name:        info.Name,
			}
		},
	})
	kod.Register(&kod.Registration{
		Name:      "github.com/go-kod/kod/tests/case1/panicCaseInterface",
		Interface: reflect.TypeOf((*panicCaseInterface)(nil)).Elem(),
		Impl:      reflect.TypeOf(panicCase{}),
		Refs:      ``,
		LocalStubFn: func(ctx context.Context, info *kod.LocalStubFnInfo) any {
			interceptors := info.Interceptors
			if h, ok := info.Impl.(interface {
				Interceptors() []interceptor.Interceptor
			}); ok {
				interceptors = append(interceptors, h.Interceptors()...)
			}

			return panicCaseInterface_local_stub{
				impl:        info.Impl.(panicCaseInterface),
				interceptor: interceptor.Chain(interceptors),
				name:        info.Name,
			}
		},
	})
	kod.Register(&kod.Registration{
		Name:      "github.com/go-kod/kod/tests/case1/panicNoRecvoeryCaseInterface",
		Interface: reflect.TypeOf((*panicNoRecvoeryCaseInterface)(nil)).Elem(),
		Impl:      reflect.TypeOf(panicNoRecvoeryCase{}),
		Refs:      ``,
		LocalStubFn: func(ctx context.Context, info *kod.LocalStubFnInfo) any {
			interceptors := info.Interceptors
			if h, ok := info.Impl.(interface {
				Interceptors() []interceptor.Interceptor
			}); ok {
				interceptors = append(interceptors, h.Interceptors()...)
			}

			return panicNoRecvoeryCaseInterface_local_stub{
				impl:        info.Impl.(panicNoRecvoeryCaseInterface),
				interceptor: interceptor.Chain(interceptors),
				name:        info.Name,
			}
		},
	})
	kod.Register(&kod.Registration{
		Name:      "github.com/go-kod/kod/tests/case1/test1Controller",
		Interface: reflect.TypeOf((*test1Controller)(nil)).Elem(),
		Impl:      reflect.TypeOf(test1ControllerImpl{}),
		Refs:      `⟦dd37e4d0:KoDeDgE:github.com/go-kod/kod/tests/case1/test1Controller→github.com/go-kod/kod/tests/case1/Test1Component⟧`,
		LocalStubFn: func(ctx context.Context, info *kod.LocalStubFnInfo) any {
			interceptors := info.Interceptors
			if h, ok := info.Impl.(interface {
				Interceptors() []interceptor.Interceptor
			}); ok {
				interceptors = append(interceptors, h.Interceptors()...)
			}

			return test1Controller_local_stub{
				impl:        info.Impl.(test1Controller),
				interceptor: interceptor.Chain(interceptors),
				name:        info.Name,
			}
		},
	})
	kod.Register(&kod.Registration{
		Name:      "github.com/go-kod/kod/tests/case1/testEchoController",
		Interface: reflect.TypeOf((*testEchoController)(nil)).Elem(),
		Impl:      reflect.TypeOf(testEchoControllerImpl{}),
		Refs:      ``,
		LocalStubFn: func(ctx context.Context, info *kod.LocalStubFnInfo) any {
			interceptors := info.Interceptors
			if h, ok := info.Impl.(interface {
				Interceptors() []interceptor.Interceptor
			}); ok {
				interceptors = append(interceptors, h.Interceptors()...)
			}

			return testEchoController_local_stub{
				impl:        info.Impl.(testEchoController),
				interceptor: interceptor.Chain(interceptors),
				name:        info.Name,
			}
		},
	})
	kod.Register(&kod.Registration{
		Name:      "github.com/go-kod/kod/tests/case1/testGinController",
		Interface: reflect.TypeOf((*testGinController)(nil)).Elem(),
		Impl:      reflect.TypeOf(testGinControllerImpl{}),
		Refs:      ``,
		LocalStubFn: func(ctx context.Context, info *kod.LocalStubFnInfo) any {
			interceptors := info.Interceptors
			if h, ok := info.Impl.(interface {
				Interceptors() []interceptor.Interceptor
			}); ok {
				interceptors = append(interceptors, h.Interceptors()...)
			}

			return testGinController_local_stub{
				impl:        info.Impl.(testGinController),
				interceptor: interceptor.Chain(interceptors),
				name:        info.Name,
			}
		},
	})
	kod.Register(&kod.Registration{
		Name:      "github.com/go-kod/kod/tests/case1/testRepository",
		Interface: reflect.TypeOf((*testRepository)(nil)).Elem(),
		Impl:      reflect.TypeOf(modelImpl{}),
		Refs:      ``,
		LocalStubFn: func(ctx context.Context, info *kod.LocalStubFnInfo) any {
			interceptors := info.Interceptors
			if h, ok := info.Impl.(interface {
				Interceptors() []interceptor.Interceptor
			}); ok {
				interceptors = append(interceptors, h.Interceptors()...)
			}

			return testRepository_local_stub{
				impl:        info.Impl.(testRepository),
				interceptor: interceptor.Chain(interceptors),
				name:        info.Name,
			}
		},
	})
	kod.Register(&kod.Registration{
		Name:      "github.com/go-kod/kod/tests/case1/testService",
		Interface: reflect.TypeOf((*testService)(nil)).Elem(),
		Impl:      reflect.TypeOf(serviceImpl{}),
		Refs:      ``,
		LocalStubFn: func(ctx context.Context, info *kod.LocalStubFnInfo) any {
			interceptors := info.Interceptors
			if h, ok := info.Impl.(interface {
				Interceptors() []interceptor.Interceptor
			}); ok {
				interceptors = append(interceptors, h.Interceptors()...)
			}

			return testService_local_stub{
				impl:        info.Impl.(testService),
				interceptor: interceptor.Chain(interceptors),
				name:        info.Name,
			}
		},
	})
}

// kod.InstanceOf checks.
var _ kod.InstanceOf[HTTPController] = (*httpControllerImpl)(nil)
var _ kod.InstanceOf[InterceptorRetry] = (*interceptorRetry)(nil)
var _ kod.InstanceOf[LazyInitComponent] = (*lazyInitComponent)(nil)
var _ kod.InstanceOf[LazyInitImpl] = (*lazyInitImpl)(nil)
var _ kod.InstanceOf[kod.Main] = (*App)(nil)
var _ kod.InstanceOf[Test1Component] = (*test1Component)(nil)
var _ kod.InstanceOf[Test2Component] = (*test2Component)(nil)
var _ kod.InstanceOf[ctxInterface] = (*ctxImpl)(nil)
var _ kod.InstanceOf[panicCaseInterface] = (*panicCase)(nil)
var _ kod.InstanceOf[panicNoRecvoeryCaseInterface] = (*panicNoRecvoeryCase)(nil)
var _ kod.InstanceOf[test1Controller] = (*test1ControllerImpl)(nil)
var _ kod.InstanceOf[testEchoController] = (*testEchoControllerImpl)(nil)
var _ kod.InstanceOf[testGinController] = (*testGinControllerImpl)(nil)
var _ kod.InstanceOf[testRepository] = (*modelImpl)(nil)
var _ kod.InstanceOf[testService] = (*serviceImpl)(nil)

// Local stub implementations.

type hTTPController_local_stub struct {
	impl        HTTPController
	name        string
	interceptor interceptor.Interceptor
}

// Check that hTTPController_local_stub implements the HTTPController interface.
var _ HTTPController = (*hTTPController_local_stub)(nil)

func (s hTTPController_local_stub) Foo(a0 http.ResponseWriter, a1 *http.Request) {
	// Because the first argument is not context.Context, so interceptors are not supported.
	s.impl.Foo(a0, a1)
	return
}

type interceptorRetry_local_stub struct {
	impl        InterceptorRetry
	name        string
	interceptor interceptor.Interceptor
}

// Check that interceptorRetry_local_stub implements the InterceptorRetry interface.
var _ InterceptorRetry = (*interceptorRetry_local_stub)(nil)

func (s interceptorRetry_local_stub) TestError(ctx context.Context) (err error) {

	if s.interceptor == nil {
		err = s.impl.TestError(ctx)
		return
	}

	call := func(ctx context.Context, info interceptor.CallInfo, req, res []any) (err error) {
		err = s.impl.TestError(ctx)
		return
	}

	info := interceptor.CallInfo{
		Impl:       s.impl,
		Component:  s.name,
		FullMethod: "github.com/go-kod/kod/tests/case1/InterceptorRetry.TestError",
		Method:     "TestError",
	}

	err = s.interceptor(ctx, info, []any{}, []any{}, call)
	return
}

func (s interceptorRetry_local_stub) TestNormal(ctx context.Context) (err error) {

	if s.interceptor == nil {
		err = s.impl.TestNormal(ctx)
		return
	}

	call := func(ctx context.Context, info interceptor.CallInfo, req, res []any) (err error) {
		err = s.impl.TestNormal(ctx)
		return
	}

	info := interceptor.CallInfo{
		Impl:       s.impl,
		Component:  s.name,
		FullMethod: "github.com/go-kod/kod/tests/case1/InterceptorRetry.TestNormal",
		Method:     "TestNormal",
	}

	err = s.interceptor(ctx, info, []any{}, []any{}, call)
	return
}

type lazyInitComponent_local_stub struct {
	impl        LazyInitComponent
	name        string
	interceptor interceptor.Interceptor
}

// Check that lazyInitComponent_local_stub implements the LazyInitComponent interface.
var _ LazyInitComponent = (*lazyInitComponent_local_stub)(nil)

func (s lazyInitComponent_local_stub) Try(ctx context.Context) (err error) {

	if s.interceptor == nil {
		err = s.impl.Try(ctx)
		return
	}

	call := func(ctx context.Context, info interceptor.CallInfo, req, res []any) (err error) {
		err = s.impl.Try(ctx)
		return
	}

	info := interceptor.CallInfo{
		Impl:       s.impl,
		Component:  s.name,
		FullMethod: "github.com/go-kod/kod/tests/case1/LazyInitComponent.Try",
		Method:     "Try",
	}

	err = s.interceptor(ctx, info, []any{}, []any{}, call)
	return
}

type lazyInitImpl_local_stub struct {
	impl        LazyInitImpl
	name        string
	interceptor interceptor.Interceptor
}

// Check that lazyInitImpl_local_stub implements the LazyInitImpl interface.
var _ LazyInitImpl = (*lazyInitImpl_local_stub)(nil)

func (s lazyInitImpl_local_stub) Try(ctx context.Context) {

	if s.interceptor == nil {
		s.impl.Try(ctx)
		return
	}

	call := func(ctx context.Context, info interceptor.CallInfo, req, res []any) (err error) {
		s.impl.Try(ctx)
		return
	}

	info := interceptor.CallInfo{
		Impl:       s.impl,
		Component:  s.name,
		FullMethod: "github.com/go-kod/kod/tests/case1/LazyInitImpl.Try",
		Method:     "Try",
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

type test1Component_local_stub struct {
	impl        Test1Component
	name        string
	interceptor interceptor.Interceptor
}

// Check that test1Component_local_stub implements the Test1Component interface.
var _ Test1Component = (*test1Component_local_stub)(nil)

func (s test1Component_local_stub) Foo(ctx context.Context, a1 *FooReq) (r0 *FooRes, err error) {

	if s.interceptor == nil {
		r0, err = s.impl.Foo(ctx, a1)
		return
	}

	call := func(ctx context.Context, info interceptor.CallInfo, req, res []any) (err error) {
		r0, err = s.impl.Foo(ctx, a1)
		res[0] = r0
		return
	}

	info := interceptor.CallInfo{
		Impl:       s.impl,
		Component:  s.name,
		FullMethod: "github.com/go-kod/kod/tests/case1/Test1Component.Foo",
		Method:     "Foo",
	}

	err = s.interceptor(ctx, info, []any{a1}, []any{r0}, call)
	return
}

type test2Component_local_stub struct {
	impl        Test2Component
	name        string
	interceptor interceptor.Interceptor
}

// Check that test2Component_local_stub implements the Test2Component interface.
var _ Test2Component = (*test2Component_local_stub)(nil)

func (s test2Component_local_stub) GetClient() (r0 *http.Client) {
	// Because the first argument is not context.Context, so interceptors are not supported.
	r0 = s.impl.GetClient()
	return
}

type ctxInterface_local_stub struct {
	impl        ctxInterface
	name        string
	interceptor interceptor.Interceptor
}

// Check that ctxInterface_local_stub implements the ctxInterface interface.
var _ ctxInterface = (*ctxInterface_local_stub)(nil)

func (s ctxInterface_local_stub) Foo(ctx context.Context) {

	if s.interceptor == nil {
		s.impl.Foo(ctx)
		return
	}

	call := func(ctx context.Context, info interceptor.CallInfo, req, res []any) (err error) {
		s.impl.Foo(ctx)
		return
	}

	info := interceptor.CallInfo{
		Impl:       s.impl,
		Component:  s.name,
		FullMethod: "github.com/go-kod/kod/tests/case1/ctxInterface.Foo",
		Method:     "Foo",
	}

	_ = s.interceptor(ctx, info, []any{}, []any{}, call)
}

type panicCaseInterface_local_stub struct {
	impl        panicCaseInterface
	name        string
	interceptor interceptor.Interceptor
}

// Check that panicCaseInterface_local_stub implements the panicCaseInterface interface.
var _ panicCaseInterface = (*panicCaseInterface_local_stub)(nil)

func (s panicCaseInterface_local_stub) TestPanic(ctx context.Context) {

	if s.interceptor == nil {
		s.impl.TestPanic(ctx)
		return
	}

	call := func(ctx context.Context, info interceptor.CallInfo, req, res []any) (err error) {
		s.impl.TestPanic(ctx)
		return
	}

	info := interceptor.CallInfo{
		Impl:       s.impl,
		Component:  s.name,
		FullMethod: "github.com/go-kod/kod/tests/case1/panicCaseInterface.TestPanic",
		Method:     "TestPanic",
	}

	_ = s.interceptor(ctx, info, []any{}, []any{}, call)
}

type panicNoRecvoeryCaseInterface_local_stub struct {
	impl        panicNoRecvoeryCaseInterface
	name        string
	interceptor interceptor.Interceptor
}

// Check that panicNoRecvoeryCaseInterface_local_stub implements the panicNoRecvoeryCaseInterface interface.
var _ panicNoRecvoeryCaseInterface = (*panicNoRecvoeryCaseInterface_local_stub)(nil)

func (s panicNoRecvoeryCaseInterface_local_stub) TestPanic(ctx context.Context) {

	if s.interceptor == nil {
		s.impl.TestPanic(ctx)
		return
	}

	call := func(ctx context.Context, info interceptor.CallInfo, req, res []any) (err error) {
		s.impl.TestPanic(ctx)
		return
	}

	info := interceptor.CallInfo{
		Impl:       s.impl,
		Component:  s.name,
		FullMethod: "github.com/go-kod/kod/tests/case1/panicNoRecvoeryCaseInterface.TestPanic",
		Method:     "TestPanic",
	}

	_ = s.interceptor(ctx, info, []any{}, []any{}, call)
}

type test1Controller_local_stub struct {
	impl        test1Controller
	name        string
	interceptor interceptor.Interceptor
}

// Check that test1Controller_local_stub implements the test1Controller interface.
var _ test1Controller = (*test1Controller_local_stub)(nil)

type testEchoController_local_stub struct {
	impl        testEchoController
	name        string
	interceptor interceptor.Interceptor
}

// Check that testEchoController_local_stub implements the testEchoController interface.
var _ testEchoController = (*testEchoController_local_stub)(nil)

func (s testEchoController_local_stub) Error(a0 echo.Context) (err error) {
	// Because the first argument is not context.Context, so interceptors are not supported.
	err = s.impl.Error(a0)
	return
}

func (s testEchoController_local_stub) Hello(a0 echo.Context) (err error) {
	// Because the first argument is not context.Context, so interceptors are not supported.
	err = s.impl.Hello(a0)
	return
}

type testGinController_local_stub struct {
	impl        testGinController
	name        string
	interceptor interceptor.Interceptor
}

// Check that testGinController_local_stub implements the testGinController interface.
var _ testGinController = (*testGinController_local_stub)(nil)

func (s testGinController_local_stub) Hello(a0 *gin.Context) {
	// Because the first argument is not context.Context, so interceptors are not supported.
	s.impl.Hello(a0)
	return
}

type testRepository_local_stub struct {
	impl        testRepository
	name        string
	interceptor interceptor.Interceptor
}

// Check that testRepository_local_stub implements the testRepository interface.
var _ testRepository = (*testRepository_local_stub)(nil)

func (s testRepository_local_stub) Foo(ctx context.Context) (err error) {

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
		Component:  s.name,
		FullMethod: "github.com/go-kod/kod/tests/case1/testRepository.Foo",
		Method:     "Foo",
	}

	err = s.interceptor(ctx, info, []any{}, []any{}, call)
	return
}

type testService_local_stub struct {
	impl        testService
	name        string
	interceptor interceptor.Interceptor
}

// Check that testService_local_stub implements the testService interface.
var _ testService = (*testService_local_stub)(nil)

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
		Component:  s.name,
		FullMethod: "github.com/go-kod/kod/tests/case1/testService.Foo",
		Method:     "Foo",
	}

	err = s.interceptor(ctx, info, []any{}, []any{}, call)
	return
}
