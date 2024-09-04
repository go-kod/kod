package helloworld

import (
	"context"
	"fmt"

	"github.com/go-kod/kod"
)

type App struct {
	kod.Implements[kod.Main]
	kod.WithGlobalConfig[Config]

	HelloWorld kod.Ref[HelloWorld]
}

type helloWorld struct {
	kod.Implements[HelloWorld]
	kod.WithConfig[Config]
}

func (h *helloWorld) SayHello(ctx context.Context) {
	fmt.Println("Hello, World!" + h.Config().Name)
}
