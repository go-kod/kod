package helloworld

import (
	"context"
	"fmt"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/interceptor"
)

type App struct {
	kod.Implements[kod.Main]
	kod.WithGlobalConfig[GlobalConfig]

	HelloWorld            kod.Ref[HelloWorld]
	HelloWorldLazy        kod.Ref[HelloWorldLazy]
	HelloWorldInterceptor kod.Ref[HelloWorldInterceptor]
}

type GlobalConfig struct {
	Name string `default:"kod"`
}

type Config struct {
	Name string `default:"-"`
}

// HelloWorld ...
type helloWorld struct {
	kod.Implements[HelloWorld]
	kod.WithConfig[Config]
}

func (h *helloWorld) Init(ctx context.Context) error {
	fmt.Println("helloWorld init")
	return nil
}

// SayHello ...
// line two
func (h *helloWorld) SayHello(ctx context.Context) {
	h.L(ctx).Info("Hello, World!")

	fmt.Println("Hello, World!" + h.Config().Name)
}

func (h *helloWorld) Shutdown(ctx context.Context) error {
	fmt.Println("helloWorld shutdown")
	return nil
}

type lazyHelloWorld struct {
	kod.Implements[HelloWorldLazy]
	kod.LazyInit
}

func (h *lazyHelloWorld) Init(ctx context.Context) error {
	fmt.Println("lazyHelloWorld init")
	return nil
}

// SayHello ...
func (h *lazyHelloWorld) SayHello(ctx context.Context) {
	fmt.Println("Hello, Lazy!")
}

func (h *lazyHelloWorld) Shutdown(ctx context.Context) error {
	fmt.Println("lazyHelloWorld shutdown")
	return nil
}

// helloWorldInterceptor ...
type helloWorldInterceptor struct {
	kod.Implements[HelloWorldInterceptor]
}

// SayHello ...
func (h *helloWorldInterceptor) SayHello(ctx context.Context) {
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
