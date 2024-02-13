package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-kod/kod"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type test1Config struct {
	A     string
	Redis struct {
		Addr    string
		Timeout time.Duration
	}
}

type test1ControllerImpl struct {
	kod.Implements[test1Controller]

	httpControllerImpl kod.Ref[HTTPController]
	test1Component     kod.Ref[Test1Component]
}

func (t *test1ControllerImpl) Foo(cccccc *gin.Context) {
	_ = t.test1Component.Get().Foo(cccccc, &FooReq{})
}

type httpControllerImpl struct {
	kod.Implements[HTTPController]
	kod.Ref[testService]
}

func (t *httpControllerImpl) Foo(w http.ResponseWriter, r http.Request) {
}

type serviceImpl struct {
	kod.Implements[testService]
	kod.Ref[testModel]
}

func (t *serviceImpl) Foo(ctx context.Context) error {
	return nil
}

type modelImpl struct {
	kod.Implements[testModel]
}

func (t *modelImpl) Foo(ctx context.Context) error {
	return nil
}

type test1Component struct {
	kod.Implements[Test1Component]
	kod.WithConfig[test1Config]
}

func (t *test1Component) Init(ctx context.Context) error {
	return nil
}

func (t *test1Component) Stop(ctx context.Context) error {
	return nil
}

type FooReq struct {
	Id int
}

func (t *test1Component) Foo(ctx context.Context, req *FooReq) error {
	t.L(ctx).InfoContext(ctx, "Foo info ", "config", t.Config())
	t.L(ctx).Debug("Foo debug:")
	fmt.Println(errors.New("test1"))
	return errors.New("test1:" + t.Config().A)
}

type fakeTest1Component struct {
	A string
}

func (f *fakeTest1Component) Foo(ctx context.Context, req *FooReq) error {
	fmt.Println(f.A)
	return errors.New("A:" + f.A)
}

type App struct {
	kod.Implements[kod.Main]
	test1ControllerImpl kod.Ref[test1Controller]
	test1               kod.Ref[Test1Component]
}

func (app *App) Run(ctx context.Context) error {
	ctx, span := otel.Tracer("").Start(ctx, "Run", trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	return app.test1.Get().Foo(ctx, &FooReq{0})
}

func Run(ctx context.Context, app *App) error {
	ctx, span := otel.Tracer("").Start(ctx, "Run", trace.WithSpanKind(trace.SpanKindInternal))
	defer span.End()

	return app.test1.Get().Foo(ctx, &FooReq{0})
}

func main() {
	Run(context.TODO(), new(App))
}
