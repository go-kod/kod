package case2

import (
	"context"
	"reflect"
	"testing"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/interceptor"
)

func TestRun(t *testing.T) {
	t.Parallel()
	t.Run("case1", func(t *testing.T) {
		err := kod.Run(context.Background(), func(ctx context.Context, t *App) error {
			return t.Run(ctx)
		})
		if err.Error() != "components [github.com/go-kod/kod/tests/case2/Test2Component] and [github.com/go-kod/kod/tests/case2/Test1Component] have cycle Ref" {
			panic(err)
		}
	})

	t.Run("case2", func(t *testing.T) {
		err := kod.Run(context.Background(), func(ctx context.Context, t *App) error {
			return t.Run(ctx)
		}, kod.WithRegistrations(
			&kod.Registration{
				Name:      "github.com/go-kod/kod/Main",
				Interface: reflect.TypeOf((*kod.Main)(nil)).Elem(),
				Impl:      reflect.TypeOf(App{}),
				Refs:      `⟦73dc6a0b:KoDeDgE:github.com/go-kod/kod/Main→github.com/go-kod/kod/tests/case2/Test1Component⟧`,
				LocalStubFn: func(ctx context.Context, info *kod.LocalStubFnInfo) any {
					var interceptors []interceptor.Interceptor
					if h, ok := info.Impl.(interface {
						Interceptors() []interceptor.Interceptor
					}); ok {
						interceptors = h.Interceptors()
					}

					return main_local_stub{
						impl:        info.Impl.(kod.Main),
						interceptor: interceptor.Chain(interceptors),
						name:        info.Name,
					}
				},
			},
			&kod.Registration{
				Name:      "github.com/go-kod/kod/tests/case2/Test1Component",
				Interface: reflect.TypeOf((*Test1Component)(nil)).Elem(),
				Impl:      reflect.TypeOf(test1Component{}),
				Refs:      `⟦3dc9f060:KoDeDgE:github.com/go-kod/kod/tests/case2/Test1Component→github.com/go-kod/kod/tests/case2/Test2Component⟧`,
				LocalStubFn: func(ctx context.Context, info *kod.LocalStubFnInfo) any {
					var interceptors []interceptor.Interceptor
					if h, ok := info.Impl.(interface {
						Interceptors() []interceptor.Interceptor
					}); ok {
						interceptors = h.Interceptors()
					}

					return test1Component_local_stub{
						impl:        info.Impl.(Test1Component),
						interceptor: interceptor.Chain(interceptors),
						name:        info.Name,
					}
				},
			},
			&kod.Registration{
				Name:      "github.com/go-kod/kod/tests/case2/Test2Component",
				Interface: reflect.TypeOf((*Test2Component)(nil)).Elem(),
				Impl:      reflect.TypeOf(test2Component{}),
				Refs:      `⟦1767cee9:KoDeDgE:github.com/go-kod/kod/tests/case2/Test2Component→github.com/go-kod/kod/tests/case2/Test1Component⟧`,
				LocalStubFn: func(ctx context.Context, info *kod.LocalStubFnInfo) any {
					var interceptors []interceptor.Interceptor
					if h, ok := info.Impl.(interface {
						Interceptors() []interceptor.Interceptor
					}); ok {
						interceptors = h.Interceptors()
					}

					return test2Component_local_stub{
						impl:        info.Impl.(Test2Component),
						interceptor: interceptor.Chain(interceptors),
						name:        info.Name,
					}
				},
			},
		))
		if err.Error() != "components [github.com/go-kod/kod/tests/case2/Test2Component] and [github.com/go-kod/kod/tests/case2/Test1Component] have cycle Ref" {
			panic(err)
		}
	})
}
