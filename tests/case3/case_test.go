package case3

import (
	"context"
	"testing"

	"github.com/go-kod/kod"
)

func TestRun(t *testing.T) {
	t.Parallel()
	t.Run("Main依赖1，2依赖3，3依赖1，有循环依赖", func(t *testing.T) {
		err := kod.Run(context.Background(), func(ctx context.Context, t *App) error {
			return t.Run(ctx)
		})
		if err.Error() != "components [github.com/go-kod/kod/tests/case3/Test3Component] and [github.com/go-kod/kod/tests/case3/Test1Component] have cycle Ref" {
			panic(err)
		}
	})
}
