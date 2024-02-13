package case2

import (
	"context"
	"errors"

	"github.com/go-kod/kod"
)

type test1Component struct {
	kod.Implements[Test1Component]
	// nolint
	test2 kod.Ref[Test2Component]
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
	test1 kod.Ref[Test1Component]
}

func (t *test2Component) Foo(ctx context.Context, req *FooReq) error {
	return errors.New("test2")
}

type App struct {
	kod.Implements[kod.Main]
	test1 kod.Ref[Test1Component]
}

func (app *App) Run(ctx context.Context) error {
	return app.test1.Get().Foo(ctx, &FooReq{})
}

func Run(ctx context.Context, app *App) error {
	return app.test1.Get().Foo(ctx, &FooReq{})
}
