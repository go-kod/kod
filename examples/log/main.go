package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/go-kod/kod"
)

type app struct {
	kod.Implements[kod.Main]
}

func main() {
	os.Setenv("OTEL_LOGS_EXPORTER", "console")

	kod.Run(context.Background(), func(ctx context.Context, app *app) error {
		// default log level is debug
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
