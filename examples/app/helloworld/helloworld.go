package helloworld

import (
	"context"

	"github.com/go-kod/kod"
)

type helloWorld struct {
	kod.Implements[HelloWorld]
}

func (h *helloWorld) SayHello(ctx context.Context) string {
	return "Hello, World!"
}
