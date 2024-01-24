package case1

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/interceptor/kaccesslog"
	"github.com/go-kod/kod/interceptor/kcircuitbreaker"
	"github.com/go-kod/kod/interceptor/kmetric"
	"github.com/go-kod/kod/interceptor/kratelimit"
	"github.com/go-kod/kod/interceptor/krecovery"
	"github.com/go-kod/kod/interceptor/ktimeout"
	"github.com/go-kod/kod/interceptor/ktrace"
	"github.com/go-kod/kod/interceptor/kvalidate"
	"github.com/samber/lo"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
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

	test1Component kod.Ref[Test1Component]
}

type serviceImpl struct {
	kod.Implements[testService]
}

func (t *serviceImpl) Foo(ctx context.Context) error {
	return nil
}

type modelImpl struct {
	kod.Implements[testRepository]
}

func (t *modelImpl) Foo(ctx context.Context) error {
	return nil
}

type test1Component struct {
	kod.Implements[Test1Component]
	kod.WithConfig[test1Config]
}

func (t *test1Component) Init(ctx context.Context) error {
	kod := kod.FromContext(ctx)
	t.L(ctx).InfoContext(ctx, "Init test1Component"+kod.Config().Name)

	return nil
}

func (t *test1Component) Interceptors() []kod.Interceptor {
	return []kod.Interceptor{
		ktrace.Interceptor(),
		kmetric.Interceptor(),
		krecovery.Interceptor(),
		kratelimit.Interceptor(),
		kaccesslog.Interceptor(),
		kcircuitbreaker.Interceptor(),
		kvalidate.Interceptor(),
		ktimeout.Interceptor(ktimeout.WithTimeout(time.Second)),
	}
}

func (t *test1Component) Stop(ctx context.Context) error {
	return nil
}

type FooReq struct {
	Id    int `validate:"lt=100"`
	Panic bool
}

type FooRes struct {
	Id int
}

func (t *test1Component) Foo(ctx context.Context, req *FooReq) (*FooRes, error) {
	if req.Panic {
		panic("test panic")
	}

	ctx = baggage.ContextWithBaggage(ctx, lo.Must(baggage.New(lo.Must(baggage.NewMember("b1", "v1")))))
	t.L(ctx).InfoContext(ctx, "Foo info ", slog.Any("config", t.Config()))
	t.L(ctx).ErrorContext(ctx, "Foo error:")
	t.L(ctx).DebugContext(ctx, "Foo debug:")
	t.L(ctx).WithGroup("test group").InfoContext(ctx, "Foo info with group")

	return &FooRes{Id: req.Id}, errors.New("test1:" + t.Config().A)
}

type fakeTest1Component struct {
	A string
}

var _ Test1Component = (*fakeTest1Component)(nil)

func (f *fakeTest1Component) Foo(ctx context.Context, req *FooReq) (*FooRes, error) {
	fmt.Println(f.A)
	return nil, errors.New("A:" + f.A)
}

type test2Component struct {
	kod.Implements[Test2Component]
	kod.WithConfig[test1Config]
}

func (t *test2Component) GetClient() *http.Client {
	slog.Info("Foo info ", "config", t.Config())
	slog.Debug("Foo debug:")
	fmt.Println(errors.New("test1"))
	return &http.Client{}
}

type App struct {
	kod.Implements[kod.Main]
	test1 kod.Ref[Test1Component]
}

func Run(ctx context.Context, app *App) error {
	_, err := app.test1.Get().Foo(ctx, &FooReq{})
	return err
}

func StartTrace(ctx context.Context) context.Context {
	var opts []sdktrace.TracerProviderOption

	provider := sdktrace.NewTracerProvider(opts...)
	otel.SetTracerProvider(provider)

	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
	if err != nil {
		panic(err)
	} else {
		provider.RegisterSpanProcessor(sdktrace.NewSimpleSpanProcessor(exporter))
	}

	ctx, span := otel.Tracer("").Start(ctx, "Run")
	defer func() {
		span.End()
		fmt.Println("!!!!!!")
	}()

	return ctx
}
