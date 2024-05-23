package main

import (
	"context"
	"log/slog"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/ext/client/kuptrace"
)

type app struct {
	kod.Implements[kod.Main]
}

func main() {
	kod.Run(context.Background(), func(ctx context.Context, app *app) error {
		s, _ := kuptrace.Config{DSN: "http://project2_secret_token@localhost:14318?grpc=14317", Debug: true}.Build(ctx)
		defer s.Stop(ctx)

		// default log level is info
		app.L(ctx).InfoContext(ctx, "hello world info")

		// wont print
		app.L(ctx).Debug("hello world info1")
		// wont print
		app.L(ctx).Info("hello world info1")
		// will print
		app.L(ctx).Error("hello world error", slog.String("error", "error message"))

		return nil
	})
}
