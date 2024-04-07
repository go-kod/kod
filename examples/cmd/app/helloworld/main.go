package main

import (
	"context"
	"fmt"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/examples/app/helloworld"
	"github.com/go-kod/kod/interceptor/kaccesslog"
	"github.com/go-kod/kod/interceptor/kmetric"
	"github.com/go-kod/kod/interceptor/krecovery"
	"github.com/go-kod/kod/interceptor/ktrace"
)

type app struct {
	kod.Implements[kod.Main]
	helloWorld kod.Ref[helloworld.HelloWorld]
}

func main() {
	kod.Run(context.Background(), func(ctx context.Context, main *app) error {
		fmt.Println(main.helloWorld.Get().SayHello(ctx))
		return nil
	}, kod.WithInterceptors(
		ktrace.Interceptor(),
		kmetric.Interceptor(),
		krecovery.Interceptor(),
		kaccesslog.Interceptor(),
	))
}
