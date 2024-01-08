package case1

import (
	"context"
	"testing"

	"github.com/go-kod/kod"
)

func TestPanicRecovery(t *testing.T) {
	t.Parallel()

	kod.RunTest(t, func(ctx context.Context, t panicCaseInterface) {
		t.TestPanic(nil)
	})
}
