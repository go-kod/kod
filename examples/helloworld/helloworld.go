package helloworld

import (
	"context"
	"fmt"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/interceptor"
)

type App struct {
	kod.Implements[kod.Main]

	HelloWorld            kod.Ref[HelloWorld]
	HelloWorldLazy        kod.Ref[HelloWorldLazy]
	HelloWorldInterceptor kod.Ref[HelloWorldInterceptor]
}

// HelloWorld ...
type helloWorld struct {
	kod.Implements[HelloWorld]
}

func (h *helloWorld) Init(_ context.Context) error {
	fmt.Println("helloWorld init")
	return nil
}

// SayHello ...
// line two
func (h *helloWorld) SayHello(ctx context.Context) {
	h.L(ctx).Info("Hello, World!")

	fmt.Println("Hello, World!")
}

func (h *helloWorld) Shutdown(_ context.Context) error {
	fmt.Println("helloWorld shutdown")
	return nil
}

type lazyHelloWorld struct {
	kod.Implements[HelloWorldLazy]
	kod.LazyInit
}

func (h *lazyHelloWorld) Init(_ context.Context) error {
	fmt.Println("lazyHelloWorld init")
	return nil
}

// SayHello ...
func (h *lazyHelloWorld) SayHello(_ context.Context) {
	fmt.Println("Hello, Lazy!")
}

func (h *lazyHelloWorld) Shutdown(_ context.Context) error {
	fmt.Println("lazyHelloWorld shutdown")
	return nil
}

// helloWorldInterceptor ...
type helloWorldInterceptor struct {
	kod.Implements[HelloWorldInterceptor]
}

// SayHello ...
func (h *helloWorldInterceptor) SayHello(_ context.Context) {
	fmt.Println("Hello, Interceptor!")
}

func (h *helloWorldInterceptor) Interceptors() []interceptor.Interceptor {
	return []interceptor.Interceptor{
		func(ctx context.Context, info interceptor.CallInfo, req, reply []any, invoker interceptor.HandleFunc) error {
			fmt.Println("Before call")
			err := invoker(ctx, info, req, reply)
			fmt.Println("After call")
			return err
		},
	}
}
