package case4

import (
	"context"
	"testing"

	"github.com/go-kod/kod"
)

func TestRun(t *testing.T) {
	t.Parallel()

	t.Run("Main依赖1和2，1依赖3，2也依赖3，没有循环依赖", func(t *testing.T) {
		err := kod.Run(context.Background(), func(ctx context.Context, t *App) error {
			return t.Run(ctx)
		})
		if err.Error() != "test1" {
			panic(err)
		}
	})
}
