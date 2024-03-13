package main

import (
	"context"
	"log/slog"

	"github.com/go-kod/kod"
)

type app struct {
	kod.Implements[kod.Main]
}

func main() {
	kod.Run(context.Background(), func(ctx context.Context, app *app) error {
		// default log level is info
		app.L(ctx).Info("hello world info")

		// set log level to error
		kod.FromContext(ctx).LevelVar().Set(slog.LevelError)

		// wont print
		app.L(ctx).Debug("hello world info1")
		// wont print
		app.L(ctx).Info("hello world info1")
		// will print
		app.L(ctx).Error("hello world error", slog.String("error", "error message"))

		return nil
	})
}
