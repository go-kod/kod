package case1

import (
	"context"
	"testing"

	"github.com/go-kod/kod"
)

func TestCtxImpl(t *testing.T) {
	t.Parallel()
	kod.RunTest(t, func(ctx context.Context, k ctxInterface) {
		k.Foo(ctx)
	})
}
