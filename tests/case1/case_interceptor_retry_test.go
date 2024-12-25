package case1

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/go-kod/kod"
)

func TestInterceptorRetry(t *testing.T) {
	kod.RunTest(t, func(ctx context.Context, k InterceptorRetry) {
		require.ErrorContains(t, k.TestError(ctx), "retry fail")
	})
}

func TestInterceptorRetry1(t *testing.T) {
	kod.RunTest(t, func(ctx context.Context, k InterceptorRetry) {
		require.Nil(t, k.TestNormal(ctx))
	})
}
