package case4

import (
	"context"
	"errors"

	"github.com/go-kod/kod"
)

type test1Component struct {
	kod.Implements[Test1Component]
	// nolint
	test2 kod.Ref[Test3Component]
}

type FooReq struct {
	Id int
}

func (t *test1Component) Foo(ctx context.Context, req *FooReq) error {
	return errors.New("test1")
}

type test2Component struct {
	kod.Implements[Test2Component]
	// nolint
	test1 kod.Ref[Test3Component]
}

func (t *test2Component) Foo(ctx context.Context, req *FooReq) error {
	return errors.New("test2")
}

type test3Component struct {
	kod.Implements[Test3Component]
}

func (t *test3Component) Foo(ctx context.Context, req *FooReq) error {
	return errors.New("test3")
}

type App struct {
	kod.Implements[kod.Main]
	test1 kod.Ref[Test1Component]
	_     kod.Ref[Test2Component]
}

func (app *App) Run(ctx context.Context) error {
	return app.test1.Get().Foo(ctx, &FooReq{})
}
