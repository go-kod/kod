package helloworld

import (
	"context"
	"fmt"

	"github.com/go-kod/kod"
)

type App struct {
	kod.Implements[kod.Main]
	kod.WithGlobalConfig[GlobalConfig]

	HelloWorld kod.Ref[HelloWorld]
	HelloBob   kod.Ref[HelloBob]
}

type GlobalConfig struct {
	Name string
}

type Config struct {
	Name string
}

type helloWorld struct {
	kod.Implements[HelloWorld]
	kod.WithConfig[Config]
}

func (h *helloWorld) Init(ctx context.Context) error {
	fmt.Println("helloWorld init")
	return nil
}

func (h *helloWorld) SayHello(ctx context.Context) {
	fmt.Println("Hello, World!" + h.Config().Name)
}

func (h *helloWorld) Shutdown(ctx context.Context) error {
	fmt.Println("helloWorld shutdown")
	return nil
}

type lazyHelloBob struct {
	kod.Implements[HelloBob]
	kod.LazyInit
}

func (h *lazyHelloBob) Init(ctx context.Context) error {
	fmt.Println("lazyHelloBob init")
	return nil
}

func (h *lazyHelloBob) SayHello(ctx context.Context) {
	fmt.Println("Hello, Bob!")
}

func (h *lazyHelloBob) Shutdown(ctx context.Context) error {
	fmt.Println("lazyHelloBob shutdown")
	return nil
}
