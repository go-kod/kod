package case1

import (
	"context"
	"fmt"
	"testing"

	"github.com/go-kod/kod"
	"github.com/stretchr/testify/require"
)

func TestTest(t *testing.T) {
	t.Parallel()
	kod.RunTest(t, func(ctx context.Context, k *test1Component) {
		_, err := k.Foo(ctx, &FooReq{})
		fmt.Println(err)
		require.Equal(t, "test1:B", err.Error())
	})
}

func TestTest2(t *testing.T) {
	t.Parallel()
	kod.RunTest2(t, func(ctx context.Context, k *test1Component, k2 Test2Component) {
		_, err := k.Foo(ctx, &FooReq{})
		fmt.Println(err)
		require.Equal(t, "test1:B", err.Error())
	})
}

func TestTest3(t *testing.T) {
	t.Parallel()

	require.Panics(t, func() {
		kod.RunTest3(t, func(ctx context.Context, k *test1Component, k2 panicNoRecvoeryCaseInterface, k3 test1Controller) {
			_, err := k.Foo(ctx, &FooReq{})
			fmt.Println(err)
			require.Equal(t, "test1:B", err.Error())

			k2.TestPanic(ctx)
		})
	})
}
