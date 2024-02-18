package main

import (
	"context"

	"github.com/go-kod/kod"
)

type app struct {
	kod.Implements[kod.Main]
	kod.WithConfig[config]
}

type config struct {
	Name string
}

func main() {
	kod.Run(context.Background(), func(ctx context.Context, app *app) error {
		app.L(ctx).Info("Hello, " + app.Config().Name)
		return nil
	})
}
