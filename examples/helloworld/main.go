package main

import (
	"context"

	"github.com/go-kod/kod"
)

type helloworld struct {
	kod.Implements[Helloworld]
}

func (h *helloworld) SayHello() string {
	return "Hello, World!"
}

type app struct {
	kod.Implements[kod.Main]
	helloworld kod.Ref[Helloworld]
}

func main() {
	kod.Run(context.Background(), func(ctx context.Context, main *app) error {
		main.helloworld.Get().SayHello()
		return nil
	})
}
