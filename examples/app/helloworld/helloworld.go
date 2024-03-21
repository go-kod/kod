package helloworld

import (
	"context"
	"fmt"

	"github.com/go-kod/kod"
)

type helloWorld struct {
	kod.Implements[HelloWorld]
}

func (h *helloWorld) SayHello() string {
	return "Hello, World!"
}

type app struct {
	kod.Implements[kod.Main]
	helloWorld kod.Ref[HelloWorld]
}

func main() {
	kod.Run(context.Background(), func(ctx context.Context, main *app) error {
		fmt.Println(main.helloWorld.Get().SayHello())
		return nil
	})
}
