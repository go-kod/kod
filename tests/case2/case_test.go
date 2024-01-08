package case2

import (
	"context"
	"testing"

	"github.com/go-kod/kod"
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
		})
		if err.Error() != "components [github.com/go-kod/kod/tests/case2/Test2Component] and [github.com/go-kod/kod/tests/case2/Test1Component] have cycle Ref" {
			panic(err)
		}
	})
}
