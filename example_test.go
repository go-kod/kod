package kod_test

import (
	"context"
	"fmt"
	"log/slog"
	"testing"

	"github.com/go-kod/kod"
	"github.com/go-kod/kod/examples/helloworld"
	"github.com/go-kod/kod/interceptor"
	"github.com/go-kod/kod/interceptor/kmetric"
	"github.com/go-kod/kod/interceptor/krecovery"
	"github.com/go-kod/kod/interceptor/ktrace"
	"github.com/knadh/koanf/v2"
	"go.uber.org/mock/gomock"
)

// This example demonstrates how to use [kod.Run] and [kod.Implements] to run a simple application.
func Example_componentRun() {
	kod.Run(context.Background(), func(ctx context.Context, app *helloworld.App) error {
		fmt.Println("Hello, World!")
		return nil
	})
	// Output:
	// helloWorld init
	// Hello, World!
	// helloWorld shutdown
}

// This example demonstrates how to use [kod.MustRun] and [kod.Implements] to run a simple application.
func Example_componentRunMust() {
	kod.MustRun(context.Background(), func(ctx context.Context, app *helloworld.App) error {
		fmt.Println("Hello, World!")
		return nil
	})
	// Output:
	// helloWorld init
	// Hello, World!
	// helloWorld shutdown
}

// This example demonstrates how to use [kod.Ref] to reference a component and call a method on it.
func Example_componentRefAndCall() {
	kod.Run(context.Background(), func(ctx context.Context, app *helloworld.App) error {
		app.HelloWorld.Get().SayHello(ctx)
		return nil
	})
	// Output:
	// helloWorld init
	// Hello, World!
	// helloWorld shutdown
}

// This example demonstrates how to use [kod.LazyInit] to defer component initialization until it is needed.
func Example_componentLazyInit() {
	kod.Run(context.Background(), func(ctx context.Context, app *helloworld.App) error {
		app.HelloWorldLazy.Get().SayHello(ctx)
		app.HelloWorld.Get().SayHello(ctx)
		return nil
	})
	// Output:
	// helloWorld init
	// lazyHelloWorld init
	// Hello, Lazy!
	// Hello, World!
	// lazyHelloWorld shutdown
	// helloWorld shutdown
}

// This example demonstrates how to use [kod.WithFakes] and [kod.Fake] to provide a mock implementation of a component.
func Example_componentMock() {
	mock := helloworld.NewMockHelloWorld(gomock.NewController(nil))
	mock.EXPECT().SayHello(gomock.Any()).DoAndReturn(func(ctx context.Context) {
		fmt.Println("Hello, Mock!")
	})

	kod.Run(context.Background(), func(ctx context.Context, app *helloworld.App) error {
		app.HelloWorld.Get().SayHello(ctx)
		return nil
	}, kod.WithFakes(kod.Fake[helloworld.HelloWorld](mock)))
	// Output:
	// Hello, Mock!
}

// This example demonstrates how to use [kod.WithConfig] to provide a configuration to the application.
func Example_configInComponent() {
	kod.Run(context.Background(), func(ctx context.Context, app *helloworld.App) error {
		app.HelloWorld.Get().SayHello(ctx)
		return nil
	}, kod.WithConfigFile("./examples/helloworld/config.toml"))
	// Output:
	// helloWorld init
	// Hello, World!config
	// helloWorld shutdown
}

// This example demonstrates how to use [kod.WithGlobalConfig] to provide a global configuration to the application.
func Example_configGlobal() {
	kod.Run(context.Background(), func(ctx context.Context, app *helloworld.App) error {
		fmt.Println(app.Config().Name)
		return nil
	}, kod.WithConfigFile("./examples/helloworld/config.toml"))
	// Output:
	// helloWorld init
	// globalConfig
	// helloWorld shutdown
}

// This example demonstrates how to use logging with OpenTelemetry.
func Example_openTelemetryLog() {
	logger, observer := kod.NewTestLogger()
	slog.SetDefault(logger)

	kod.RunTest(&testing.T{}, func(ctx context.Context, app *helloworld.App) {
		app.L(ctx).Debug("Hello, World!")
		app.L(ctx).Info("Hello, World!")
		app.L(ctx).Warn("Hello, World!")
		app.L(ctx).Error("Hello, World!")
		app.HelloWorld.Get().SayHello(ctx)
	})

	fmt.Println(observer.RemoveKeys("trace_id", "span_id", "time"))

	// Output:
	// helloWorld init
	// Hello, World!
	// helloWorld shutdown
	// {"component":"github.com/go-kod/kod/Main","level":"INFO","msg":"Hello, World!"}
	// {"component":"github.com/go-kod/kod/Main","level":"WARN","msg":"Hello, World!"}
	// {"component":"github.com/go-kod/kod/Main","level":"ERROR","msg":"Hello, World!"}
	// {"component":"github.com/go-kod/kod/examples/helloworld/HelloWorld","level":"INFO","msg":"Hello, World!"}
}

// This example demonstrates how to use metrics with OpenTelemetry.
func Example_openTelemetryMetric() {
	kod.Run(context.Background(), func(ctx context.Context, app *helloworld.App) error {
		metric, _ := app.Meter().Int64Counter("example")
		metric.Add(ctx, 1)

		return nil
	})

	// Output:
	// helloWorld init
	// helloWorld shutdown
}

// This example demonstrates how to use [kod.WithInterceptors] to provide global interceptors to the application.
func Example_interceptorGlobal() {
	interceptor := interceptor.Interceptor(func(ctx context.Context, info interceptor.CallInfo, req, res []interface{}, next interceptor.HandleFunc) error {
		fmt.Println("Before call")
		err := next(ctx, info, req, res)
		fmt.Println("After call")
		return err
	})

	kod.Run(context.Background(), func(ctx context.Context, app *helloworld.App) error {
		app.HelloWorld.Get().SayHello(ctx)
		return nil
	}, kod.WithInterceptors(interceptor))
	// Output:
	// helloWorld init
	// Before call
	// Hello, World!
	// After call
	// helloWorld shutdown
}

// This example demonstrates how to use [kod.WithInterceptors] to provide component-specific interceptors to the application.
func Example_interceptorComponent() {
	kod.Run(context.Background(), func(ctx context.Context, app *helloworld.App) error {
		app.HelloWorldInterceptor.Get().SayHello(ctx)
		return nil
	})
	// Output:
	// helloWorld init
	// Before call
	// Hello, Interceptor!
	// After call
	// helloWorld shutdown
}

// This example demonstrates how to use built-in interceptors with [kod.WithInterceptors].
// Such as [krecovery.Interceptor], [ktrace.Interceptor], and [kmetric.Interceptor] ...
func Example_interceptorBuiltin() {
	kod.Run(context.Background(), func(ctx context.Context, app *helloworld.App) error {
		app.HelloWorld.Get().SayHello(ctx)
		return nil
	}, kod.WithInterceptors(krecovery.Interceptor(), ktrace.Interceptor(), kmetric.Interceptor()))
	// Output:
	// helloWorld init
	// Hello, World!
	// helloWorld shutdown
}

// This example demonstrates how to use [kod.RunTest] to run a test function.
func Example_testRun() {
	kod.RunTest(&testing.T{}, func(ctx context.Context, app *helloworld.App) {
		app.HelloWorld.Get().SayHello(ctx)
	})
	// Output:
	// helloWorld init
	// Hello, World!
	// helloWorld shutdown
}

// This example demonstrates how to use [kod.RunTest], [kod.Fake] and [kod.WithFakes] to run a test function with a mock component.
func Example_testWithMockComponent() {
	mock := helloworld.NewMockHelloWorld(gomock.NewController(nil))
	mock.EXPECT().SayHello(gomock.Any()).DoAndReturn(func(ctx context.Context) {
		fmt.Println("Hello, Mock!")
	})

	kod.RunTest(&testing.T{}, func(ctx context.Context, app *helloworld.App) {
		app.HelloWorld.Get().SayHello(ctx)
	}, kod.WithFakes(kod.Fake[helloworld.HelloWorld](mock)))
	// Output:
	// Hello, Mock!
}

// This example demonstrates how to use [kod.RunTest] and [kod.WithConfigFile] to run a test function with a configuration.
func Example_testWithConfig() {
	kod.RunTest(&testing.T{}, func(ctx context.Context, app *helloworld.App) {
		fmt.Println(app.Config().Name)
		app.HelloWorld.Get().SayHello(ctx)
	}, kod.WithConfigFile("./examples/helloworld/config.toml"))
	// Output:
	// helloWorld init
	// globalConfig
	// Hello, World!config
	// helloWorld shutdown
}

// This example demonstrates how to use kod.WithGlobalConfig with default configuration.
func Example_testWithDefaultConfig() {
	kod.RunTest(&testing.T{}, func(ctx context.Context, app *helloworld.App) {
		fmt.Println(app.Config().Name)
		app.HelloWorld.Get().SayHello(ctx)
	})
	// Output:
	// helloWorld init
	// kod
	// Hello, World!
	// helloWorld shutdown
}

// This example demonstrates how to use [kod.RunTest], [kod.NewTestLogger] to run a test function with a custom logger.
func Example_testWithLogObserver() {
	logger, observer := kod.NewTestLogger()
	slog.SetDefault(logger)

	t := &testing.T{}
	kod.RunTest(t, func(ctx context.Context, app *helloworld.App) {
		app.L(ctx).Debug("Hello, World!")
		app.L(ctx).Info("Hello, World!")
		app.L(ctx).Warn("Hello, World!")
		app.L(ctx).Error("Hello, World!")
	})

	fmt.Println(observer.Len())
	fmt.Println(observer.ErrorCount())
	fmt.Println(observer.Clean().Len())

	// Output:
	// helloWorld init
	// helloWorld shutdown
	// 3
	// 1
	// 0
}

// This example demonstrates how to use [kod.RunTest], [kod.WithKoanf] to run a test function with a custom koanf instance.
func Example_testWithKoanf() {
	c := koanf.New("_")
	c.Set("name", "testName")

	kod.RunTest(&testing.T{}, func(ctx context.Context, app *helloworld.App) {
		fmt.Println(app.Config().Name)
		app.HelloWorld.Get().SayHello(ctx)
	}, kod.WithKoanf(c))
	// Output:
	// helloWorld init
	// testName
	// Hello, World!
	// helloWorld shutdown
}
