package case1

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/samber/lo"
	"go.opentelemetry.io/otel/baggage"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/interceptor"
	"github.com/go-kod/kod/interceptor/kaccesslog"
	"github.com/go-kod/kod/interceptor/kcircuitbreaker"
	"github.com/go-kod/kod/interceptor/kmetric"
	"github.com/go-kod/kod/interceptor/kratelimit"
	"github.com/go-kod/kod/interceptor/krecovery"
	"github.com/go-kod/kod/interceptor/ktimeout"
	"github.com/go-kod/kod/interceptor/ktrace"
	"github.com/go-kod/kod/interceptor/kvalidate"
)

type test1Config struct {
	A     string
	Redis struct {
		Addr    string
		Timeout time.Duration
	}
}

func newTest1Config() test1Config {
	cfg := test1Config{A: "B"}
	cfg.Redis.Addr = "localhost:6379"
	cfg.Redis.Timeout = 2 * time.Second
	return cfg
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
	config test1Config
}

func (t *test1Component) Init(ctx context.Context) error {
	t.config = newTest1Config()
	t.L(ctx).InfoContext(ctx, "Init test1Componenttestapp")

	return nil
}

func (t *test1Component) Interceptors() []interceptor.Interceptor {
	return []interceptor.Interceptor{
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

func (t *test1Component) Shutdown(ctx context.Context) error {
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
	t.L(ctx).InfoContext(ctx, "Foo info ", slog.Any("config", t.config))
	t.L(ctx).ErrorContext(ctx, "Foo error:")
	t.L(ctx).DebugContext(ctx, "Foo debug:")
	t.L(ctx).WithGroup("test group").InfoContext(ctx, "Foo info with group")

	return &FooRes{Id: req.Id}, errors.New("test1:" + t.config.A)
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
	config test1Config
}

func (t *test2Component) Init(context.Context) error {
	t.config = newTest1Config()
	return nil
}

func (t *test2Component) GetClient() *http.Client {
	slog.Info("Foo info ", "config", t.config)
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

// func StartTrace(ctx context.Context) context.Context {
// 	var opts []sdktrace.TracerProviderOption

// 	provider := sdktrace.NewTracerProvider(opts...)
// 	otel.SetTracerProvider(provider)

// 	exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
// 	if err != nil {
// 		panic(err)
// 	} else {
// 		provider.RegisterSpanProcessor(sdktrace.NewSimpleSpanProcessor(exporter))
// 	}

// 	ctx, span := otel.Tracer("").Start(ctx, "Run")
// 	defer func() {
// 		span.End()
// 		fmt.Println("!!!!!!")
// 	}()

// 	return ctx
// }
