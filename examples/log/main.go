package main

import (
	"context"
	"log/slog"

	"github.com/go-kod/kod"
	"go.opentelemetry.io/otel/exporters/stdout/stdoutlog"
	"go.opentelemetry.io/otel/log/global"
	"go.opentelemetry.io/otel/sdk/log"
)

type app struct {
	kod.Implements[kod.Main]
}

func main() {
	exp, err := stdoutlog.New()
	if err != nil {
		panic(err)
	}

	processor := log.NewSimpleProcessor(exp)
	provider := log.NewLoggerProvider(log.WithProcessor(processor))
	defer func() {
		if err := provider.Shutdown(context.Background()); err != nil {
			panic(err)
		}
	}()

	global.SetLoggerProvider(provider)

	kod.Run(context.Background(), func(ctx context.Context, app *app) error {
		// default log level is info
		app.L(ctx).Info("hello world info")

		// wont print
		app.L(ctx).Debug("hello world info1")
		// wont print
		app.L(ctx).Info("hello world info1")
		// will print
		app.L(ctx).Error("hello world error", slog.String("error", "error message"))

		return nil
	})
}
