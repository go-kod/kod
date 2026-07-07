package case1

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-kod/kod"
)

func TestTest(t *testing.T) {
	t.Parallel()

	kod.RunTest(t, func(ctx context.Context, k *test1Component) {
		got, err := kod.Get[*test1Component](ctx)
		require.NoError(t, err)
		require.Same(t, k, got)

		_, err = k.Foo(ctx, &FooReq{})
		fmt.Println(err)
		require.Equal(t, "test1:B", err.Error())
	})
}

func TestTest2(t *testing.T) {
	t.Parallel()

	kod.RunTest(t, func(ctx context.Context, k *test1Component, k2 Test2Component) {
		got, err := kod.Get[Test2Component](ctx)
		require.NoError(t, err)
		require.IsType(t, k2, got)
		require.NotNil(t, got.GetClient())

		_, err = k.Foo(ctx, &FooReq{})
		fmt.Println(err)
		require.Equal(t, "test1:B", err.Error())
	})
}

func TestInterfaceThenImpl(t *testing.T) {
	t.Parallel()

	kod.RunTest(t, func(ctx context.Context, intf Test1Component, impl *test1Component) {
		require.NotNil(t, intf)
		require.NotNil(t, impl)
	})
}

func TestTest3(t *testing.T) {
	t.Parallel()

	require.Panics(t, func() {
		kod.RunTest(t, func(ctx context.Context, k *test1Component, k2 panicNoRecvoeryCaseInterface, k3 test1Controller) {
			_, err := k.Foo(ctx, &FooReq{})
			fmt.Println(err)
			require.Equal(t, "test1:B", err.Error())

			k2.TestPanic(ctx)
		})
	})
}
