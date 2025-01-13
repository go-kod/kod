package case1

import (
	"context"
	"testing"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/interceptor/krecovery"
)

func TestPanicRecovery(t *testing.T) {
	t.Parallel()

	kod.RunTest(t, func(ctx context.Context, t panicCaseInterface) {
		t.TestPanic(ctx)
	})
}

func TestRunWithInterceptor(t *testing.T) {
	t.Parallel()

	t.Run("panicNoRecvoeryCase with interceptor", func(t *testing.T) {
		kod.RunTest(t, func(ctx context.Context, t panicNoRecvoeryCaseInterface) {
			kod.FromContext(ctx).SetInterceptors(krecovery.Interceptor())

			t.TestPanic(ctx)
		})
	})
}
