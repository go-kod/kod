package main

import (
	"context"
	"fmt"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/examples/app/helloworld"
)

type app struct {
	kod.Implements[kod.Main]
	helloWorld kod.Ref[helloworld.HelloWorld]
}

func main() {
	kod.Run(context.Background(), func(ctx context.Context, main *app) error {
		fmt.Println(main.helloWorld.Get().SayHello())
		return nil
	})
}
